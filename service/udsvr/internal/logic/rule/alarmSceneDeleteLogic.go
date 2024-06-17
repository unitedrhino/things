package rulelogic

import (
	"context"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmSceneDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmSceneDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmSceneDeleteLogic {
	return &AlarmSceneDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmSceneDeleteLogic) AlarmSceneDelete(in *ud.AlarmSceneDeleteReq) (*ud.Empty, error) {
	err := relationDB.NewAlarmSceneRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.AlarmSceneFilter{
		AlarmID: in.AlarmID,
		SceneID: in.SceneID,
	})
	return &ud.Empty{}, err
}
