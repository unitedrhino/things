package menulogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuDeleteLogic {
	return &MenuDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MenuDeleteLogic) MenuDelete(in *sys.MenuDeleteReq) (*sys.Response, error) {
	err := l.svcCtx.MenuModel.DeleteMenu(in.Id)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &sys.Response{}, nil
}
