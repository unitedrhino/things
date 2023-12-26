package tenantmanagelogic

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppApiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppApiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppApiCreateLogic {
	return &TenantAppApiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppApiCreateLogic) TenantAppApiCreate(in *sys.TenantApiInfo) (*sys.WithID, error) {
	// todo: add your logic here and delete this line

	return &sys.WithID{}, nil
}
