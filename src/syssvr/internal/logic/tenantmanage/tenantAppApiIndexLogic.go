package tenantmanagelogic

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppApiIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppApiIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppApiIndexLogic {
	return &TenantAppApiIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppApiIndexLogic) TenantAppApiIndex(in *sys.TenantAppApiIndexReq) (*sys.ApiInfoIndexResp, error) {
	// todo: add your logic here and delete this line

	return &sys.ApiInfoIndexResp{}, nil
}
