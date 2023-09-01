package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/rulesvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmLogCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AlDB *relationDB.AlarmLogRepo
}

func NewAlarmLogCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmLogCreateLogic {
	return &AlarmLogCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AlDB:   relationDB.NewAlarmLogRepo(ctx),
	}
}

func (l *AlarmLogCreateLogic) AlarmLogCreate(in *rule.AlarmLog) (*rule.WithID, error) {
	err := l.AlDB.Insert(l.ctx, &relationDB.RuleAlarmLog{
		AlarmRecordID: in.AlarmRecordID,
		Serial:        in.Serial,
		SceneName:     in.SceneName,
		SceneID:       in.SceneID,
		Desc:          in.Desc,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &rule.WithID{}, nil
}
