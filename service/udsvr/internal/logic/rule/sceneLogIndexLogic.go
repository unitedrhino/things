package rulelogic

import (
	"context"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/udsvr/internal/logic"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type SceneLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSceneLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SceneLogIndexLogic {
	return &SceneLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SceneLogIndexLogic) SceneLogIndex(in *ud.SceneLogIndexReq) (*ud.SceneLogIndexResp, error) {
	f := relationDB.SceneLogFilter{SceneID: in.SceneID, Status: in.Status,
		Time: logic.ToTimeRange(in.TimeRange)}
	list, err := relationDB.NewSceneLogRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page).
		WithDefaultOrder(stores.OrderBy{
			Filed: "createdTime",
			Sort:  2,
		}))
	if err != nil {
		return nil, err
	}
	total, err := relationDB.NewSceneLogRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	return &ud.SceneLogIndexResp{List: utils.CopySlice[ud.SceneLog](list), Total: total}, nil
}
