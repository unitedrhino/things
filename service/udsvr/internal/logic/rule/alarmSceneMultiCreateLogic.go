package rulelogic

import (
	"context"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmSceneMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmSceneMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmSceneMultiCreateLogic {
	return &AlarmSceneMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 告警关联场景联动
func (l *AlarmSceneMultiCreateLogic) AlarmSceneMultiCreate(in *ud.AlarmSceneMultiSaveReq) (*ud.Empty, error) {
	var pos []*relationDB.UdAlarmScene
	for _, v := range in.SceneIDs {
		pos = append(pos, &relationDB.UdAlarmScene{SceneID: v, AlarmID: in.AlarmID})
	}
	err := relationDB.NewAlarmSceneRepo(l.ctx).MultiInsert(l.ctx, pos)
	return &ud.Empty{}, err
}
