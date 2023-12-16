package tenantmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppMultiUpdateLogic {
	return &TenantAppMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppMultiUpdateLogic) TenantAppMultiUpdate(in *sys.TenantAppMultiUpdateReq) (*sys.Response, error) {
	err := logic.IsSupperAdmin(l.ctx, def.TenantCodeDefault)
	if err != nil {
		return nil, err
	}
	err = relationDB.NewTenantAppRepo(l.ctx).MultiUpdate(l.ctx, in.Code, in.AppCodes)

	return &sys.Response{}, err
}
