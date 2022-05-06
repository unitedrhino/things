package logic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/rootAuth"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type RootCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRootCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RootCheckLogic {
	return &RootCheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 鉴定是否是root账号
func (l *RootCheckLogic) RootCheck(in *dm.RootCheckReq) (*dm.Response, error) {
	l.Infof("RootCheck|req=%+v", in)
	if rootAuth.IsAdmin(l.svcCtx.Config.AuthWhite, in.Ip) {
		return &dm.Response{}, nil
	}
	return &dm.Response{}, errors.Permissions
}
