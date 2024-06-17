package rulelogic

import (
	"context"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoDeleteLogic {
	return &AlarmInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmInfoDeleteLogic) AlarmInfoDelete(in *ud.WithID) (*ud.Empty, error) {
	err := relationDB.NewAlarmInfoRepo(l.ctx).Delete(l.ctx, in.Id)

	return &ud.Empty{}, err
}
