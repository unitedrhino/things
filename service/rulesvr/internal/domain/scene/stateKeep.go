package scene

import (
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
)

type StateKeepType string

const (
	StateKeepTypeDuration  = "duration"
	StateKeepTypeRepeating = "repeating"
)

// StateKeep 状态保持
type StateKeep struct {
	Type  StateKeepType `json:"type"`  //持续时间: duration  重复次数 repeating
	Value int64         `json:"value"` //持续的时间或重复的次数
	Unit  TimeUnit      `json:"unit"`  //时间单位 second:秒  minute:分钟  hour:小时  week:星期 month:月
}

func (s StateKeepType) Validate() error {
	if !utils.SliceIn(s, StateKeepTypeDuration, StateKeepTypeRepeating) {
		return errors.Parameter.AddMsgf("状态保持 类型不支持:%v", string(s))
	}
	return nil
}

func (s *StateKeep) Validate() error {
	if s == nil {
		return nil
	}
	if err := s.Type.Validate(); err != nil {
		return err
	}
	if err := s.Unit.Validate(); err != nil {
		return err
	}
	return nil
}
