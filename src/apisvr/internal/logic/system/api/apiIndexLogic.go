package api

import (
	"context"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiIndexLogic {
	return &ApiIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiIndexLogic) ApiIndex(req *types.ApiIndexReq) (resp *types.ApiIndexResp, err error) {
	var page sys.PageInfo
	copier.Copy(&page, req.Page)
	info, err := l.svcCtx.ApiRpc.ApiIndex(l.ctx, &sys.ApiIndexReq{
		Page:   &page,
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
	var apiInfo []*types.ApiIndexData
	apiInfo = make([]*types.ApiIndexData, 0, len(apiInfo))
	for _, i := range info.List {
		apiInfo = append(apiInfo, &types.ApiIndexData{
			ID:           i.Id,
			Route:        i.Route,
			Method:       i.Method,
			Group:        i.Group,
			Name:         i.Name,
			BusinessType: i.BusinessType,
		})
	}
	return &types.ApiIndexResp{List: apiInfo, Total: total}, nil
}
