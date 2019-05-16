package u

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ParseInt(s string) int64 {
	if strings.IndexByte(s, '.') != -1 {
		i, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return int64(i)
		}
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return i
	}
	return 0
}

func ParseUint(s string) uint64 {
	if strings.IndexByte(s, '.') != -1 {
		i, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return uint64(i)
		}
	}
	i, err := strconv.ParseUint(s, 10, 64)
	if err == nil {
		return i
	}
	return 0
}

func Int(value interface{}) int {
	return int(Int64(value))
}
func Int64(value interface{}) int64 {
	switch realValue := value.(type) {
	case int:
		return int64(realValue)
	case int8:
		return int64(realValue)
	case int16:
		return int64(realValue)
	case int32:
		return int64(realValue)
	case int64:
		return int64(realValue)
	case float32:
		return int64(realValue)
	case float64:
		return int64(realValue)
	case bool:
		if realValue {
			return 1
		} else {
			return 0
		}
	case []byte:
		return ParseInt(string(realValue))
	case string:
		return ParseInt(realValue)
	}
	return 0
}

func Uint(value interface{}) uint {
	return uint(Uint64(value))
}
func Uint64(value interface{}) uint64 {
	switch realValue := value.(type) {
	case uint:
		return uint64(realValue)
	case uint8:
		return uint64(realValue)
	case uint16:
		return uint64(realValue)
	case uint32:
		return uint64(realValue)
	case uint64:
		return uint64(realValue)
	case float32:
		return uint64(realValue)
	case float64:
		return uint64(realValue)
	case bool:
		if realValue {
			return 1
		} else {
			return 0
		}
	case []byte:
		return ParseUint(string(realValue))
	case string:
		return ParseUint(realValue)
	}
	return 0
}

func Float(value interface{}) float32 {
	return float32(Float64(value))
}
func Float64(value interface{}) float64 {
	switch realValue := value.(type) {
	case int:
		return float64(realValue)
	case int8:
		return float64(realValue)
	case int16:
		return float64(realValue)
	case int32:
		return float64(realValue)
	case int64:
		return float64(realValue)
	case float32:
		return float64(realValue)
	case float64:
		return float64(realValue)
	case bool:
		if realValue {
			return 1
		} else {
			return 0
		}
	case []byte:
		i, err := strconv.ParseFloat(string(realValue), 10)
		if err == nil {
			return i
		}
	case string:
		i, err := strconv.ParseFloat(realValue, 10)
		if err == nil {
			return i
		}
	}
	return 0
}

func Bytes(value interface{}) []byte {
	return []byte(String(value))
}
func String(value interface{}) string {
	switch realValue := value.(type) {
	case int:
		return strconv.FormatInt(int64(realValue), 10)
	case int8:
		return strconv.FormatInt(int64(realValue), 10)
	case int16:
		return strconv.FormatInt(int64(realValue), 10)
	case int32:
		return strconv.FormatInt(int64(realValue), 10)
	case int64:
		return strconv.FormatInt(realValue, 10)
	case bool:
		if realValue {
			return "true"
		} else {
			return "false"
		}
	case string:
		return realValue
	case []byte:
		return string(realValue)
	}
	t := reflect.TypeOf(value)
	if t != nil {
		if t.Kind() == reflect.Struct || t.Kind() == reflect.Map || t.Kind() == reflect.Slice {
			return Json(value)
		}
	}
	return fmt.Sprint(value)
}

func Bool(value interface{}) bool {
	switch realValue := value.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		return realValue != 0
	case bool:
		return realValue
	case []byte:
		switch string(realValue) {
		case "1", "t", "T", "true", "TRUE", "True":
			return true
		}
	case string:
		switch realValue {
		case "1", "t", "T", "true", "TRUE", "True":
			return true
		}
	}
	return false
}

func Ints(value interface{}) []int64 {
	switch realValue := value.(type) {
	case []interface{}:
		result := make([]int64, len(realValue))
		for i, v := range realValue {
			result[i] = Int64(v)
		}
		return result
	default:
		return []int64{Int64(value)}
	}
	return make([]int64, 0)
}

func Floats(value interface{}) []float64 {
	switch realValue := value.(type) {
	case []interface{}:
		result := make([]float64, len(realValue))
		for i, v := range realValue {
			result[i] = Float64(v)
		}
		return result
	default:
		return []float64{Float64(value)}
	}
	return make([]float64, 0)
}

func Strings(value interface{}) []string {
	switch realValue := value.(type) {
	case []interface{}:
		result := make([]string, len(realValue))
		for i, v := range realValue {
			result[i] = String(v)
		}
		return result
	default:
		return []string{String(value)}
	}
	return make([]string, 0)
}

func GetLowerName(s string) string {
	buf := []byte(s)
	if buf[0] >= 'A' && buf[0] <= 'Z' {
		buf[0] += 32
	}
	return string(buf)
}

func GetUpperName(s string) string {
	buf := []byte(s)
	if buf[0] >= 'a' && buf[0] <= 'z' {
		buf[0] -= 32
	}
	return string(buf)
}

func FixUpperCase(data []byte, excludesKeys []string) {
	n := len(data)
	types := make([]bool, 0)
	keys := make([]string, 0)
	tpos := -1

	for i := 0; i < n-1; i++ {
		if tpos+1 >= len(types) {
			types = append(types, false)
			keys = append(keys, "")
		}

		if data[i] == '{' {
			tpos++
			types[tpos] = true
			keys[tpos] = ""
			//log.Println(" >>>1 ", types, tpos)
		} else if data[i] == '}' {
			tpos--
			//log.Println(" >>>2 ", types, tpos)
		}
		if data[i] == '[' {
			tpos++
			types[tpos] = false
			keys[tpos] = ""
			//log.Println(" >>>3 ", types, tpos)
		} else if data[i] == ']' {
			tpos--
			//log.Println(" >>>4 ", types, tpos)
		}
		if data[i] == '"' {
			keyPos := -1
			if data[i-1] == '{' || (data[i-1] == ',' && tpos >= 0 && types[tpos]) {
				keyPos = i + 1
			}
			// skip string
			i++
			for ; i < n-1; i++ {
				if data[i] == '\\' {
					i++
					continue
				}
				if data[i] == '"' {
					if keyPos >= 0 && excludesKeys != nil {
						keys[tpos] = string(data[keyPos:i])
					}
					break
				}
			}

			if keyPos >= 0 && (data[keyPos] >= 'A' && data[keyPos] <= 'Z') {
				if excludesKeys != nil {
					// 是否排除
					excluded := false
					for _, ek := range excludesKeys {
						for j := tpos - 1; j >= 0; j-- {
							if strings.Index(keys[j], ek) != -1 {
								excluded = true
								break
							}
						}
						if excluded {
							break
						}
					}
					//fmt.Println("  ** >>", keys, excluded)
					if !excluded {
						data[keyPos] += 32
					}
				} else {
					// 不进行排除判断
					data[keyPos] += 32
				}
			}
			continue
		}
	}
}

func If(i bool, a, b interface{}) interface{} {
	if i {
		return a
	}
	return b
}

func StringIf(i bool, a, b string) string {
	if i {
		return a
	}
	return b
}

func Switch(i uint, args ...interface{}) interface{} {
	if i < uint(len(args)) {
		return args[i]
	}
	return nil
}

func StringIn(arr []string, s string) bool {
	for _, d := range arr {
		if d == s {
			return true
		}
	}
	return false
}

func In(arr []interface{}, s interface{}) bool {
	for _, d := range arr {
		if d == s {
			return true
		}
	}
	return false
}

//func Elem(v reflect.Value) reflect.Value {
//	for v.Kind() == reflect.Ptr {
//		v = v.Elem()
//	}
//	return v
//}

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
	for j := len(keys) - 1; j >= 0; j-- {
		keyMap[strings.ToLower(keys[j].String())] = &keys[j]
	}

	toType := to.Type()
	for i := toType.NumField() - 1; i >= 0; i-- {
		f := toType.Field(i)
		if f.Anonymous {
			convertMapToStruct(from, to.Field(i))
			continue
		}

		k := keyMap[strings.ToLower(f.Name)]
		var v reflect.Value
		if k != nil {
			v = from.MapIndex(*k)
		}

		if v.IsValid() && !v.IsNil() {
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
		convert(v, valueItem)
		to.SetMapIndex(keyItem, valueItem)
	}
}

func convertStructToMap(from, to reflect.Value) {
	toType := to.Type()
	for i := from.NumField() - 1; i >= 0; i-- {
		k := from.Type().Field(i).Name
		v := from.Field(i)
		keyItem := reflect.New(toType.Key()).Elem()
		valueItem := reflect.New(toType.Elem()).Elem()
		convert(k, keyItem)
		convert(v, valueItem)
		to.SetMapIndex(keyItem, valueItem)
	}
}

func convertSliceToSlice(from, to reflect.Value) *reflect.Value {
	toType := to.Type()
	fromNum := from.Len()
	for i := 0; i < fromNum; i++ {
		valueItem := reflect.New(toType.Elem()).Elem()
		convert(from.Index(i), valueItem)
		to = reflect.Append(to, valueItem)
	}
	return &to
}

func Convert(from, to interface{}) interface{} {
	r := convert(from, to)
	if r == nil {
		return to
	} else {
		return r.Interface()
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
		if fromType.Kind() == reflect.Slice {
			return convertSliceToSlice(fromValue, toValue)
		}
		if toType.Kind() == reflect.Slice && toType.Elem().Kind() == reflect.Uint8 {
			toValue.SetBytes(Bytes(from))
		}
	case reflect.Struct:
		switch fromType.Kind() {
		case reflect.Map:
			convertMapToStruct(fromValue, toValue)
		case reflect.Struct:
			//convertStructToStruct(fromValue, toValue)
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

func Json(value interface{}) string {
	j, err := json.Marshal(value)
	if err == nil {
		return string(j)
	}
	return fmt.Sprint(value)
}

func JsonP(value interface{}) string {
	j, err := json.MarshalIndent(value, "", "  ")
	if err == nil {
		return string(j)
	}
	return fmt.Sprint(value)
}
