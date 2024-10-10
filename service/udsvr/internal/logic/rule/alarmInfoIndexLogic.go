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

type AlarmInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoIndexLogic {
	return &AlarmInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmInfoIndexLogic) AlarmInfoIndex(in *ud.AlarmInfoIndexReq) (*ud.AlarmInfoIndexResp, error) {
	f := relationDB.AlarmInfoFilter{Name: in.Name}
	pos, err := relationDB.NewAlarmInfoRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page).
		WithDefaultOrder(stores.OrderBy{
			Field: "createdTime",
			Sort:  stores.OrderDesc,
		}))
	if err != nil {
		return nil, err
	}
	total, err := relationDB.NewAlarmInfoRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	var list []*ud.AlarmInfo
	for _, po := range pos {
		v := utils.Copy[ud.AlarmInfo](po)
		for _, s := range po.Scenes {
			v.SceneIDs = append(v.SceneIDs, s.SceneID)
		}
		list = append(list, v)
	}
	return &ud.AlarmInfoIndexResp{List: list, Total: total}, nil
}
