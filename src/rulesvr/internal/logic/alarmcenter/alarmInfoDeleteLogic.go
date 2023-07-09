package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/rulesvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.AlarmInfoRepo
}

func NewAlarmInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoDeleteLogic {
	return &AlarmInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewAlarmInfoRepo(ctx),
	}
}

func (l *AlarmInfoDeleteLogic) AlarmInfoDelete(in *rule.WithID) (*rule.Empty, error) {
	err := l.AiDB.Delete(l.ctx, in.Id)
	//todo 要把日志等删除
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &rule.Empty{}, nil
}
