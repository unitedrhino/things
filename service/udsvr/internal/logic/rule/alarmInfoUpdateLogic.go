package rulelogic

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoUpdateLogic {
	return &AlarmInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmInfoUpdateLogic) AlarmInfoUpdate(in *ud.AlarmInfo) (*ud.Empty, error) {
	old, err := relationDB.NewAlarmInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.Name != "" {
		old.Name = in.Name
	}
	if in.Desc != "" {
		old.Desc = in.Desc
	}
	if in.Level != 0 {
		old.Level = in.Level
	}
	if in.Status != 0 {
		old.Status = in.Status
	}
	if len(in.Notifies) != 0 {
		old.Notifies = utils.CopySlice[relationDB.UdAlarmNotify](in.Notifies)
	}
	err = relationDB.NewAlarmInfoRepo(l.ctx).Update(l.ctx, old)
	return &ud.Empty{}, err
}
