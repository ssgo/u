package u

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
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
	if value == nil {
		return ""
	}
	value = FixPtr(value)
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
		j, err := json.Marshal(value)
		if err == nil {
			return string(FixJsonBytes(j))
		}
		return fmt.Sprint(value)
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

//func FixMap(interface{}) interface{} {
//
//}

func JsonBytes(value interface{}) []byte {
	//fmt.Println("  &&&&&&&&&& 111", reflect.TypeOf(value))
	if v1, ok := value.(map[interface{}]interface{}); ok {
		//fmt.Println("  &&&&&&&&&& 222")
		v2 := map[string]interface{}{}
		for k, v := range v1 {
			v2[String(k)] = v
		}
		value = v2
		//fmt.Println("  &&&&&&&&&& 333", reflect.TypeOf(value), "@@@@@@")
	}
	j, err := json.Marshal(value)
	if err == nil {
		//fmt.Println("  &&&&&&&&&& 444", string(j), "@@@@@@")
		return FixJsonBytes(j)
	} else {
		//fmt.Println("  &&&&&&&&&& 555", String(value), "@@@@@@")
		//fmt.Println("error", err.Error())
		return Bytes(value)
	}
}

func Json(value interface{}) string {
	return string(JsonBytes(value))
}

func FixedJson(value interface{}) string {
	buf := JsonBytes(value)
	FixUpperCase(buf, nil)
	return string(buf)
}

func JsonBytesP(value interface{}) []byte {
	j, err := json.MarshalIndent(value, "", "  ")
	if err == nil {
		return FixJsonBytes(j)
	}
	return Bytes(value)
}

func JsonP(value interface{}) string {
	return string(JsonBytesP(value))
}

func FixedJsonP(value interface{}) string {
	buf := JsonBytesP(value)
	FixUpperCase(buf, nil)
	return string(buf)
}

func UnJsonBytes(data []byte, value interface{}) interface{} {
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
