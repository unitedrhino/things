package module

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

func (l *CreateLogic) Create(req *types.TenantModuleCreateReq) error {
	_, err := l.svcCtx.TenantRpc.TenantAppModuleCreate(l.ctx, &sys.TenantModuleCreateReq{
		Code:       req.Code,
		AppCode:    req.AppCode,
		ModuleCode: req.ModuleCode,
		MenuIDs:    req.MenuIDs,
	})
	return err
}
