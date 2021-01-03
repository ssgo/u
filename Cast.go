package u

import (
	"encoding/json"
	"fmt"
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
	if t != nil && (t.Kind() == reflect.Struct || t.Kind() == reflect.Map || t.Kind() == reflect.Slice) {
		return Json(value)
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
