package u

import (
	"reflect"
	"strings"
)

func FinalType(v reflect.Value) reflect.Type {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Interface {
		return reflect.TypeOf(v.Interface())
	} else {
		return v.Type()
	}
}

func RealValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func FinalValue(v reflect.Value) reflect.Value {
	v = RealValue(v)
	if v.Kind() == reflect.Interface {
		return v.Elem()
	} else {
		return v
	}
}

func FixNilValue(v reflect.Value) {
	t := v.Type()
	for t.Kind() == reflect.Ptr && v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
		v = v.Elem()
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice && v.IsNil() {
		v.Set(reflect.MakeSlice(v.Type(), 0, 0))
	}
	if t.Kind() == reflect.Map && v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}
}

func convertMapToStruct(from, to reflect.Value) {
	keys := from.MapKeys()
	keyMap := map[string]*reflect.Value{}
	fixedKeyMap := map[string]*reflect.Value{} // match with fixed '-' & '_'
	for j := len(keys) - 1; j >= 0; j-- {
		keyStr := ValueToString(keys[j])
		keyMap[strings.ToLower(keyStr)] = &keys[j]
		if strings.ContainsRune(keyStr, '-') || strings.ContainsRune(keyStr, '_') {
			fixedKeyMap[strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(keyStr), "-", ""), "_", "")] = &keys[j]
		}
	}

	toType := to.Type()
	for i := toType.NumField() - 1; i >= 0; i-- {
		f := toType.Field(i)
		if f.Anonymous {
			convertMapToStruct(from, to.Field(i))
			continue
		}

		if toType.Field(i).Name[0] > 90 {
			continue
		}

		k := keyMap[strings.ToLower(f.Name)]
		if k == nil {
			k = fixedKeyMap[strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(f.Name), "-", ""), "_", "")]
		}
		var v reflect.Value
		if k != nil {
			v = from.MapIndex(*k)
		}

		if v.IsValid() {
			r := convert(v, to.Field(i))
			if r != nil {
				to.Field(i).Set(*r)
			}
		}
	}
}

func convertStructToStruct(from, to reflect.Value) {
	keyMap := map[string]int{}
	fixedKeyMap := map[string]int{}
	fromType := from.Type()
	for i := fromType.NumField() - 1; i >= 0; i-- {
		if fromType.Field(i).Name[0] > 90 {
			continue
		}
		keyMap[strings.ToLower(fromType.Field(i).Name)] = i + 1
		if strings.ContainsRune(fromType.Field(i).Name, '_') {
			fixedKeyMap[strings.ReplaceAll(strings.ToLower(fromType.Field(i).Name), "_", "")] = i + 1
		}
	}

	toType := to.Type()
	for i := toType.NumField() - 1; i >= 0; i-- {
		f := toType.Field(i)
		if f.Anonymous {
			convertStructToStruct(from, to.Field(i))
			continue
		}

		if f.Name[0] > 90 {
			continue
		}

		k := keyMap[strings.ToLower(f.Name)]
		if k == 0 {
			k = fixedKeyMap[strings.ReplaceAll(strings.ToLower(f.Name), "_", "")]
		}
		var v reflect.Value
		if k != 0 {
			v = from.Field(k - 1)
		}

		if v.IsValid() {
			r := convert(v, to.Field(i))
			if r != nil {
				to.Field(i).Set(*r)
			}
		}
	}
}

func convertMapToMap(from, to reflect.Value) {
	toType := to.Type()
	keys := from.MapKeys()
	keyNum := len(keys)
	for i := 0; i < keyNum; i++ {
		k := keys[i]
		v := from.MapIndex(k)
		keyItem := reflect.New(toType.Key()).Elem()
		valueItem := reflect.New(toType.Elem()).Elem()
		convert(k, keyItem)
		newItem := convert(v, valueItem)
		if newItem != nil {
			to.SetMapIndex(keyItem, *newItem)
		} else {
			to.SetMapIndex(keyItem, valueItem)
		}
	}
}

func convertStructToMap(from, to reflect.Value) {
	toType := to.Type()
	for i := from.NumField() - 1; i >= 0; i-- {
		k := from.Type().Field(i).Name
		v := from.Field(i)
		if k[0] > 90 {
			continue
		}
		keyItem := reflect.New(toType.Key()).Elem()
		valueItem := reflect.New(toType.Elem()).Elem()
		convert(k, keyItem)
		if keyItem.Kind() == reflect.String {
			// Struct转Map时自动将首字母改为小写
			keyStr := keyItem.String()
			if len(keyStr) > 0 && keyStr[0] >= 'A' && keyStr[0] <= 'Z' {
				keyBuf := []byte(keyStr)
				keyBuf[0] += 32
				keyItem = reflect.ValueOf(string(keyBuf))
			}
		}
		newItem := convert(v, valueItem)
		if newItem != nil {
			to.SetMapIndex(keyItem, *newItem)
		} else {
			to.SetMapIndex(keyItem, valueItem)
		}
	}
}

func convertSliceToSlice(from, to reflect.Value) *reflect.Value {
	toType := to.Type()
	fromNum := from.Len()
	for i := 0; i < fromNum; i++ {
		valueItem := reflect.New(toType.Elem()).Elem()
		newItem := convert(from.Index(i), valueItem)
		if newItem != nil {
			to = reflect.Append(to, *newItem)
		} else {
			to = reflect.Append(to, valueItem)
		}
	}
	return &to
}

func Convert(from, to interface{}) {
	r := convert(from, to)
	if r != nil {
		toValue := reflect.ValueOf(to)
		var prevValue reflect.Value
		for toValue.Kind() == reflect.Ptr {
			prevValue = toValue
			toValue = toValue.Elem()
		}
		if prevValue.IsValid() {
			prevValue.Elem().Set(*r)
		}
	}
}

func convert(from, to interface{}) *reflect.Value {
	var fromValue reflect.Value
	var toValue reflect.Value
	if v, ok := from.(reflect.Value); ok {
		from = v.Interface()
		fromValue = v
	} else {
		fromValue = reflect.ValueOf(from)
	}
	if v, ok := to.(reflect.Value); ok {
		toValue = v
	} else {
		toValue = reflect.ValueOf(to)
	}
	//originToValue := toValue
	FixNilValue(toValue)

	fromValue = FinalValue(fromValue)
	toValue = RealValue(toValue)
	if !fromValue.IsValid() || !toValue.IsValid() {
		return nil
	}

	fromType := FinalType(fromValue)
	toType := toValue.Type()

	switch toType.Kind() {
	case reflect.Bool:
		toValue.SetBool(Bool(from))
	case reflect.Interface:
		toValue.Set(reflect.ValueOf(from))
	case reflect.String:
		toValue.SetString(String(from))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		toValue.SetInt(Int64(from))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		toValue.SetUint(Uint64(from))
	case reflect.Float32, reflect.Float64:
		toValue.SetFloat(Float64(from))
	case reflect.Slice:
		if toType.Kind() == reflect.Slice && toType.Elem().Kind() == reflect.Uint8 {
			toValue.SetBytes(Bytes(from))
		} else if fromType.Kind() == reflect.Slice {
			return convertSliceToSlice(fromValue, toValue)
		} else {
			tmpSlice := reflect.MakeSlice(reflect.SliceOf(fromType), 1, 1)
			tmpSlice.Index(0).Set(fromValue)
			return convertSliceToSlice(tmpSlice, toValue)
		}
	case reflect.Struct:
		switch fromType.Kind() {
		case reflect.Map:
			convertMapToStruct(fromValue, toValue)
		case reflect.Struct:
			convertStructToStruct(fromValue, toValue)
		}
	case reflect.Map:
		switch fromType.Kind() {
		case reflect.Map:
			convertMapToMap(fromValue, toValue)
		case reflect.Struct:
			convertStructToMap(fromValue, toValue)
		}
	default:
		//fmt.Println(" !!!!!!2", fromType.Kind(), toType.Kind(), toType.Elem().Kind())
	}
	return nil
}

func ToInterfaceArray(in interface{}) []interface{} {
	v := FinalValue(reflect.ValueOf(in))
	out := make([]interface{}, 0)
	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			if v.Index(i).CanInterface() {
				out = append(out, v.Index(i).Interface())
			}
		}
	}
	return out
}
