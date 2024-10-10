package rulelogic

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/udsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/udsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/udsvr/pb/ud"

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

func (l *AlarmSceneIndexLogic) AlarmSceneIndex(in *ud.AlarmSceneIndexReq) (*ud.AlarmSceneIndexResp, error) {
	pos, err := relationDB.NewAlarmSceneRepo(l.ctx).FindByFilter(l.ctx, relationDB.AlarmSceneFilter{AlarmID: in.AlarmID}, nil)
	if err != nil {
		return nil, err
	}
	sceneIDs := utils.ToSliceWithFunc(pos, func(in *relationDB.UdAlarmScene) int64 {
		return in.SceneID
	})
	list, err := relationDB.NewSceneInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.SceneInfoFilter{IDs: sceneIDs}, nil)
	if err != nil {
		return nil, err
	}
	return &ud.AlarmSceneIndexResp{Scenes: PoToSceneInfoPbs(l.ctx, l.svcCtx, list), AlarmID: in.AlarmID}, nil
}
