package u

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"reflect"
	"strconv"
	"strings"
	"time"
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
	if value == nil {
		return 0
	}
	value = FixPtr(value)
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
		return realValue
	case uint:
		return int64(realValue)
	case uint8:
		return int64(realValue)
	case uint16:
		return int64(realValue)
	case uint32:
		return int64(realValue)
	case uint64:
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
	if value == nil {
		return 0
	}
	value = FixPtr(value)
	switch realValue := value.(type) {
	case int:
		return uint64(realValue)
	case int8:
		return uint64(realValue)
	case int16:
		return uint64(realValue)
	case int32:
		return uint64(realValue)
	case int64:
		return uint64(realValue)
	case uint:
		return uint64(realValue)
	case uint8:
		return uint64(realValue)
	case uint16:
		return uint64(realValue)
	case uint32:
		return uint64(realValue)
	case uint64:
		return realValue
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
	if value == nil {
		return 0
	}
	value = FixPtr(value)
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
	case uint:
		return float64(realValue)
	case uint8:
		return float64(realValue)
	case uint16:
		return float64(realValue)
	case uint32:
		return float64(realValue)
	case uint64:
		return float64(realValue)
	case float32:
		return float64(realValue)
	case float64:
		return realValue
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
	return _string(value, false)
}

func StringP(value interface{}) string {
	return _string(value, true)
}

func _string(value interface{}, p bool) string {
	if value == nil {
		return ""
	}
	value = FixPtr(value)
	if value == nil {
		return ""
	}
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
	case uint:
		return strconv.FormatInt(int64(realValue), 10)
	case uint8:
		return strconv.FormatInt(int64(realValue), 10)
	case uint16:
		return strconv.FormatInt(int64(realValue), 10)
	case uint32:
		return strconv.FormatInt(int64(realValue), 10)
	case uint64:
		return strconv.FormatInt(int64(realValue), 10)
	case float32:
		return strconv.FormatFloat(float64(realValue), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(realValue, 'f', -1, 64)
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
	if t != nil && (t.Kind() == reflect.Struct || t.Kind() == reflect.Map || t.Kind() == reflect.Slice) {
		//j, err := json.Marshal(value)
		//if err == nil {
		//	return string(FixJsonBytes(j))
		//}
		//return fmt.Sprint(value)
		if p {
			return JsonP(value)
		} else {
			return Json(value)
		}
	}
	return fmt.Sprint(value)
}

func Bool(value interface{}) bool {
	value = FixPtr(value)
	switch realValue := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return Uint64(realValue) != 0
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
	value = FixPtr(value)
	switch realValue := value.(type) {
	case []interface{}:
		result := make([]int64, len(realValue))
		for i, v := range realValue {
			result[i] = Int64(v)
		}
		return result
	case string:
		if strings.HasPrefix(realValue, "[") {
			result := make([]int64, 0)
			UnJson(realValue, &result)
			return result
		} else {
			return []int64{Int64(value)}
		}
	default:
		return []int64{Int64(value)}
	}
}

func Floats(value interface{}) []float64 {
	value = FixPtr(value)
	switch realValue := value.(type) {
	case []interface{}:
		result := make([]float64, len(realValue))
		for i, v := range realValue {
			result[i] = Float64(v)
		}
		return result
	case string:
		if strings.HasPrefix(realValue, "[") {
			result := make([]float64, 0)
			UnJson(realValue, &result)
			return result
		} else {
			return []float64{Float64(value)}
		}
	default:
		return []float64{Float64(value)}
	}
}

func Strings(value interface{}) []string {
	value = FixPtr(value)
	switch realValue := value.(type) {
	case []interface{}:
		result := make([]string, len(realValue))
		for i, v := range realValue {
			result[i] = String(v)
		}
		return result
	case string:
		if strings.HasPrefix(realValue, "[") {
			result := make([]string, 0)
			UnJson(realValue, &result)
			return result
		} else {
			return []string{String(value)}
		}
	default:
		return []string{String(value)}
	}
}

func Duration(value string) time.Duration {
	result, err := time.ParseDuration(value)
	if err != nil {
		return time.Duration(Int64(value)) * time.Millisecond
	} else {
		return result
	}
}

func ValueToString(v reflect.Value) string {
	if v.Kind() == reflect.String {
		return v.String()
	} else {
		return fmt.Sprint(v)
	}
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

func StringFromValue(v reflect.Value) string {
	if v.CanInterface() {
		return String(v.Interface())
	} else {
		return v.String()
	}
}

func MakeExcludeUpperKeys(data interface{}, prefix string) []string {
	oriPrefix := prefix
	if prefix != "" {
		prefix += "."
	}
	outs := make([]string, 0)

	var dataV reflect.Value
	if v, ok := data.(reflect.Value); ok {
		dataV = v
	} else {
		dataV = reflect.ValueOf(data)
	}
	inValue := FinalValue(dataV)
	if !inValue.IsValid() {
		return nil
	}

	inType := inValue.Type()
	switch inType.Kind() {
	case reflect.Map:
		for _, k := range inValue.MapKeys() {
			r := MakeExcludeUpperKeys(inValue.MapIndex(k), prefix+StringFromValue(k))
			if len(r) > 0 {
				outs = append(outs, r...)
			}
		}
	case reflect.Slice:
		if inType.Elem().Kind() != reflect.Uint8 {
			for i := inValue.Len() - 1; i >= 0; i-- {
				r := MakeExcludeUpperKeys(inValue.Index(i), prefix)
				if len(r) > 0 {
					outs = append(outs, r...)
				}
			}
		}
	case reflect.Struct:
		for i := inType.NumField() - 1; i >= 0; i-- {
			f := inType.Field(i)
			if f.Anonymous {
				r := MakeExcludeUpperKeys(inValue.Field(i), oriPrefix)
				if len(r) > 0 {
					outs = append(outs, r...)
				}
			} else {
				if strings.Contains(String(f.Tag), "keepKey") {
					outs = append(outs, prefix+f.Name)
				}
				if strings.Contains(String(f.Tag), "keepSubKey") {
					outs = append(outs, prefix+f.Name+".")
				}
				r := MakeExcludeUpperKeys(inValue.Field(i), prefix+f.Name)
				if len(r) > 0 {
					outs = append(outs, r...)
				}
			}
		}
	}
	return outs
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
			if i > 0 && (data[i-1] == '{' || (data[i-1] == ',' && tpos >= 0 && types[tpos])) {
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
					if keyPos >= 0 && excludesKeys != nil && len(excludesKeys) > 0 {
						keys[tpos] = string(data[keyPos:i])
					}
					break
				}
			}

			if keyPos >= 0 && (data[keyPos] >= 'A' && data[keyPos] <= 'Z') {
				if excludesKeys != nil && len(excludesKeys) > 0 {
					// 是否排除
					excluded := false
					checkStr := strings.Join(keys[0:tpos+1], ".")
					for _, ek := range excludesKeys {
						//fmt.Println(".  >", ek, checkStr)
						//if checkStr == ek {
						//	excluded = true
						//}
						//if !excluded && strings.HasPrefix(checkStr, ek+".") {
						//	excluded = true
						//}
						// if set "requestHeaders" mean requestHeaders is excluded, but children is not excluded
						// if set "requestHeaders." mean requestHeaders is not excluded, but children is excluded
						if strings.HasSuffix(ek, ".") {
							excluded = strings.HasPrefix(checkStr, ek)
						} else {
							excluded = checkStr == ek
						}
						//for j := tpos - 1; j >= 0; j-- {
						//	if strings.Index(keys[j], ek) != -1 {
						//		excluded = true
						//		break
						//	}
						//}
						if excluded {
							break
						}
					}
					if !excluded {
						keyStr := keys[tpos]
						hasLower := false
						for c := len(keyStr) - 1; c >= 0; c-- {
							if keyStr[c] >= 'a' && keyStr[c] <= 'z' {
								hasLower = true
								break
							}
						}
						// 不转换全大写的Key
						if hasLower {
							data[keyPos] += 32
						}
					}
				} else {
					// 不进行排除判断
					hasLower := false
					dataLen := len(data)
					for c := keyPos; c < dataLen; c++ {
						if data[c] == '"' {
							break
						}
						if data[c] >= 'a' && data[c] <= 'z' {
							hasLower = true
							break
						}
					}
					// 不转换全大写的Key
					if hasLower {
						data[keyPos] += 32
					}
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

func SplitTrim(s, sep string) []string {
	ss := strings.Split(s, sep)
	for i, s1 := range ss {
		ss[i] = strings.TrimSpace(s1)
	}
	return ss
}

func SplitWithoutNone(s, sep string) []string {
	if s == "" {
		return []string{}
	} else {
		return SplitTrim(s, sep)
	}
}

func SplitArgs(s string) []string {
	a := make([]string, 0)
	chars := []rune(s)
	inQuote := false
	for i := range chars {
		c := chars[i]
		prevC := rune(0)
		if i > 0 {
			prevC = chars[i-1]
		}
		if c == '"' && prevC != '\\' && ((!inQuote && (i == 0 || chars[i-1] == ' ')) || (inQuote && (i == len(s)-1 || len(chars) <= i+1 || chars[i+1] == ' '))) {
			inQuote = !inQuote
		} else {
			a = append(a, StringIf(c == ' ' && inQuote, "__SPACE__", string(c)))
		}
	}

	s = strings.Join(a, "")
	s = strings.ReplaceAll(s, "\\\"", "\"")
	a = strings.Split(s, " ")
	for i := range a {
		if strings.Contains(a[i], "__SPACE__") {
			a[i] = strings.ReplaceAll(a[i], "__SPACE__", " ")
		}
	}
	return a
}

func JoinArgs(arr []string, sep string) string {
	a := make([]string, 0)
	for _, s := range arr {
		if strings.ContainsRune(s, ' ') || strings.ContainsRune(s, '"') {
			s = `"` + strings.ReplaceAll(s, "\"", "\\\"") + `"`
		}
		a = append(a, s)
	}
	return strings.Join(a, sep)
}

func AppendUniqueString(to []string, from string) []string {
	found := false
	for _, str := range to {
		if str == from {
			found = true
			break
		}
	}
	if !found {
		to = append(to, from)
	}
	return to
}

func AppendUniqueStrings(to []string, from []string) []string {
	if from != nil {
		for _, fromStr := range from {
			to = AppendUniqueString(to, fromStr)
		}
	}
	return to
}

// 修复Golang中Json默认处理HTML转义 < > & 的问题
func FixJsonBytes(b []byte) []byte {
	l := len(b)
	i := 0
	for j := 0; j < l; j++ {
		if b[j] == '\\' && j < l-6 && b[j+1] == 'u' && b[j+2] == '0' && b[j+3] == '0' {
			// 替换
			var c byte = '0'
			if b[j+4] == '3' && b[j+5] == 'c' {
				c = '<'
			} else if b[j+4] == '3' && b[j+5] == 'e' {
				c = '>'
			} else if b[j+4] == '2' && b[j+5] == '6' {
				c = '&'
			}
			if c != '0' {
				b[i] = c
				j += 5
				i++
				continue
			}
		}

		// 复制
		if i != j {
			b[i] = b[j]
		}
		i++
	}

	if i != l {
		return b[0:i]
	} else {
		return b
	}
}

// 支持 map[interface{}]interface{}
func makeJsonType(inValue reflect.Value) *reflect.Value {
	if inValue.Kind() == reflect.Interface {
		inValue = inValue.Elem()
	}
	for inValue.Kind() == reflect.Ptr {
		inValue = inValue.Elem()
	}

	if !inValue.IsValid() {
		return nil
	}

	inType := inValue.Type()

	switch inType.Kind() {
	case reflect.Map:
		if inType.Key().Kind() == reflect.Interface {
			// 测试是否为数组
			isMap := false
			for i := len(inValue.MapKeys()); i > 0; i-- {
				if inValue.MapIndex(reflect.ValueOf(float64(i))).Kind() == reflect.Invalid {
					isMap = true
					break
				}
			}
			if isMap {
				// 处理字典
				newMap := reflect.MakeMap(reflect.MapOf(reflect.TypeOf(""), inType.Elem()))
				for _, k := range inValue.MapKeys() {
					v1 := inValue.MapIndex(k)
					v2 := makeJsonType(v1)
					k2 := reflect.ValueOf(StringFromValue(k))
					if v2 != nil {
						newMap.SetMapIndex(k2, *v2)
					} else {
						newMap.SetMapIndex(k2, v1)
					}
				}
				return &newMap
			} else {
				// 处理数组
				newArray := reflect.MakeSlice(reflect.SliceOf(inType.Elem()), inValue.Len(), inValue.Len())
				for i, k := range inValue.MapKeys() {
					v1 := inValue.MapIndex(k)
					v2 := makeJsonType(v1)
					if v2 != nil {
						newArray.Index(i).Set(*v2)
					} else {
						newArray.Index(i).Set(v1)
					}
				}
				return &newArray
			}
		} else {
			for _, k := range inValue.MapKeys() {
				v := makeJsonType(inValue.MapIndex(k))
				if v != nil {
					inValue.SetMapIndex(k, *v)
				}
			}
			return nil
		}
	case reflect.Slice:
		if inType.Elem().Kind() != reflect.Uint8 {
			for i := inValue.Len() - 1; i >= 0; i-- {
				v := makeJsonType(inValue.Index(i))
				if v != nil {
					inValue.Index(i).Set(*v)
				}
			}
		}
		return nil
	case reflect.Struct:
		for i := inType.NumField() - 1; i >= 0; i-- {
			f := inType.Field(i)
			if f.Anonymous {
				v := makeJsonType(inValue.Field(i))
				if v != nil {
					inValue.Field(i).Set(*v)
				}
			} else {
				if f.Name[0] <= 90 {
					v := makeJsonType(inValue.Field(i))
					if v != nil {
						inValue.Field(i).Set(*v)
					}
				}
			}
		}
		return nil
	default:
		return nil
	}
}

func JsonBytes(value interface{}) []byte {
	//fmt.Println("  &&&&&&&&&& 111", reflect.TypeOf(value))
	//if v1, ok := value.(map[interface{}]interface{}); ok {
	//	//fmt.Println("  &&&&&&&&&& 222")
	//	v2 := map[string]interface{}{}
	//	for k, v := range v1 {
	//		v2[String(k)] = v
	//	}
	//	value = v2
	//	//fmt.Println("  &&&&&&&&&& 333", reflect.TypeOf(value), "@@@@@@")
	//}
	if j, err := json.Marshal(value); err != nil {
		// 支持 interface{} 下标
		v2 := makeJsonType(reflect.ValueOf(value))
		var value2 interface{}
		if v2 != nil {
			value2 = makeJsonType(reflect.ValueOf(value)).Interface()
		} else {
			value2 = value
		}
		if r, err := json.Marshal(value2); err != nil {
			return []byte(fmt.Sprint(value))
		} else {
			return FixJsonBytes(r)
		}
	} else {
		return FixJsonBytes(j)
	}
}

func Json(value interface{}) string {
	return string(JsonBytes(value))
}

func FixedJson(value interface{}) string {
	buf := JsonBytes(value)
	excludeKeys := MakeExcludeUpperKeys(buf, "")
	bytesResult, err := json.Marshal(buf)
	if err != nil || (len(bytesResult) == 4 && string(bytesResult) == "null") {
		t := reflect.TypeOf(buf)
		if t.Kind() == reflect.Slice {
			bytesResult = []byte("[]")
		}
		if t.Kind() == reflect.Map {
			bytesResult = []byte("{}")
		}
	}
	FixUpperCase(buf, excludeKeys)
	return string(buf)
}

func JsonBytesP(value interface{}) []byte {
	j := JsonBytes(value)
	r := bytes.Buffer{}
	err := json.Indent(&r, j, "", "  ")
	if err == nil {
		return FixJsonBytes(r.Bytes())
	}
	return j
}

func JsonP(value interface{}) string {
	return string(JsonBytesP(value))
}

func FixedJsonP(value interface{}) string {
	buf := JsonBytes(value)
	excludeKeys := MakeExcludeUpperKeys(buf, "")
	bytesResult, err := json.Marshal(buf)
	if err != nil || (len(bytesResult) == 4 && string(bytesResult) == "null") {
		t := reflect.TypeOf(buf)
		if t.Kind() == reflect.Slice {
			bytesResult = []byte("[]")
		}
		if t.Kind() == reflect.Map {
			bytesResult = []byte("{}")
		}
	}
	FixUpperCase(buf, excludeKeys)
	r := bytes.Buffer{}
	json.Indent(&r, buf, "", "  ")
	return r.String()
}

func UnJsonBytes(data []byte, value interface{}) interface{} {
	if value == nil {
		var v interface{}
		value = &v
	}
	_ = json.Unmarshal(data, value)
	return value
}

func UnJson(str string, value interface{}) interface{} {
	return UnJsonBytes([]byte(str), value)
}

func UnJsonMap(str string) map[string]interface{} {
	value := map[string]interface{}{}
	UnJsonBytes([]byte(str), &value)
	return value
}

func UnJsonArr(str string) []interface{} {
	value := make([]interface{}, 0)
	UnJsonBytes([]byte(str), &value)
	return value
}

func Yaml(value interface{}) string {
	j, err := yaml.Marshal(value)
	if err == nil {
		return string(j)
	}
	return String(value)
}

func UnYaml(data string, value interface{}) interface{} {
	_ = yaml.Unmarshal([]byte(data), value)
	return value
}

func UnYamlMap(data string) map[string]interface{} {
	value := map[string]interface{}{}
	UnYaml(data, &value)
	return value
}

func UnYamlArr(data string) []interface{} {
	value := make([]interface{}, 0)
	UnYaml(data, &value)
	return value
}
