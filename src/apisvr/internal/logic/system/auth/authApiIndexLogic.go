package auth

import (
	"context"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthApiIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthApiIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthApiIndexLogic {
	return &AuthApiIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthApiIndexLogic) AuthApiIndex(req *types.AuthApiIndexReq) (resp *types.AuthApiIndexResp, err error) {
	info, err := l.svcCtx.RoleRpc.RoleApiIndex(l.ctx, &sys.RoleApiIndexReq{
		RoleID: req.RoleID,
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

	return &types.AuthApiIndexResp{List: authInfo, Total: info.Total}, nil
}
