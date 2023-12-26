package tenantmanagelogic

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppMenuIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppMenuIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppMenuIndexLogic {
	return &TenantAppMenuIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppMenuIndexLogic) TenantAppMenuIndex(in *sys.MenuInfoIndexReq) (*sys.MenuInfoIndexResp, error) {
	// todo: add your logic here and delete this line

	return &sys.MenuInfoIndexResp{}, nil
}
