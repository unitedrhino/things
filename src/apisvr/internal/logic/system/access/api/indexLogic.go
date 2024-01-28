package api

import (
	"context"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

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

func (l *IndexLogic) Index(req *types.ApiInfoIndexReq) (resp *types.ApiInfoIndexResp, err error) {
	rst, err := l.svcCtx.AccessRpc.ApiInfoIndex(l.ctx, &sys.ApiInfoIndexReq{
		Page:         logic.ToSysPageRpc(req.Page),
		Route:        req.Route,
		Method:       req.Method,
		Name:         req.Name,
		AccessCode:   req.AccessCode,
		IsAuthTenant: req.IsAuthTenant,
	})
	if err != nil {
		return nil, err
	}
	return &types.ApiInfoIndexResp{
		List:  ToApiInfosTypes(rst.List),
		Total: rst.Total,
	}, nil
}
