package module

import (
	"context"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.TenantModuleWithIDOrCode) error {
	_, err := l.svcCtx.TenantRpc.TenantAppModuleDelete(l.ctx, &sys.TenantModuleWithIDOrCode{
		Code:       req.Code,
		AppCode:    req.AppCode,
		ModuleCode: req.ModuleCode,
		Id:         req.ID,
	})
	return err
}
