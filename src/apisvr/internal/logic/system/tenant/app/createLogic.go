package app

import (
	"context"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.TenantAppCreateReq) error {
	_, err := l.svcCtx.TenantRpc.TenantAppCreate(l.ctx, &sys.TenantAppCreateReq{
		Code:    req.Code,
		AppCode: req.AppCode,
		Modules: ToTenantAppModulesPb(req.Modules),
	})
	return err
}
