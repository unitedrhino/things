package alarmcenterlogic

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"github.com/i-Things/things/src/rulesvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmSceneDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AsDB *relationDB.AlarmSceneRepo
}

func NewAlarmSceneDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmSceneDeleteLogic {
	return &AlarmSceneDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AsDB:   relationDB.NewAlarmSceneRepo(ctx),
	}
}

func (l *AlarmSceneDeleteLogic) AlarmSceneDelete(in *rule.AlarmSceneDeleteReq) (*rule.Empty, error) {
	err := l.AsDB.DeleteByFilter(l.ctx, relationDB.AlarmSceneFilter{
		AlarmID: in.AlarmID,
		SceneID: in.SceneID,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &rule.Empty{}, nil
}
