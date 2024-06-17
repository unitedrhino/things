package rulelogic

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoReadLogic {
	return &AlarmInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmInfoReadLogic) AlarmInfoRead(in *ud.WithID) (*ud.AlarmInfo, error) {
	po, err := relationDB.NewAlarmInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	return utils.Copy[ud.AlarmInfo](po), err
}
