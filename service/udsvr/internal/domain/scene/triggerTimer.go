package scene

import (
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/observerly/dusk/pkg/dusk"
	"time"
)

// TriggerTimer 定时器类型
type TriggerTimer struct {
	ExecType      ExecType   `json:"execType"`                //执行方式
	ExecAdd       int64      `json:"execAdd,omitempty"`       //如果是日出日落模式,则为日出日落前后的秒数
	ExecAt        int64      `json:"execAt,omitempty"`        //执行的时间点 从0点加起来的秒数 如 1点就是 1*60*60
	ExecLoopStart int64      `json:"execLoopStart,omitempty"` //循环执行起始时间配置
	ExecLoopEnd   int64      `json:"execLoopEnd,omitempty"`
	ExecLoop      int64      `json:"execLoop,omitempty"` //循环时间间隔
	RepeatType    RepeatType `json:"repeatType,omitempty"`
	ExecRepeat    string     `json:"execRepeat,omitempty"` //二进制周一到周日 11111111 或二进制月一到月末
}

type TriggerTimers []*TriggerTimer

func (t *TriggerTimer) Validate(repo CheckRepo) error {
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
	case ExecTypeSunSet, ExecTypeSunRises:
		if t.ExecAt > 3*60*60 {
			return errors.Parameter.AddMsg("最晚只能三个小时后")
		}
		if t.ExecAt < (-3 * 60 * 60) {
			return errors.Parameter.AddMsg("最早只能提前三个小时")
		}
		pi, err := repo.ProjectCache.GetData(repo.Ctx, repo.Info.ProjectID)
		if err != nil {
			return err
		}
		if pi.Position == nil || pi.Position.Latitude == 0 || pi.Position.Longitude == 0 {
			return errors.Parameter.AddMsg("需要填写地理位置才可以使用日出日落触发")
		}
		twilight, _, err := dusk.GetLocalCivilTwilight(time.Now(), pi.Position.Longitude, pi.Position.Latitude, 0)
		if err != nil {
			return errors.System.AddMsg("计算日出日落时间失败").AddDetail(err)
		}
		switch t.ExecType {
		case ExecTypeSunRises:
			t.ExecAt = utils.TimeToDaySec(twilight.Until)
		case ExecTypeSunSet:
			t.ExecAt = utils.TimeToDaySec(twilight.From)
		}
		t.ExecAt += t.ExecAdd
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

func (t TriggerTimers) Validate(repo CheckRepo) error {
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
