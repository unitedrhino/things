package alarmcenterlogic

import (
	"context"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmSceneIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmSceneIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmSceneIndexLogic {
	return &AlarmSceneIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmSceneIndexLogic) AlarmSceneIndex(in *rule.AlarmSceneIndexReq) (*rule.AlarmSceneIndexResp, error) {
	// todo: add your logic here and delete this line

	return &rule.AlarmSceneIndexResp{}, nil
}
