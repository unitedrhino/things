package rulelogic

import (
	"context"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/udsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/udsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/udsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/udsvr/pb/ud"

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
	f := relationDB.SceneLogFilter{SceneID: in.SceneID, Status: in.Status, WithSceneInfo: true,
		Time: logic.ToTimeRange(in.TimeRange)}
	pos, err := relationDB.NewSceneLogRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page).
		WithDefaultOrder(stores.OrderBy{
			Field: "createdTime",
			Sort:  2,
		}))
	if err != nil {
		return nil, err
	}
	total, err := relationDB.NewSceneLogRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	var list []*ud.SceneLog
	for _, v := range pos {
		one := utils.Copy[ud.SceneLog](v)
		if v.SceneInfo != nil {
			one.SceneName = v.SceneInfo.Name
		}
		list = append(list, one)
	}
	return &ud.SceneLogIndexResp{List: list, Total: total}, nil
}
