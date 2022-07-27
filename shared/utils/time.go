package utils

import "time"

// GetMonthDays 获取指定年月的天数
func GetMonthDays(years int, month time.Month) int {
	startTime := time.Date(years, month, 1, 0, 0, 0, 0, time.Local)
	endTime := startTime.AddDate(0, 1, 0)
	return endTime.Add(-time.Second).Day()
}
