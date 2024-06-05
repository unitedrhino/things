package scene

import (
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
)

type StateKeepType = string

const (
	StateKeepTypeDuration  = "duration"
	StateKeepTypeRepeating = "repeat"
)

// StateKeep 状态保持
type StateKeep struct {
	Type  StateKeepType `json:"type"`  //持续时间: duration  重复次数 repeating
	Value int64         `json:"value"` //持续的时间(秒)或重复的次数
}

func (s *StateKeep) Validate() error {
	if s == nil {
		return nil
	}
	if !utils.SliceIn(s.Type, StateKeepTypeDuration, StateKeepTypeRepeating) {
		return errors.Parameter.AddMsg("状态保持 类型不支持:" + string(s.Type))
	}
	return nil
}
