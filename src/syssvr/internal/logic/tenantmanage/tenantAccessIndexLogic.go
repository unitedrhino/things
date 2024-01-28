package tenantmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAccessIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAccessIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAccessIndexLogic {
	return &TenantAccessIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAccessIndexLogic) TenantAccessIndex(in *sys.TenantAccessIndexReq) (*sys.TenantAccessIndexResp, error) {
	tas, err := relationDB.NewTenantAccessRepo(l.ctx).FindByFilter(l.ctx, relationDB.TenantAccessFilter{TenantCode: in.Code}, nil)
	if err != nil {
		return nil, err
	}
	var accessCodes []string
	for _, v := range tas {
		accessCodes = append(accessCodes, v.AccessCode)
	}

	return &sys.TenantAccessIndexResp{AccessCodes: accessCodes}, nil
}
