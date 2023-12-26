package scene

import (
	"github.com/i-Things/things/shared/crons"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
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
	if t.Exec < 0 || t.Exec > 24*60*60 {
		return errors.Parameter.AddMsg("时间执行时间范围只能在0到24小时之间")
	}
	if t.Repeat > 0b1111111 {
		return errors.Parameter.AddMsg("时间重复模式只能在0 7个二进制为高位")
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
