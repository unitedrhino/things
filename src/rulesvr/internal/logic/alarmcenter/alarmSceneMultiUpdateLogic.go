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

type AlarmSceneMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmSceneMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmSceneMultiUpdateLogic {
	return &AlarmSceneMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 告警关联场景联动
func (l *AlarmSceneMultiUpdateLogic) AlarmSceneMultiUpdate(in *rule.AlarmSceneMultiUpdateReq) (*rule.Response, error) {
	//检查数据是否存在
	_, err := l.svcCtx.AlarmInfoRepo.FindOne(l.ctx, in.AlarmID)
	if err != nil {
		return nil, mysql.ToError(err)
	}
	//先删除绑定的信息
	err = l.svcCtx.AlarmSceneRepo.DeleteByFilter(l.ctx, alarm.SceneFilter{AlarmID: in.AlarmID})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	err = l.svcCtx.AlarmSceneRepo.InsertMulti(l.ctx, in.AlarmID, in.SceneIDs)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &rule.Response{}, nil
}
