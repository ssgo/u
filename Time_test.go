package u_test

import (
	"testing"
	"time"

	"github.com/ssgo/u"
)

func TestParseTime(t *testing.T) {
	// 测试时间点（固定用于所有测试）
	refTime := time.Date(2025, 6, 23, 15, 30, 45, 123456789, time.Local)

	tests := []struct {
		input    interface{}
		expected time.Time
		name     string
	}{
		// 1. 数字时间戳
		{input: refTime.Unix(), expected: refTime.Truncate(time.Second), name: "秒级时间戳"},
		{input: refTime.UnixMilli(), expected: refTime.Truncate(time.Millisecond), name: "毫秒级时间戳"},
		{input: refTime.UnixMicro(), expected: refTime.Truncate(time.Microsecond), name: "微秒级时间戳"},
		{input: refTime.UnixNano(), expected: refTime, name: "纳秒级时间戳"},

		// 2. 纯数字格式
		{input: "20250623153045", expected: refTime.Truncate(time.Second), name: "YYYYMMDDHHmmSS格式"},
		{input: "20250623", expected: time.Date(2025, 6, 23, 0, 0, 0, 0, time.Local), name: "YYYYMMDD格式"},
		{input: "250623", expected: time.Date(2025, 6, 23, 0, 0, 0, 0, time.Local), name: "YYMMDD格式"},
		{input: "153045", expected: time.Date(0, 1, 1, 15, 30, 45, 0, time.Local), name: "HHmmSS格式"},

		// 3. RFC3339格式
		{input: "2025-06-23T15:30:45Z", expected: refTime.UTC().Truncate(time.Second), name: "RFC3339 UTC"},
		{input: "25-06-23T15:30:45Z", expected: refTime.UTC().Truncate(time.Second), name: "RFC3339 UTC"},
		{input: "2025-06-23T15:30:45.123Z", expected: refTime.UTC().Truncate(time.Millisecond), name: "RFC3339 MS UTC"},
		{input: "25-06-23T15:30:45.123Z", expected: refTime.UTC().Truncate(time.Millisecond), name: "RFC3339 MS UTC"},
		{input: "2025-06-23T15:30:45.1Z", expected: refTime.UTC().Truncate(time.Second).Add(100 * time.Millisecond), name: "RFC3339 MS/100 UTC"},
		{input: "25-06-23T15:30:45.1Z", expected: refTime.UTC().Truncate(time.Second).Add(100 * time.Millisecond), name: "RFC3339 MS/100 UTC"},
		{input: "2025-06-23T15:30:45+08:00", expected: refTime.Truncate(time.Second), name: "RFC3339 带时区"},
		{input: "2025.06.23T15:30:45+08:00", expected: refTime.Truncate(time.Second), name: "RFC3339 带时区"},
		{input: "25-06-23T15:30:45+08:00", expected: refTime.Truncate(time.Second), name: "RFC3339 带时区"},
		{input: "2025-06-23T15:30:45.123456789+08:00", expected: refTime, name: "RFC3339Nano"},
		{input: "25-06-23T15:30:45.123456789+08:00", expected: refTime, name: "RFC3339Nano"},

		// 4. JavaScript格式
		{input: "Mon Jun 23 2025 15:30:45 GMT+0800", expected: refTime.Truncate(time.Second), name: "JS格式1"},
		{input: "Mon Jun 23 2025 15:30:45 GMT-0700", expected: refTime.Add(15 * time.Hour).Truncate(time.Second), name: "JS格式2"},
		{input: "Mon, 23 Jun 2025 15:30:45 CST", expected: refTime.Truncate(time.Second), name: "RFC1123"},

		// 5. 常见日期时间格式
		{input: "2025-06-23 15:30:45", expected: refTime.Truncate(time.Second), name: "日期时间空格分隔"},
		{input: "2025-06-23 15:30", expected: refTime.Truncate(time.Minute), name: "日期时间空格分隔，无秒"},
		{input: "25-06-23 15:30", expected: refTime.Truncate(time.Minute), name: "日期时间空格分隔，无秒"},
		{input: "06-23 15:30", expected: time.Date(0, 6, 23, 15, 30, 0, 0, time.Local), name: "日期时间空格分隔，无秒"},
		{input: "25-06-23 15:30:45", expected: refTime.Truncate(time.Second), name: "日期时间空格分隔"},
		{input: "2025/06/23 15:30:45", expected: refTime.Truncate(time.Second), name: "斜杠分隔"},
		{input: "25/06/23 15:30:45", expected: refTime.Truncate(time.Second), name: "斜杠分隔"},
		{input: "06/23/2025 15:30:45", expected: refTime.Truncate(time.Second), name: "美式日期格式"},
		{input: "2025-06-23", expected: time.Date(2025, 6, 23, 0, 0, 0, 0, time.Local), name: "日期格式"},
		{input: "2025/06/23", expected: time.Date(2025, 6, 23, 0, 0, 0, 0, time.Local), name: "日期格式"},
		{input: "2025.06.23", expected: time.Date(2025, 6, 23, 0, 0, 0, 0, time.Local), name: "日期格式"},
		{input: "25-06-23", expected: time.Date(2025, 6, 23, 0, 0, 0, 0, time.Local), name: "日期格式"},
		{input: "25/06/23", expected: time.Date(2025, 6, 23, 0, 0, 0, 0, time.Local), name: "日期格式"},
		{input: "25.06.23", expected: time.Date(2025, 6, 23, 0, 0, 0, 0, time.Local), name: "日期格式"},
		{input: "15:30:45", expected: time.Date(0, 1, 1, 15, 30, 45, 0, time.Local), name: "时间格式"},

		// 6. 带小数秒的格式
		{input: "15:30:45.123", expected: time.Date(0, 1, 1, 15, 30, 45, 123000000, time.Local), name: "毫秒时间"},
		{input: "15:30:45.123456", expected: time.Date(0, 1, 1, 15, 30, 45, 123456000, time.Local), name: "微秒时间"},
		{input: "2025-06-23 15:30:45.123", expected: refTime.Truncate(time.Millisecond), name: "日期时间毫秒"},
		{input: "2025-06-23 15:30:45.123456", expected: refTime.Truncate(time.Microsecond), name: "日期时间毫秒"},
		{input: "2025.06.23 15:30:45.123456", expected: refTime.Truncate(time.Microsecond), name: "日期时间毫秒"},

		// 7. 边界和错误情况
		{input: "", expected: time.Now(), name: "空字符串"},
		{input: "invalid-time", expected: time.Now(), name: "无效格式"},
		{input: "Mon Jan 01 2024 00:00:00 GMT+0800 (中国标准时间)", expected: time.UnixMilli(1704038400000), name: "JS日期解析"},

		// 中文日期
		{input: "2025年06月23日 15点30分45秒", expected: refTime.Truncate(time.Second), name: "中文日期1"},
		{input: "2025年06月23日 15时30分45秒", expected: refTime.Truncate(time.Second), name: "中文日期2"},
		{input: "25年06月23日15时30分45秒", expected: refTime.Truncate(time.Second), name: "中文日期3"},
		{input: "2025年6月23日 下午3点30分45秒", expected: refTime.Truncate(time.Second), name: "中文日期4"},
		{input: "2025年6月23日下午3点30分", expected: refTime.Truncate(time.Minute), name: "中文日期5"},
		{input: "6月23日15点30分", expected: time.Date(0, 6, 23, 15, 30, 0, 0, time.Local), name: "中文日期4"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := u.ParseTime(tc.input)

			// 允许1秒内的误差（对于当前时间的情况）
			if tc.name == "空字符串" || tc.name == "无效格式" {
				if diff := time.Since(result); diff < -time.Second || diff > time.Second {
					t.Errorf("当前时间误差过大，期望误差<1秒，实际差异: %v", diff)
				}
				return
			}

			// 精确比较
			if !result.Equal(tc.expected) {
				t.Errorf("解析结果不匹配\n输入: %v\n期望: %v\n实际: %v", tc.input, tc.expected, result)
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	// 固定测试时间点：2006-01-02 15:04:05 UTC (星期一)
	baseTime := time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)

	tests := []struct {
		name   string
		layout string
		want   string
	}{
		// 基础格式
		{
			name:   "Full date/time",
			layout: "YYYY-MM-DD HH:mm:ss",
			want:   "2006-01-02 15:04:05",
		},
		// 年份
		{
			name:   "4-digit year",
			layout: "YYYY",
			want:   "2006",
		},
		{
			name:   "2-digit year",
			layout: "YY",
			want:   "06",
		},
		// 月份
		{
			name:   "2-digit month",
			layout: "MM",
			want:   "01",
		},
		{
			name:   "Single-digit month",
			layout: "M",
			want:   "1", // January without leading zero
		},
		// 日期
		{
			name:   "2-digit day",
			layout: "DD",
			want:   "02",
		},
		{
			name:   "Single-digit day",
			layout: "D",
			want:   "2", // Day 2 without leading zero
		},
		// 时间 (24小时制)
		{
			name:   "24-hour (2-digit)",
			layout: "HH",
			want:   "15",
		},
		{
			name:   "24-hour (no leading zero)",
			layout: "H",
			want:   "15", // Go的Format("15")固定返回两位数字
		},
		// 时间 (12小时制)
		{
			name:   "12-hour (2-digit)",
			layout: "hh",
			want:   "03", // 下午3点转12小时制=03
		},
		{
			name:   "12-hour (no leading zero)",
			layout: "h",
			want:   "3", // 无前导零的12小时制
		},
		// AM/PM标识
		{
			name:   "Lowercase am/pm",
			layout: "a",
			want:   "pm", // 15点=下午
		},
		{
			name:   "Uppercase AM/PM",
			layout: "A",
			want:   "PM",
		},
		// 时区
		{
			name:   "Timezone offset with colon",
			layout: "Z",
			want:   "+00:00", // UTC时区
		},
		{
			name:   "Timezone offset without colon",
			layout: "ZZ",
			want:   "+0000",
		},
		// 复杂组合格式
		{
			name:   "Mixed format with literals",
			layout: "Date: YYYY/M/D at hh:mm A Z",
			want:   "Date: 2006/1/2 at 03:04 PM +00:00",
		},
		// 无格式字符
		{
			name:   "Plain text",
			layout: "Hello World!",
			want:   "Hello World!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := u.FormatTime(tt.layout, baseTime)
			if got != tt.want {
				t.Errorf("FormatTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddTime(t *testing.T) {
	baseTime := time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)

	tests := []struct {
		name    string
		addExpr string
		want    time.Time
		wantErr bool
	}{
		// 基础加减
		{
			name:    "Add 1 day",
			addExpr: "+1D",
			want:    time.Date(2006, 1, 3, 15, 4, 5, 0, time.UTC),
		},
		{
			name:    "Subtract 2 days",
			addExpr: "-2D",
			want:    time.Date(2005, 12, 31, 15, 4, 5, 0, time.UTC),
		},
		// 年月组合
		{
			name:    "Add 1 year and 1 month",
			addExpr: "+1Y1M",
			want:    time.Date(2007, 2, 2, 15, 4, 5, 0, time.UTC),
		},
		{
			name:    "Subtract 3 months",
			addExpr: "-3M",
			want:    time.Date(2005, 10, 2, 15, 4, 5, 0, time.UTC),
		},
		// 时间单位
		{
			name:    "Add 2 hours",
			addExpr: "+2h",
			want:    time.Date(2006, 1, 2, 17, 4, 5, 0, time.UTC),
		},
		{
			name:    "Add 90 minutes",
			addExpr: "+90m",
			want:    time.Date(2006, 1, 2, 16, 34, 5, 0, time.UTC),
		},
		{
			name:    "Add 30 seconds",
			addExpr: "+30s",
			want:    time.Date(2006, 1, 2, 15, 4, 35, 0, time.UTC),
		},
		// 微秒级精度
		{
			name:    "Add 500 milliseconds",
			addExpr: "+500ms",
			want:    time.Date(2006, 1, 2, 15, 4, 5, 5e8, time.UTC),
		},
		// 混合操作
		{
			name:    "Complex: +1Y-2M+3D-4h+5m",
			addExpr: "+1Y-2M+3D-4h+5m",
			want:    time.Date(2006, 11, 5, 11, 9, 5, 0, time.UTC), // 2006-01-02 +1年=2007, -2月=2006-11-02, +3天=2006-11-05, -4小时=11:04, +5分=11:09
		},
		// 边界测试
		{
			name:    "End of month rollover",
			addExpr: "+31D", // 2006-01-02 +31天
			want:    time.Date(2006, 2, 2, 15, 4, 5, 0, time.UTC),
		},
		{
			name:    "Leap year transition",
			addExpr: "+365D", // 2006不是闰年
			want:    time.Date(2007, 1, 2, 15, 4, 5, 0, time.UTC),
		},
		// 错误处理
		{
			name:    "No number",
			addExpr: "Y", // 缺少数字
			wantErr: true,
		},
		{
			name:    "Invalid unit",
			addExpr: "+1X", // 无效单位X
			wantErr: true,
		},
		{
			name:    "Empty expression",
			addExpr: "",
			want:    baseTime,
		},
		{
			name:    "Only sign",
			addExpr: "+",
			want:    baseTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := u.AddTime(tt.addExpr, baseTime)
			if !got.Equal(tt.want) && !tt.wantErr {
				t.Errorf("AddTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
