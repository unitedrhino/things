package rulelogic

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/udsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/udsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/udsvr/pb/ud"

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
	v := utils.Copy[ud.AlarmInfo](po)
	for _, s := range po.Scenes {
		v.SceneIDs = append(v.SceneIDs, s.SceneID)
	}
	return v, err
}
