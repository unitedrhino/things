package tenantmanagelogic

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppMenuCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppMenuCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppMenuCreateLogic {
	return &TenantAppMenuCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppMenuCreateLogic) TenantAppMenuCreate(in *sys.TenantMenuInfo) (*sys.WithID, error) {
	// todo: add your logic here and delete this line

	return &sys.WithID{}, nil
}
