package menumanagelogic

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuInfoIndexLogic {
	return &MenuInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MenuInfoIndexLogic) MenuInfoIndex(in *sys.MenuInfoIndexReq) (*sys.MenuInfoIndexResp, error) {
	// todo: add your logic here and delete this line

	return &sys.MenuInfoIndexResp{}, nil
}
