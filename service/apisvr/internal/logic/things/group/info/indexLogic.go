package info

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.GroupInfoIndexReq) (resp *types.GroupInfoIndexResp, err error) {
	res, err := l.svcCtx.DeviceG.GroupInfoIndex(l.ctx, &dm.GroupInfoIndexReq{
		Page:     logic.ToDmPageRpc(req.Page),
		ParentID: req.ParentID,
		AreaID:   req.AreaID,
		Name:     req.Name,
		Tags:     logic.ToTagsMap(req.Tags),
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.GroupInfoIndex req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	glist := make([]*types.GroupInfo, 0, len(res.List))
	for _, v := range res.List {
		glist = append(glist, ToGroupInfoTypes(v))
	}

	return &types.GroupInfoIndexResp{
		List:  glist,
		Total: res.Total,
	}, nil
}
