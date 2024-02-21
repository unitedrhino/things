package scene

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type Then struct {
	Actions Actions `json:"actions"` //执行内容
}

func (t *Then) Validate() error {
	if t == nil || len(t.Actions) == 0 {
		return errors.Parameter.AddMsg("需要填写执行内容")
	}
	return t.Actions.Validate()
}
func (t *Then) Execute(ctx context.Context, repo ActionRepo) error {
	for _, v := range t.Actions {
		err := v.Execute(ctx, repo)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return err
		}
	}
	return nil
}
