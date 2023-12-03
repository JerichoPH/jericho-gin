package tools

import (
	"fmt"
	"jericho-gin/types"
	"jericho-gin/wrongs"
	"time"
)

// Time 时间帮助函数
type Time struct {
	t time.Time
}

// NewTime 构造函数
func NewTime() Time { return Time{} }

// GetTime 获取时间
func (receiver Time) GetTime() time.Time {
	return receiver.t
}

// ToDateString 转换字符串（日期）
func (receiver Time) ToDateString() string {
	return receiver.t.Format(string(types.TIME_FORMAT_DATE))
}

// ToTimeString 转换字符串（时间）
func (receiver Time) ToTimeString() string {
	return receiver.t.Format(string(types.TIME_FORMAT_TIME))
}

// ToDateTimeString ToTimeString 转换字符串（日期时间）
func (receiver Time) ToDateTimeString() string {
	return receiver.t.Format(string(types.TIME_FORMAT_DATETIME))
}

// Parse 解析文字格式时间
func (receiver Time) Parse(format string, t string, more ...string) Time {
	if ti, err := time.Parse(format, fmt.Sprintf(t, more)); err != nil {
		wrongs.ThrowForbidden("时间格式解析错误(%s)：%s", t, err.Error())
	} else {
		receiver.t = ti
	}
	return receiver
}

// ParseDate 解析文字格式时间（日期）
func (receiver Time) ParseDate(t string, more ...string) Time {
	if ti, err := time.Parse(string(types.TIME_FORMAT_DATE), fmt.Sprintf(t, more)); err != nil {
		wrongs.ThrowForbidden("时间格式解析错误(%s)：%s", t, err.Error())
	} else {
		receiver.t = ti
	}
	return receiver
}

// ParseDatetime 解析文字格式时间（日期时间）
func (receiver Time) ParseDatetime(t string, more ...string) Time {
	if ti, err := time.Parse(string(types.TIME_FORMAT_DATETIME), fmt.Sprintf(t, more)); err != nil {
		wrongs.ThrowForbidden("时间格式解析错误(%s)：%s", t, err.Error())
	} else {
		receiver.t = ti
	}
	return receiver
}

// SetTime 设置时间
func (receiver Time) SetTime(t time.Time) Time {
	receiver.t = t
	return receiver
}

// SetTimeNowAdd8Hour 设置当前时间+8小时
func (receiver Time) SetTimeNowAdd8Hour() Time {
	receiver.t = time.Now().Add(8 * time.Hour)
	return receiver
}

// SetTimeNow 设置当前时间
func (receiver Time) SetTimeNow() Time {
	receiver.t = time.Now()
	return receiver
}

// SetTimeNowAdd 设置当前时间增加任意时间
func (receiver Time) SetTimeNowAdd(duration time.Duration) Time {
	receiver.t = time.Now().Add(duration)
	return receiver
}

// SetTimeYear 设置时间（年）
func (receiver Time) SetTimeYear(year int) Time {
	receiver.t = time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	return receiver
}

// SetTimeMonth 设置时间（年、月）
func (receiver Time) SetTimeMonth(year, month int) Time {
	receiver.t = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	return receiver
}

// SetTimeDay 设置时间（年、月、日）
func (receiver Time) SetTimeDay(year, month, day int) Time {
	receiver.t = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return receiver
}

// SetTimeHour 设置时间（年、月、日、时）
func (receiver Time) SetTimeHour(year, month, day, hour int) Time {
	receiver.t = time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.Local)
	return receiver
}

// SetTimeMinute 设置时间（年、月、日、时、分）
func (receiver Time) SetTimeMinute(year, month, day, hour, minute int) Time {
	receiver.t = time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.Local)
	return receiver
}

// SetTimeSecond 设置时间（年、月、日、时、分、秒）
func (receiver Time) SetTimeSecond(year, month, day, hour, minute, second int) Time {
	receiver.t = time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
	return receiver
}

// SetTimeNanoSecond 设置时间（年、月、日、时、分、秒、纳秒）
func (receiver Time) SetTimeNanoSecond(year, month, day, hour, minute, second, nanoSecond int) Time {
	receiver.t = time.Date(year, time.Month(month), day, hour, minute, second, nanoSecond, time.Local)
	return receiver
}

// Copy 复制时间对象
func (receiver Time) Copy() Time {
	return NewTime().SetTime(receiver.t)
}

// Add 增加时间
func (receiver Time) Add(duration time.Duration) Time {
	receiver.t = receiver.t.Add(duration)
	return receiver
}

// AddYears 增加年
func (receiver Time) AddYears(year int) Time {
	receiver.t = receiver.t.AddDate(year, 0, 0)
	return receiver
}

// AddMonths 增加月份
func (receiver Time) AddMonths(month int) Time {
	receiver.t = receiver.t.AddDate(0, month, 0)
	return receiver
}

// AddDays 增加日期
func (receiver Time) AddDays(day int) Time {
	receiver.t = receiver.t.AddDate(0, 0, day)
	return receiver
}

// StartOfYear 获取年起点
func (receiver Time) StartOfYear() Time {
	receiver.t = time.Date(receiver.t.Year(), 1, 1, 0, 0, 0, 0, receiver.t.Location())
	return receiver
}

// EndOfYear 获取年终点
func (receiver Time) EndOfYear() Time {
	receiver.t = time.Date(receiver.t.Year(), 12, 31, 23, 59, 59, 1000, receiver.t.Location())
	return receiver
}

// StartOfMonth 获取月起点
func (receiver Time) StartOfMonth() Time {
	receiver.t = receiver.StartOfDay().GetTime().AddDate(0, 0, -receiver.t.Day()+1)
	return receiver
}

// EndOfMonth 获取月终点
func (receiver Time) EndOfMonth() Time {
	receiver.t = NewTime().SetTime(receiver.StartOfMonth().GetTime().AddDate(0, 1, -1)).EndOfDay().GetTime()
	return receiver.EndOfDay()
}

// StartOfDay 获取日起点
func (receiver Time) StartOfDay() Time {
	receiver.t = time.Date(receiver.t.Year(), receiver.t.Month(), receiver.t.Day(), 0, 0, 0, 0, receiver.t.Location())
	return receiver
}

// EndOfDay 获取日终点
func (receiver Time) EndOfDay() Time {
	receiver.t = time.Date(receiver.t.Year(), receiver.t.Month(), receiver.t.Day(), 23, 59, 59, 1000, receiver.t.Location())
	return receiver
}

// StartOfWeek 获取周起点
func (receiver Time) StartOfWeek() Time {
	now := time.Now()

	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	receiver.t = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	return receiver
}

// EndOfWeek 获取周终点
func (receiver Time) EndOfWeek() Time {
	startOfWeek := receiver.StartOfWeek().ToDateString()
	if time, err := time.Parse(string(types.TIME_FORMAT_DATE), startOfWeek); err != nil {
		wrongs.ThrowForbidden("解析周起点时间失败 %", startOfWeek)
	} else {
		receiver.t = NewTime().SetTime(time.AddDate(0, 0, 7)).EndOfDay().GetTime()
	}
	return receiver
}

// GetCstZone 获取当前时区
func (receiver Time) GetCstZone() *time.Location {
	var cstZone = time.FixedZone("CST", 8*3600)
	return cstZone
}
