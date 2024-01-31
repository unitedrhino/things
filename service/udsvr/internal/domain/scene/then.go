package scene

import "gitee.com/i-Things/core/shared/errors"

type Then struct {
	Actions Actions `json:"actions"` //执行内容
}

func (t *Then) Validate() error {
	if t == nil || len(t.Actions) == 0 {
		return errors.Parameter.AddMsg("需要填写执行内容")
	}
	return t.Actions.Validate()
}
