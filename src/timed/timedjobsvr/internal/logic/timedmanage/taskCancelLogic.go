package timedmanagelogic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/src/timed/internal/domain"
	"github.com/i-Things/things/src/timed/timedjobsvr/internal/svc"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskCancelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskCancelLogic {
	return &TaskCancelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var (
	retryKeys []string
)

func init() {
	for _, priority := range domain.Prioritys {
		retryKeys = append(retryKeys, fmt.Sprintf("asynq:{%s}:retry", priority))
	}
}

func (l *TaskCancelLogic) TaskCancel(in *timedjob.TaskWithTaskID) (*timedjob.Response, error) {
	var hashKeys []string
	for _, priority := range domain.Prioritys {
		hashKeys = append(hashKeys, fmt.Sprintf("asynq:{%s}:t:%s", priority, in.TaskID))
	}
	_, err := l.svcCtx.Store.Del(hashKeys...)
	if err != nil {
		return nil, err
	}
	for _, v := range retryKeys {
		l.svcCtx.Store.ZremCtx(l.ctx, v, in.TaskID)
	}
	return &timedjob.Response{}, err
}
