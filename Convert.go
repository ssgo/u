package u

import (
	"encoding/json"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

func FixPtr(value interface{}) interface{} {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = FinalValue(v)
		if v.IsValid() {
			return v.Interface()
		} else {
			return nil
		}
	}
	return value
}

func FinalType(v reflect.Value) reflect.Type {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Interface {
		t := reflect.TypeOf(v.Interface())
		for t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		return t
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
		return RealValue(v.Elem())
	} else {
		return v
	}
}

func FixNilValue(v reflect.Value) {
	t := v.Type()
	for t.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
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
		//fmt.Println("##########", f.Name, k.Interface(), FinalType(v))

		if v.IsValid() {
			// 支持 Parse 方法对数据进行转换
			parsed := false
			if to.CanAddr() {
				toP := to.Addr()
				if m, ok := toP.Type().MethodByName("Parse" + toType.Field(i).Name); ok && m.Type.NumIn() == 2 && m.Type.NumOut() == 1 {
					if m.Type.In(0).String() == toP.Type().String() && m.Type.Out(0).String() == to.Field(i).Type().String() {
						//fmt.Println(" ===== Parse"+toType.Field(i).Name, m.Type.In(0).String(), m.Type.In(1), m.Type.Out(0), to.Field(i).Type().String())
						argP := reflect.New(m.Type.In(1))
						vF := FinalValue(v)
						r := convert(vF, argP)
						var argV reflect.Value
						if r != nil {
							argV = *r
						} else {
							argV = argP.Elem()
						}
						out := m.Func.Call([]reflect.Value{toP, argV})
						//fmt.Println("  >>>", JsonP(out[0].Interface()))
						to.Field(i).Set(out[0])
						parsed = true
					}
				}
			}
			if !parsed {
				r := convert(v, to.Field(i))
				if r != nil {
					FinalSet(*r, to.Field(i))
					//to.Field(i).Set(*r)
				}
			}
		}
	}
}

func FinalSet(from, to reflect.Value) {
	fv := FinalValue(from)
	tv := FinalValue(to)
	tv.Set(fv)
}

func convertStructToStruct(from, to reflect.Value) {
	keyMap := map[string]int{}
	fixedKeyMap := map[string]int{}
	fromType := from.Type()
	toType := to.Type()

	// copy when same type
	ft := FinalType(from)
	tt := FinalType(to)
	if ft == tt {
		FinalSet(from, to)
		//fv := FinalValue(from)
		//tv := FinalValue(to)
		//tv.Set(fv)
		return
	}

	for i := fromType.NumField() - 1; i >= 0; i-- {
		if fromType.Field(i).Name[0] > 90 {
			continue
		}
		keyMap[strings.ToLower(fromType.Field(i).Name)] = i + 1
		if strings.ContainsRune(fromType.Field(i).Name, '_') {
			fixedKeyMap[strings.ReplaceAll(strings.ToLower(fromType.Field(i).Name), "_", "")] = i + 1
		}
	}

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
		convert(k, keyItem)
		valueItem := to.MapIndex(keyItem)
		if !valueItem.IsValid() {
			valueItem = reflect.New(toType.Elem()).Elem()
		}
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
		// valueItem := reflect.New(toType.Elem()).Elem()
		convert(k, keyItem)
		valueItem := to.MapIndex(keyItem)
		if !valueItem.IsValid() {
			valueItem = reflect.New(toType.Elem()).Elem()
		}
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
	var newValue *reflect.Value = nil

	if toValue.CanAddr() {
		jsonUM, jsonUMOk := toValue.Addr().Interface().(json.Unmarshaler)
		yamlUM, yamlUMOk := toValue.Addr().Interface().(yaml.Unmarshaler)
		if jsonUMOk {
			jsonUM.UnmarshalJSON(Bytes(fromValue.Interface()))
			return nil
		} else if yamlUMOk {
			yamlUM.UnmarshalYAML(&yaml.Node{
				Value: String(fromValue.Interface()),
			})
			return nil
		}
	}

	toValueP := toValue
	toType := toValue.Type()
	if toValue.IsValid() && toValue.Type().Kind() == reflect.Interface {
		if toValue.Elem().IsValid() {
			toValue = toValue.Elem()
			toType = toValue.Type()
		}
	}

	switch toType.Kind() {
	case reflect.Bool:
		if toValue != toValueP {
			if toValueP.CanAddr() {
				toValueP.Set(fromValue)
			} else {
				newValue = &fromValue
			}
		} else {
			toValue.SetBool(Bool(fromValue.Interface()))
		}
	case reflect.Interface:
		toValueP.Set(reflect.ValueOf(fromValue.Interface()))
	case reflect.String:
		if toValue != toValueP {
			if toValueP.CanAddr() {
				toValueP.Set(fromValue)
			} else {
				newValue = &fromValue
			}
		} else {
			if toValueP.CanAddr() {
				toValue.SetString(String(fromValue.Interface()))
				// newValue = &toValue
			} else {
				newValue = &fromValue
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if toValue != toValueP {
			if toValueP.CanAddr() {
				toValueP.Set(fromValue)
			} else {
				newValue = &fromValue
			}
		} else {
			toValue.SetInt(Int64(fromValue.Interface()))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if toValue != toValueP {
			if toValueP.CanAddr() {
				toValueP.Set(fromValue)
			} else {
				newValue = &fromValue
			}
		} else {
			toValue.SetUint(Uint64(fromValue.Interface()))
		}
	case reflect.Float32, reflect.Float64:
		if toValue != toValueP {
			if toValueP.CanAddr() {
				toValueP.Set(fromValue)
			} else {
				newValue = &fromValue
			}
		} else {
			toValue.SetFloat(Float64(fromValue.Interface()))
		}
	case reflect.Slice:
		if toType.Elem().Kind() == reflect.Uint8 {
			toValue.SetBytes(Bytes(fromValue.Interface()))
		} else if fromType.Kind() == reflect.Slice {
			return convertSliceToSlice(fromValue, toValue)
		} else if fromType.Kind() == reflect.String {
			return convertSliceToSlice(reflect.ValueOf(UnJsonArr(fromValue.String())), toValue)
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
		case reflect.String:
			convertMapToStruct(reflect.ValueOf(UnJsonMap(fromValue.String())), toValue)
		}
	case reflect.Map:
		if toValue.IsNil() {
			toValue = reflect.MakeMap(toType)
			newValue = &toValue
		}
		switch fromType.Kind() {
		case reflect.Map:
			convertMapToMap(fromValue, toValue)
		case reflect.Struct:
			convertStructToMap(fromValue, toValue)
		case reflect.String:
			convertMapToMap(reflect.ValueOf(UnJsonMap(fromValue.String())), toValue)
		}
	case reflect.Func:
		if fromType.Kind() == reflect.Func {
			toValue.Set(reflect.MakeFunc(toType, func(goArgs []reflect.Value) []reflect.Value {
				ins := make([]reflect.Value, 0)
				j := 0
				for i := 0; i < len(goArgs); i++ {
					if j >= fromType.NumIn() {
						break
					}
					var jV interface{}
					if fromType.IsVariadic() && j == fromType.NumIn()-1 && fromType.In(j).Kind() == reflect.Slice {
						jV = reflect.New(fromType.In(j).Elem()).Interface()
					} else {
						jV = reflect.New(fromType.In(j)).Interface()
						j++
					}
					convert(goArgs[i].Interface(), jV)
					ins = append(ins, reflect.ValueOf(jV).Elem())
				}
				out := fromValue.Call(ins)
				outs := make([]reflect.Value, 0)
				j = 0
				for i := 0; i < toType.NumOut(); i++ {
					iV := reflect.New(toType.Out(i)).Interface()
					var jV interface{}
					if toType.NumOut() > len(out) {
						if j == len(out)-1 && out[j].Kind() == reflect.Slice {
							if out[j].Len() > i-j {
								jV = out[j].Index(i - j).Interface()
								convert(jV, iV)
							}
							outs = append(outs, reflect.ValueOf(iV).Elem())
						} else {
							// not match, use default value
							if toType.Kind() == reflect.Ptr {
								outs = append(outs, reflect.ValueOf(nil))
							} else {
								outs = append(outs, reflect.New(toType.Out(i)).Elem())
							}
						}
					} else {
						jV = out[j].Interface()
						convert(jV, iV)
						j++
						outs = append(outs, reflect.ValueOf(iV).Elem())
					}
				}
				return outs
			}))

			//jsFunc := args[i]
			//funcType := needArgType
			//argValue = reflect.MakeFunc(funcType, func(goArgs []reflect.Value) []reflect.Value {
			//	ins := make([]quickjs.Value, 0)
			//	for _, goArg := range goArgs {
			//		ins = append(ins, MakeJsValue(ctx, goArg.Interface(), false))
			//	}
			//	outs := make([]reflect.Value, 0)
			//	for j := 0; j < funcType.NumOut(); j++ {
			//		outs = append(outs, reflect.New(funcType.Out(j)).Elem())
			//	}
			//	jsResult := jsCtx.Invoke(jsFunc, jsCtx.Null(), ins...)
			//	if !jsResult.IsUndefined() && len(outs) > 0 {
			//		out0P := outs[0].Interface()
			//		u.Convert(MakeFromJsValue(jsResult), out0P)
			//		outs[0] = reflect.ValueOf(out0P).Elem()
			//	}
			//	return outs
			//})
		}
	default:
		//fmt.Println(" !!!!!!2", fromType.Kind(), toType.Kind(), toType.Elem().Kind())
	}
	return newValue
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

func SetValue(to, from reflect.Value) {
	if to.CanSet() {
		if from.Kind() == to.Kind() {
			to.Set(from)
		} else if to.Kind() == reflect.Ptr && from.Kind() == to.Type().Elem().Kind() {
			newValue := reflect.New(to.Type().Elem())
			newValue.Elem().Set(from)
			to.Set(newValue)
		} else if from.Kind() == reflect.Ptr && from.Elem().Kind() == to.Kind() {
			to.Set(from.Elem())
		} else {
			newValue := reflect.New(to.Type())
			Convert(from.Interface(), newValue.Interface())
			if to.Kind() == reflect.Ptr {
				to.Set(newValue)
			} else {
				to.Set(newValue.Elem())
			}
		}
	}
}

type StructInfo struct {
	Fields       []reflect.StructField
	Values       map[string]reflect.Value
	Methods      []reflect.Method
	MethodValues map[string]reflect.Value
}

func FlatStruct(data interface{}) *StructInfo {
	return flatStruct(data, true)
}

func FlatStructWithUnexported(data interface{}) *StructInfo {
	return flatStruct(data, false)
}

func flatStruct(data interface{}, onlyExported bool) *StructInfo {
	out := &StructInfo{
		Fields:       make([]reflect.StructField, 0),
		Values:       make(map[string]reflect.Value),
		Methods:      make([]reflect.Method, 0),
		MethodValues: make(map[string]reflect.Value),
	}
	if v, ok := data.(reflect.Value); ok {
		makeStructInfo(v, out, onlyExported)
	} else {
		makeStructInfo(reflect.ValueOf(data), out, onlyExported)
	}
	return out
}

func makeStructInfo(v reflect.Value, out *StructInfo, onlyExported bool) {
	fv := v
	for v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct {
		fv = v.Elem()
	}
	t := v.Type()
	ft := fv.Type()

	if fv.Kind() == reflect.Struct {
		if v.Kind() == reflect.Ptr {
			for i := 0; i < v.NumMethod(); i++ {
				if onlyExported && !t.Method(i).IsExported() {
					continue
				}
				out.Methods = append(out.Methods, t.Method(i))
				out.MethodValues[t.Method(i).Name] = v.Method(i)
			}
		}
		for i := 0; i < ft.NumField(); i++ {
			if onlyExported && !ft.Field(i).IsExported() {
				continue
			}
			if ft.Field(i).Anonymous {
				makeStructInfo(fv.Field(i), out, onlyExported)
			} else {
				out.Fields = append(out.Fields, ft.Field(i))
				out.Values[ft.Field(i).Name] = fv.Field(i)
			}
		}
	}
}
