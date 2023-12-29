package module

import (
	"context"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MulitUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMulitUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MulitUpdateLogic {
	return &MulitUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MulitUpdateLogic) MulitUpdate(req *types.RoleModuleMultiUpdateReq) error {
	_, err := l.svcCtx.RoleRpc.RoleModuleMultiUpdate(l.ctx, &sys.RoleModuleMultiUpdateReq{
		Id: req.ID, AppCode: req.AppCode, ModuleCodes: req.ModuleCodes})
	return err
}
