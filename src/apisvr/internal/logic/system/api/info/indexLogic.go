package info

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
	info, err := l.svcCtx.ApiRpc.ApiInfoIndex(l.ctx, &sys.ApiInfoIndexReq{
		Page:   logic.ToSysPageRpc(req.Page),
		Route:  req.Route,
		Method: req.Method,
		Group:  req.Group,
		Name:   req.Name,
	})
	if err != nil {
		return nil, err
	}

	var total int64
	total = info.Total
	var apiInfo []*types.ApiInfo
	apiInfo = make([]*types.ApiInfo, 0, len(apiInfo))
	for _, i := range info.List {
		apiInfo = append(apiInfo, ToApiInfoTypes(i))
	}
	return &types.ApiInfoIndexResp{List: apiInfo, Total: total}, nil
}
