package scene

import (
	"github.com/i-Things/things/shared/crons"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"time"
)

var secondParser = crons.NewParser(crons.Second | crons.Minute | crons.Hour | crons.Dom | crons.Month | crons.DowOptional | crons.Descriptor)

const (
	TimeRangeTypeAllDay = "allDay"
	TimeRangeTypeLight  = "light"
	TimeRangeTypeNight  = "night"
	TimeRangeTypeCustom = "custom"
)
const (
	DateRangeTypeWorkDay = "workday"
	DateRangeTypeWeekend = "weekend"
	DateRangeTypeHoliday = "holiday"
	DateRangeTypeCustom  = "custom"
)

// TimeRange 时间范围 只支持后面几种特殊字符:*  - ,
type TimeRange struct {
	Type      string `json:"type"`      //时间类型  allDay:全天 light:白天(从日出到日落) night:夜间(从日落到日出) custom:自定义
	StartTime int64  `json:"startTime"` //自定义开始时间 从0点加起来的秒数
	EndTime   int64  `json:"endTime"`   //自定义结束时间 从0点加起来的秒数
}

type DateRange struct {
	Type      string `json:"type"`      //日期类型 workday: 工作日 weekend: 周末 holiday: 节假日 custom:自定义
	StartDate string `json:"startDate"` //开始日期 2006-01-02
	EndDate   string `json:"endDate"`   //结束日期 2006-01-02
}

type Timers []*Timer

// Timer 定时器类型
type Timer struct {
	Exec   int64 `json:"exec"`   //从0点加起来的秒数 如 1点就是 1*60*60
	Repeat int64 `json:"repeat"` //二进制周一到周日 11111111
}

type TimeUnit string

const (
	TimeUnitSeconds TimeUnit = "seconds"
	TimeUnitMinutes TimeUnit = "minutes"
	TimeUnitHours   TimeUnit = "hours"
)

type UnitTime struct {
	Time int64    `json:"time"` //延迟时间
	Unit TimeUnit `json:"unit"` //时间单位 second:秒  minute:分钟  hour:小时  week:星期 month:月
}

func (t *TimeRange) Validate() error {
	if t == nil {
		return errors.Parameter.AddMsg("时间范围需要填写时间内容")
	}
	if !utils.SliceIn(t.Type, TimeRangeTypeAllDay, TimeRangeTypeLight, TimeRangeTypeNight, TimeRangeTypeCustom) {
		return errors.Parameter.AddMsg("时间范围类型不正确")
	}
	if t.Type == TimeRangeTypeCustom {
		if t.StartTime < 0 || t.StartTime > 24*60*60 || t.EndTime < 0 || t.EndTime > 24*60*60 || t.StartTime > t.EndTime {
			return errors.Parameter.AddMsg("自定义时间范围只能在0到24小时之间")
		}
	}
	return nil
}

func (t *DateRange) Validate() error {
	if t == nil {
		return errors.Parameter.AddMsg("日期范围需要填写日期内容")
	}
	if !utils.SliceIn(t.Type, DateRangeTypeWorkDay, DateRangeTypeWeekend, DateRangeTypeHoliday, DateRangeTypeCustom) {
		return errors.Parameter.AddMsg("日期范围类型不正确")
	}
	if t.Type == DateRangeTypeCustom {
		start := utils.FmtNilDateStr(t.StartDate)
		if start == nil {
			return errors.Parameter.AddMsg("日期范围开始时间的格式为:2006-01-02")
		}
		end := utils.FmtNilDateStr(t.EndDate)
		if end == nil {
			return errors.Parameter.AddMsg("日期范围结束时间的格式为:2006-01-02")
		}
	}
	return nil
}
func (t *TimeRange) IsHit(tim time.Time) bool {
	//todo 等待实现
	return true
}

func (t *Timer) Validate() error {
	if t == nil {
		return errors.Parameter.AddMsg("时间触发模式需要填写时间内容")
	}
	if t.Exec < 0 || t.Exec > 24*60*60 {
		return errors.Parameter.AddMsg("时间执行时间范围只能在0到24小时之间")
	}
	if t.Repeat > 0b1111111 {
		return errors.Parameter.AddMsg("时间重复模式只能在0 7个二进制为高位")
	}
	return nil
}

func (t Timers) Validate() error {
	if len(t) == 0 {
		return nil
	}
	for _, v := range t {
		err := v.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (t TimeUnit) Validate() error {
	if !utils.SliceIn(t, TimeUnitSeconds, TimeUnitMinutes, TimeUnitHours) {
		return errors.Parameter.AddMsg("时间单位不支持:" + string(t))
	}
	return nil
}

func (a *UnitTime) Validate() error {
	if a == nil {
		return nil
	}
	return a.Unit.Validate()
}

func (a *UnitTime) Execute() {
	var delayTime = time.Duration(a.Time)
	switch a.Unit {
	case TimeUnitSeconds:
		delayTime *= time.Second
	case TimeUnitMinutes:
		delayTime *= time.Minute
	case TimeUnitHours:
		delayTime *= time.Hour
	}
	time.Sleep(delayTime)
}
