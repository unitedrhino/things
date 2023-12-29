package menumanagelogic

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuInfoDeleteLogic {
	return &MenuInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MenuInfoDeleteLogic) MenuInfoDelete(in *sys.WithID) (*sys.Response, error) {
	// todo: add your logic here and delete this line

	return &sys.Response{}, nil
}
