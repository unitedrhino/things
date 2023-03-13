package alarmcenterlogic

import (
	"context"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoDeleteLogic {
	return &AlarmInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmInfoDeleteLogic) AlarmInfoDelete(in *rule.AlarmInfoDeleteReq) (*rule.Response, error) {
	// todo: add your logic here and delete this line

	return &rule.Response{}, nil
}
