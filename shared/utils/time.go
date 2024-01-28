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

// GetFirstDateOfMonth 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

// GetLastDateOfMonth 获取传入的时间所在月份的最后一天，即某月最后一天的23:59:59。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	d = GetEndTime(d)
	return d.AddDate(0, 1, -1)
}

// GetZeroTime 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// GetEndTime 获取某一天的结束点时间
func GetEndTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location())
}

// UnixSecondToTime 秒级时间戳转time
func UnixSecondToTime(second int64) time.Time {
	return time.Unix(second, 0)
}

// UnixMilliToTime 毫秒级时间戳转time
func UnixMilliToTime(milli int64) time.Time {
	return time.Unix(milli/1000, (milli%1000)*(1000*1000))
}

// UnixNanoToTime 纳秒级时间戳转time
func UnixNanoToTime(nano int64) time.Time {
	return time.Unix(nano/(1000*1000*1000), nano%(1000*1000*1000))
}

func FmtDateStr(t string) time.Time {
	if t == "" {
		return time.Now()
	}
	ret, err := time.ParseInLocation("2006-01-02", t, time.Local)
	if err != nil {
		ret, err = time.ParseInLocation("2006-01", t, time.Local)
		if err != nil {
			return time.Now()
		}
	}
	return ret
}

func FmtNilDateStr(t string) *time.Time {
	if t == "" {
		return nil
	}
	ret, err := time.ParseInLocation("2006-01-02", t, time.Local)
	if err != nil {
		return nil
	}
	return &ret
}

func TimeInt64ToStr(t int64) string {
	return ToDateStr(time.Unix(t, 0))
}

func ToDateStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}
func ToTimeStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// TimeToUnix 时间转成时间戳
func TimeToUnix(e time.Time) int64 {
	timeUnix, _ := time.Parse("2006-01-02 15:04:05", e.Format("2006-01-02 15:04:05"))
	return timeUnix.UnixNano() / 1e6
}

// GetDiffDays 计算日期相差的天数
func GetDiffDays(t1 time.Time, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)

	return int(t1.Sub(t2).Hours() / 24)
}
