package scene

import (
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
)

// TriggerTimer 定时器类型
type TriggerTimer struct {
	ExecType      ExecType   `json:"execType"`      //执行方式
	ExecAt        int64      `json:"execAt"`        //执行的时间点 从0点加起来的秒数 如 1点就是 1*60*60
	ExecLoopStart int64      `json:"execLoopStart"` //循环执行起始时间配置
	ExecLoopEnd   int64      `json:"execLoopEnd"`
	ExecLoop      int64      `json:"execLoop"` //循环时间间隔
	RepeatType    RepeatType `json:"repeatType"`
	ExecRepeat    string     `json:"execRepeat"` //二进制周一到周日 11111111 或二进制月一到月末
}

type TriggerTimers []*TriggerTimer

func (t *TriggerTimer) Validate() error {
	if t == nil {
		return errors.Parameter.AddMsg("时间触发模式需要填写时间内容")
	}
	if t.ExecType == "" {
		t.ExecType = ExecTypeAt
	}
	if t.RepeatType == "" {
		t.RepeatType = RepeatTypeWeek
	}
	switch t.ExecType {
	case ExecTypeLoop:
		if t.ExecLoopStart < 0 || t.ExecLoopStart > 24*60*60 {
			return errors.Parameter.AddMsg("时间执行时间范围只能在0到24小时之间")
		}
		if t.ExecLoopEnd < 0 || t.ExecLoopEnd > 24*60*60 {
			return errors.Parameter.AddMsg("时间执行时间范围只能在0到24小时之间")
		}
		if t.ExecLoopEnd < t.ExecLoopStart {
			return errors.Parameter.AddMsg("时间执行时间范围只能在0到24小时之间")
		}
	default:
		if t.ExecAt < 0 || t.ExecAt > 24*60*60 {
			return errors.Parameter.AddMsg("时间执行时间范围只能在0到24小时之间")
		}
	}
	repeat := utils.BStrToInt64(t.ExecRepeat)
	switch t.RepeatType {
	case RepeatTypeWeek:
		if repeat > 0b1111111 {
			return errors.Parameter.AddMsg("时间重复模式只能在0 7个二进制为高位")
		}
	case RepeatTypeMount:
		if repeat > 0b1111111111111111111111111111111 {
			return errors.Parameter.AddMsg("时间重复模式只能在0 31个二进制为高位")
		}
	}
	return nil
}

func (t TriggerTimers) Validate() error {
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
