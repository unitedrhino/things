package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/rulesvr/internal/domain/alarm"
	"github.com/i-Things/things/src/rulesvr/internal/repo/mysql"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmSceneMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmSceneMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmSceneMultiCreateLogic {
	return &AlarmSceneMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 告警关联场景联动
func (l *AlarmSceneMultiCreateLogic) AlarmSceneMultiCreate(in *rule.AlarmSceneMultiCreateReq) (*rule.Response, error) {
	//检查数据是否存在
	_, err := l.svcCtx.AlarmInfoRepo.FindOne(l.ctx, in.AlarmID)
	if err != nil {
		return nil, mysql.ToError(err)
	}
	err = l.svcCtx.AlarmSceneRepo.DeleteByFilter(l.ctx, alarm.SceneFilter{
		AlarmID: in.AlarmID,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	err = l.svcCtx.AlarmSceneRepo.InsertMulti(l.ctx, in.AlarmID, in.SceneID)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &rule.Response{}, nil
}
