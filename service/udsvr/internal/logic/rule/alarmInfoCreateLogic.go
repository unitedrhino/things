package rulelogic

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoCreateLogic {
	return &AlarmInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmInfoCreateLogic) AlarmInfoCreate(in *ud.AlarmInfo) (*ud.WithID, error) {
	po := utils.Copy[relationDB.UdAlarmInfo](in)
	po.ID = 0
	err := relationDB.NewAlarmInfoRepo(l.ctx).Insert(l.ctx, po)
	if err != nil {
		return nil, err
	}
	return &ud.WithID{Id: po.ID}, nil
}
