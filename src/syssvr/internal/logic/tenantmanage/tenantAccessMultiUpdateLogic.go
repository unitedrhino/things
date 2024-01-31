package tenantmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAccessMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAccessMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAccessMultiUpdateLogic {
	return &TenantAccessMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAccessMultiUpdateLogic) TenantAccessMultiUpdate(in *sys.TenantAccessMultiUpdateReq) (*sys.Response, error) {
	err := relationDB.NewTenantAccessRepo(l.ctx).MultiUpdate(l.ctx, in.Code, in.AccessCodes)
	return &sys.Response{}, err
}
