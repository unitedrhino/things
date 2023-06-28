package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/rulesvr/internal/domain/alarm"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmSceneDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmSceneDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmSceneDeleteLogic {
	return &AlarmSceneDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmSceneDeleteLogic) AlarmSceneDelete(in *rule.AlarmSceneDeleteReq) (*rule.Empty, error) {
	err := l.svcCtx.AlarmSceneRepo.DeleteByFilter(l.ctx, alarm.SceneFilter{
		AlarmID: in.AlarmID,
		SceneID: in.SceneID,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &rule.Empty{}, nil
}
