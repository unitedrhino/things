package menumanagelogic

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuInfoCreateLogic {
	return &MenuInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MenuInfoCreateLogic) MenuInfoCreate(in *sys.MenuInfo) (*sys.WithID, error) {
	// todo: add your logic here and delete this line

	return &sys.WithID{}, nil
}
