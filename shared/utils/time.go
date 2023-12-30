package utils

import "time"

// GetMonthDays 获取指定年月的天数
func GetMonthDays(years int, month time.Month) int {
	startTime := time.Date(years, month, 1, 0, 0, 0, 0, time.Local)
	endTime := startTime.AddDate(0, 1, 0)
	return endTime.Add(-time.Second).Day()
}

// GetSubDay 比较两个时间戳差的天数
func GetSubDay(base time.Time, cmp time.Time) int64 {
	return (cmp.Unix() - base.Unix()) / int64(24*time.Hour/time.Second)
}

// 转换为07-02 01:02 这种格式
func ToMMddHHSS(timeStamp int64) string {
	return time.UnixMilli(timeStamp).Format("01-02 15:04")
}

func ToddHHSS(timeStamp int64) string {
	return time.UnixMilli(timeStamp).Format("15:04:03")
}

// 转换为07-02 01:02 这种格式
func ToYYMMddHHSS(timeStamp int64) string {
	return time.UnixMilli(timeStamp).Format("2006-01-02 15:04")
}

// 转换为07-02 01:02 这种格式
func ToYYMMdd(timeStamp int64) string {
	return time.UnixMilli(timeStamp).Format("2006 01-02")
}

// 转换为07-02 01:02 这种格式
func ToYYMMdd2(timeStamp int64) string {
	return time.UnixMilli(timeStamp).Format("2006-01-02")
}
