package menu

import (
	"context"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/module/menu"
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

func (l *CreateLogic) Create(req *types.TenantAppMenu) (resp *types.WithID, err error) {
	ret, err := l.svcCtx.TenantRpc.TenantAppMenuCreate(l.ctx, &sys.TenantAppMenu{
		TemplateID: req.TemplateID,
		Code:       req.Code,
		AppCode:    req.AppCode,
		Info:       menu.ToMenuInfoRpc(&req.MenuInfo),
	})
	if err != nil {
		return nil, err
	}
	return &types.WithID{
		ID: ret.Id,
	}, nil
}
