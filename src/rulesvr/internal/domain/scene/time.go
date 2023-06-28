package scene

import (
	"github.com/i-Things/things/shared/crons"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/robfig/cron/v3"
	"time"
)

var secondParser = crons.NewParser(crons.Second | crons.Minute | crons.Hour | crons.Dom | crons.Month | crons.DowOptional | crons.Descriptor)

// TimeRange 时间范围 只支持后面几种特殊字符:*  - ,
type TimeRange struct {
	Type string `json:"type"` //时间类型 cron
	Cron string `json:"cron"` //  cron表达式
}

// Timer 定时器类型
type Timer struct {
	Type string `json:"type"` //时间类型 cron
	Cron string `json:"cron"` //  cron表达式
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
	_, err := secondParser.Parse(t.Cron)
	if err != nil {
		return errors.Parameter.AddMsgf("cron表达式检验不通过:%v", t.Cron).AddDetail(err)
	}
	return nil
}
func (t *TimeRange) IsHit(tim time.Time) bool {
	s, err := secondParser.Parse(t.Cron)
	if err != nil {
		return false
	}
	return s.Parse(tim)
}

func (t *Timer) Validate() error {
	if t == nil {
		return errors.Parameter.AddMsg("时间触发模式需要填写时间内容")
	}
	p := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	_, err := p.Parse(t.Cron)
	return errors.IfNotNil(errors.Parameter.AddMsg("时间cron表达式解析失败"), err)
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
