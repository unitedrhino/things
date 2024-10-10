package scene

import (
	"gitee.com/unitedrhino/share/errors"
)

// TriggerWeather 定时器类型
type TriggerWeather struct {
	TermWeather
}

type TriggerWeathers []*TriggerWeather

func (t *TriggerWeather) Validate(repo CheckRepo) error {
	if t == nil {
		return errors.Parameter.AddMsg("天气触发模式需要填写触发内容")
	}
	return t.Validate(repo)
}

func (t TriggerWeathers) Validate(repo CheckRepo) error {
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
