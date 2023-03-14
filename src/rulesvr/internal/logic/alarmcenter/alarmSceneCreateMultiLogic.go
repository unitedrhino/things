package alarmcenterlogic

import (
	"context"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmSceneCreateMultiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmSceneCreateMultiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmSceneCreateMultiLogic {
	return &AlarmSceneCreateMultiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 告警关联场景联动
func (l *AlarmSceneCreateMultiLogic) AlarmSceneCreateMulti(in *rule.AlarmSceneCreateMultiReq) (*rule.Response, error) {
	// todo: add your logic here and delete this line

	return &rule.Response{}, nil
}
