package scene

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/robfig/cron/v3"
	"time"
)

// TimeRange 时间范围 只支持后面几种特殊字符:*  - ,
type TimeRange struct {
	Type string `json:"type"` //时间类型 cron
	Cron string `json:"cron"` //  cron表达式
}

func (t TimeRange) Validate() error {
	return nil
}

// Timer 定时器类型
type Timer struct {
	Type string //时间类型 cron
	Cron string //  cron表达式
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

func (t *Timer) Validate() error {
	if t == nil {
		return errors.Parameter.AddMsg("时间触发模式需要填写时间内容")
	}
	p := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	_, err := p.Parse(t.Cron)
	return errors.Parameter.AddMsg("时间cron表达式解析失败").AddDetail(err)
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
