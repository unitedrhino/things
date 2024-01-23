package info

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/logic/system"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/user"
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

func (l *CreateLogic) Create(req *types.TenantInfoCreateReq) (*types.WithID, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	if req.AdminUserInfo.UserName == "" {
		return nil, errors.Parameter.AddMsgf("需要填写管理员账号")
	}
	resp, err := l.svcCtx.TenantRpc.TenantInfoCreate(l.ctx, &sys.TenantInfoCreateReq{Info: system.ToTenantInfoRpc(req.Info), AdminUserInfo: user.UserInfoToRpc(req.AdminUserInfo)})
	return logic.SysToWithIDTypes(resp), err
}
