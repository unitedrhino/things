package rulelogic

import (
	"context"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmRecordDealLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmRecordDealLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmRecordDealLogic {
	return &AlarmRecordDealLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmRecordDealLogic) AlarmRecordDeal(in *ud.AlarmRecordDealReq) (*ud.Empty, error) {
	po, err := relationDB.NewAlarmRecordRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	switch in.Handle {
	case 1:
		po.DealState = 2
		err := relationDB.NewAlarmRecordRepo(l.ctx).Update(l.ctx, po)
		if err != nil {
			return nil, err
		}
	case 2: //todo 添加到工作流
		po.DealState = 3
		err := relationDB.NewAlarmRecordRepo(l.ctx).Update(l.ctx, po)
		if err != nil {
			return nil, err
		}
	}

	return &ud.Empty{}, nil
}
