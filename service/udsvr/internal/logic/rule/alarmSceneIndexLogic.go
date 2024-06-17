package rulelogic

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmSceneIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmSceneIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmSceneIndexLogic {
	return &AlarmSceneIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmSceneIndexLogic) AlarmSceneIndex(in *ud.AlarmSceneIndexReq) (*ud.AlarmSceneMultiSaveReq, error) {
	pos, err := relationDB.NewAlarmSceneRepo(l.ctx).FindByFilter(l.ctx, relationDB.AlarmSceneFilter{AlarmID: in.AlarmID}, nil)
	if err != nil {
		return nil, err
	}
	sceneIDs := utils.ToSliceWithFunc(pos, func(in *relationDB.UdAlarmScene) int64 {
		return in.SceneID
	})
	return &ud.AlarmSceneMultiSaveReq{SceneIDs: sceneIDs, AlarmID: in.AlarmID}, nil
}
