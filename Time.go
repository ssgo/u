package u

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseTime(v interface{}) time.Time {
	if v == nil {
		return time.Now()
	}

	if tm, ok := v.(time.Time); ok {
		return tm
	}

	str := String(v)
	var tm time.Time
	var err error

	if num := Int64(v); num > 0 {
		// 20060102150405
		if len(str) == 14 {
			if tm, err = time.ParseInLocation("20060102150405", str, time.Local); err == nil {
				return tm
			}
		}
		// 20060102
		if len(str) == 8 {
			if tm, err = time.ParseInLocation("20060102", str, time.Local); err == nil {
				return tm
			}
		}
		// 150405
		if len(str) == 6 {
			if tm, err = time.ParseInLocation("150405", str, time.Local); err == nil {
				return tm
			}
		}

		switch {
		// 秒级时间戳
		case num < 1e10:
			return time.Unix(num, 0)
		// 毫秒级时间戳
		case num < 1e13:
			return time.UnixMilli(num)
		// 微秒级时间戳
		case num < 1e16:
			return time.UnixMicro(num)
		// 纳秒级时间戳
		default:
			return time.Unix(0, num)
		}
	}

	// 2006-01-02T15:04:05.999999999Z07:00
	if len(str) >= 35 && str[10] == 'T' && str[19] == '.' {
		if tm, err = time.Parse(time.RFC3339Nano, str); err == nil {
			return tm.In(time.Local)
		}
	}

	// 2006-01-02T15:04:05Z07:00
	if len(str) >= 25 && str[10] == 'T' && str[19] == 'Z' {
		if tm, err = time.Parse(time.RFC3339, str); err == nil {
			return tm.In(time.Local)
		}
	}

	// 2006-01-02 15:04:05.999999、2006-01-02T15:04:05.999999、2006/01/02 15:04:05.999999、2006/01/02T15:04:05.999999
	if len(str) >= 26 && (str[4] == '-' || str[4] == '/') && (str[10] == ' ' || str[10] == 'T') && str[19] == '.' {
		if tm, err = time.ParseInLocation(fmt.Sprintf("2006%c01%c02%c15:04:05.999999", str[4], str[4], str[10]), str[0:26], time.Local); err == nil {
			return tm
		}
	}
	// 01/02/2006 15:04:05.999999、01/02/2006T15:04:05.999999
	if len(str) >= 26 && str[2] == '/' && (str[10] == ' ' || str[10] == 'T') && str[19] == '.' {
		if tm, err = time.ParseInLocation(fmt.Sprintf("01/02/2006%c15:04:05.999999", str[10]), str[0:26], time.Local); err == nil {
			return tm
		}
	}
	// 2006-01-02 15:04:05.999、2006-01-02T15:04:05.999、2006/01/02 15:04:05.999、2006/01/02T15:04:05.999
	if len(str) >= 23 && (str[4] == '-' || str[4] == '/') && (str[10] == ' ' || str[10] == 'T') && str[19] == '.' {
		if tm, err = time.ParseInLocation(fmt.Sprintf("2006%c01%c02%c15:04:05.999999", str[4], str[4], str[10]), str[0:23], time.Local); err == nil {
			return tm
		}
	}
	// 01/02/2006 15:04:05.999、01/02/2006T15:04:05.999
	if len(str) >= 23 && str[2] == '/' && (str[10] == ' ' || str[10] == 'T') && str[19] == '.' {
		if tm, err = time.ParseInLocation(fmt.Sprintf("01/02/2006%c15:04:05.999", str[10]), str[0:23], time.Local); err == nil {
			return tm
		}
	}
	// 2006-01-02 15:04:05、2006-01-02T15:04:05、2006/01/02 15:04:05、2006/01/02T15:04:05
	if len(str) >= 19 && (str[4] == '-' || str[4] == '/') && (str[10] == ' ' || str[10] == 'T') {
		if tm, err = time.ParseInLocation(fmt.Sprintf("2006%c01%c02%c15:04:05", str[4], str[4], str[10]), str[0:19], time.Local); err == nil {
			return tm
		}
	}
	// 01/02/2006 15:04:05、01/02/2006T15:04:05
	if len(str) >= 19 && str[2] == '/' && (str[10] == ' ' || str[10] == 'T') {
		if tm, err = time.ParseInLocation(fmt.Sprintf("01/02/2006%c15:04:05", str[10]), str[0:19], time.Local); err == nil {
			return tm
		}
	}
	// 15:04:05.999999
	if len(str) >= 15 && str[2] == ':' && str[8] == '.' {
		if tm, err = time.ParseInLocation("15:04:05.999999", str[:15], time.Local); err == nil {
			return tm
		}
	}
	// 15:04:05.999
	if len(str) >= 12 && str[2] == ':' && str[8] == '.' {
		if tm, err = time.ParseInLocation("15:04:05.999", str[:12], time.Local); err == nil {
			return tm
		}
	}
	// 15:04:05
	if len(str) >= 8 && str[2] == ':' {
		if tm, err = time.ParseInLocation("15:04:05", str[:8], time.Local); err == nil {
			return tm
		}
	}
	// 2006-01-02、2006/01/02
	if len(str) >= 10 && (str[4] == '-' || str[4] == '/') {
		if tm, err = time.ParseInLocation(fmt.Sprintf("2006%c01%c02", str[4], str[4]), str[:10], time.Local); err == nil {
			return tm
		}
	}
	// 01/02/2006
	if len(str) >= 10 && str[2] == '/' {
		if tm, err = time.ParseInLocation("01/02/2006", str[0:10], time.Local); err == nil {
			return tm
		}
	}

	if len(str) > 24 && str[3] == ' ' {
		// Javascript：Mon Jun 23 2025 20:56:30 GMT+0800
		tzStr := "GMT"
		if strings.Contains(str, "MST") {
			// tzStr = "MST"
			str = strings.Replace(str, "MST", "-0700", 1)
		} else if strings.Contains(str, "CST") {
			// tzStr = "CST"
			str = strings.Replace(str, "CST", "+0800", 1)
		}
		if tm, err = time.Parse("Mon Jan 02 2006 15:04:05 "+tzStr+"-0700", str); err == nil {
			return tm.In(time.Local)
		}
		if tm, err = time.Parse("Mon Jan _2 2006 15:04:05 "+tzStr, str); err == nil {
			return tm.In(time.Local)
		}
		// UnixDate: Mon Jan _2 15:04:05 MST 2006
		if tm, err = time.Parse("Mon Jan _2 15:04:05 "+tzStr+" 2006", str); err == nil {
			return tm.In(time.Local)
		}
		// RubyDate: Mon Jan 02 15:04:05 -0700 2006
		if tm, err = time.Parse("Mon Jan 02 15:04:05 -0700 2006", str); err == nil {
			return tm.In(time.Local)
		}
		// ANSIC: Mon Jan _2 15:04:05 2006
		if tm, err = time.ParseInLocation("Mon Jan _2 15:04:05 2006", str, time.Local); err == nil {
			return tm
		}
	}

	// HTTP 头格式 (Mon, 02 Jan 2006 15:04:05 MST)
	if len(str) > 20 && str[3] == ',' {
		tzStr := "GMT"
		if strings.Contains(str, "MST") {
			// tzStr = "MST"
			str = strings.Replace(str, "MST", "-0700", 1)
		} else if strings.Contains(str, "CST") {
			// tzStr = "CST"
			str = strings.Replace(str, "CST", "+0800", 1)
		}
		if tm, err = time.Parse("Mon, 02 Jan 2006 15:04:05 "+tzStr, str); err == nil {
			return tm.In(time.Local)
		}
		if tm, err = time.Parse(time.RFC1123Z, str); err == nil {
			return tm.In(time.Local)
		}
	}

	return time.Now()
}

var dateFormatPattern = regexp.MustCompile(`[a-zA-Z]+`)

func FormatTime(layout string, timeValue interface{}) string {
	layout = dateFormatPattern.ReplaceAllStringFunc(layout, func(m string) string {
		switch m {
		case "YYYY":
			return "2006"
		case "YY":
			return "06"
		case "MM":
			return "01"
		case "M":
			return "1"
		case "DD":
			return "02"
		case "D":
			return "2"
		case "HH":
			return "15"
		case "H":
			return "15" // 注意：改为15，因为Go没有单数字24小时制
		case "hh":
			return "03"
		case "h":
			return "3"
		case "mm":
			return "04"
		case "ss":
			return "05"
		case "a":
			return "pm"
		case "A":
			return "PM"
		case "ZZ":
			return "-0700"
		case "Z":
			return "-07:00"
		default:
			return m
		}
	})
	tm := ParseTime(timeValue)
	return tm.Format(layout)
}

func AddTime(timeStr string, timeValue interface{}) time.Time {
	tm := ParseTime(timeValue)

	// 处理空字符串
	if timeStr == "" {
		return tm
	}

	i := 0
	years := 0
	months := 0
	days := 0
	duration := time.Duration(0)

	for i < len(timeStr) {
		// 处理每部分的符号（默认正数）
		sign := 1
		if timeStr[i] == '+' {
			i++
		} else if timeStr[i] == '-' {
			sign = -1
			i++
		}

		// 解析数字部分
		j := i
		for j < len(timeStr) && timeStr[j] >= '0' && timeStr[j] <= '9' {
			j++
		}
		numStr := ""
		if j == i { // 没有数字
			// return tm, fmt.Errorf("missing number at position %d", i)
			numStr = "1"
		} else {
			numStr = timeStr[i:j]
		}
		value, err := strconv.Atoi(numStr)
		if err != nil {
			// return tm, err
			value = 1
		}
		value *= sign

		// 解析单位部分
		i = j
		var unit string

		// 尝试匹配双字符单位
		if i+2 <= len(timeStr) {
			unit = timeStr[i : i+2]
			switch unit {
			case "ms", "us", "ns", "µs", "μs":
				i += 2
			default:
				unit = ""
			}
		}

		// 尝试匹配单字符单位
		if unit == "" && i < len(timeStr) {
			unit = timeStr[i : i+1]
			switch unit {
			case "Y", "M", "D", "h", "m", "s":
				i++
			default:
				// return tm, fmt.Errorf("unknown unit '%s' at position %d", unit, i)
				// 无单位时认为是秒
				unit = "s"
				i++
			}
		}

		// 单位处理
		switch unit {
		case "Y":
			years += value
		case "M":
			months += value
		case "D":
			days += value
		default:
			d, err := time.ParseDuration(fmt.Sprintf("%d%s", value, unit))
			if err == nil {
				duration += d
				// return tm, fmt.Errorf("invalid duration %d%s: %w", value, unit, err)
			}
		}
	}

	// 应用时间加减
	if years != 0 || months != 0 || days != 0 {
		tm = tm.AddDate(years, months, days)
	}
	if duration != 0 {
		tm = tm.Add(duration)
	}
	return tm
}
