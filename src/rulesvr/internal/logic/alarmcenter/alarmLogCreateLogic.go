package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/rulesvr/internal/repo/mysql"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmLogCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmLogCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmLogCreateLogic {
	return &AlarmLogCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmLogCreateLogic) AlarmLogCreate(in *rule.AlarmLog) (*rule.Response, error) {
	_, err := l.svcCtx.AlarmInfoRepo.FindOne(l.ctx, in.AlarmID)
	if !(err == mysql.ErrNotFound) {
		return nil, errors.Parameter.AddMsg("告警名称重复").AddDetail(err)
	}
	_, err = l.svcCtx.AlarmLogRepo.Insert(l.ctx, &mysql.RuleAlarmLog{
		AlarmID:   in.AlarmID,
		Serial:    in.Serial,
		SceneName: in.SceneName,
		SceneID:   in.SceneID,
		Desc:      in.Desc,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &rule.Response{}, nil
}
