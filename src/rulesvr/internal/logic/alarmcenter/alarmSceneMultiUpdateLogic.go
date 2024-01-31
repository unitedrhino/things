package alarmcenterlogic

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"github.com/i-Things/things/src/rulesvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmSceneMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.AlarmInfoRepo
	AsDB *relationDB.AlarmSceneRepo
}

func NewAlarmSceneMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmSceneMultiUpdateLogic {
	return &AlarmSceneMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewAlarmInfoRepo(ctx),
		AsDB:   relationDB.NewAlarmSceneRepo(ctx),
	}
}

// 告警关联场景联动
func (l *AlarmSceneMultiUpdateLogic) AlarmSceneMultiUpdate(in *rule.AlarmSceneMultiUpdateReq) (*rule.Empty, error) {
	//检查数据是否存在
	_, err := l.AiDB.FindOne(l.ctx, in.AlarmID)
	if err != nil {
		return nil, err
	}
	//先删除绑定的信息
	err = l.AsDB.DeleteByFilter(l.ctx, relationDB.AlarmSceneFilter{AlarmID: in.AlarmID})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	err = l.AsDB.MultiInsert(l.ctx, in.AlarmID, in.SceneIDs)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &rule.Empty{}, nil
}
