package scene

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type Then struct {
	Actions Actions `json:"actions"` //执行内容
}

func (t *Then) GetFlowPath() (ret []*FlowInfo) {
	if t == nil {
		return nil
	}
	for _, v := range t.Actions {
		ret = append(ret, v.GetFlowInfo())
	}
	return ret
}

func (t *Then) Validate(repo ValidateRepo) error {
	if t == nil || len(t.Actions) == 0 {
		return errors.Parameter.AddMsg("需要填写执行内容")
	}
	return t.Actions.Validate(repo)
}

func (t *Then) Execute(ctx context.Context, repo ActionRepo) error {
	repo.Info.Log = NewLog(repo.Info)
	for _, v := range t.Actions {
		err := v.Execute(ctx, repo)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return err
		}
	}
	err := repo.SaveLog(ctx, repo.Info.Log)
	if err != nil {
		logx.WithContext(ctx).Error(err)
	}
	return nil
}
