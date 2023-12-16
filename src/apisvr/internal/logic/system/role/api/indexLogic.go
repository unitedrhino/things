package api

import (
	"context"
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

func (l *IndexLogic) Index(req *types.RoleApiIndexReq) (resp *types.RoleApiIndexResp, err error) {
	info, err := l.svcCtx.RoleRpc.RoleApiIndex(l.ctx, &sys.RoleApiIndexReq{
		Id:      req.ID,
		AppCode: req.AppCode,
	})
	if err != nil {
		return nil, err
	}
	authInfo := make([]*types.AuthApiInfo, 0, len(info.List))
	for _, i := range info.List {
		authInfo = append(authInfo, &types.AuthApiInfo{
			Route:  i.Route,
			Method: i.Method,
		})
	}
	return &types.RoleApiIndexResp{List: authInfo, Total: info.Total}, nil
}
