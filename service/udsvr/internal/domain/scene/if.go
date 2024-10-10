// Package scene 触发器
package scene

import (
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
)

type If struct {
	Triggers Triggers `json:"triggers,omitempty"`
}

type TriggerType = string

const (
	TriggerTypeDevice  TriggerType = "device"
	TriggerTypeTimer   TriggerType = "timer"
	TriggerTypeWeather TriggerType = "weather"
)

type Triggers []*Trigger
type Trigger struct {
	Type    TriggerType     `json:"type"`
	Order   int64           `json:"order"`
	AreaID  int64           `json:"areaID,string"`    //涉及到的区域ID
	Device  *TriggerDevice  `json:"device,omitempty"` //设备触发
	Timer   *TriggerTimer   `json:"timer,omitempty"`  //定时触发
	Weather *TriggerWeather `json:"weather,omitempty"`
}

func (t Triggers) Validate(repo CheckRepo) error {
	if len(t) == 0 {
		return nil
	}
	for _, v := range t {
		err := v.Validate(repo)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Trigger) Validate(repo CheckRepo) error {
	if t == nil {
		return errors.Parameter.AddMsg("需要填写触发内容")
	}
	if !utils.SliceIn(t.Type, TriggerTypeTimer, TriggerTypeDevice, TriggerTypeWeather) {
		return errors.Parameter.AddMsg("触发器不支持的类型:" + string(t.Type))
	}
	switch t.Type {
	case TriggerTypeDevice:
		return t.Device.Validate(repo)
	case TriggerTypeTimer:
		return t.Timer.Validate(repo)
	}
	return nil
}

func (i *If) Validate(t SceneType, repo CheckRepo) error {
	switch t {
	case SceneTypeManual:
		return nil
	case SceneTypeAuto:
		//if len(i.Triggers) == 0 {
		//	return errors.Parameter.AddMsg("自动触发类型需要填写至少一项设备触发或者定时触发")
		//}
		err := i.Triggers.Validate(repo)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.Parameter.AddMsg("触发类型只支持: auto:自动触发 manual:手动触发")
	}
}
