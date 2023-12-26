package tenantmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantInfoReadLogic {
	return &TenantInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取区域信息详情
func (l *TenantInfoReadLogic) TenantInfoRead(in *sys.WithIDCode) (*sys.TenantInfo, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	ti, err := relationDB.NewTenantInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantInfoFilter{Codes: []string{in.Code}, ID: in.Id})

	return ToTenantInfoRpc(ti), err
}
