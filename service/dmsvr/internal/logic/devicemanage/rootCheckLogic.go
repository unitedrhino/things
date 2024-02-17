package devicemanagelogic

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

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

// 鉴定是否是root账号(提供给mqtt broker)
func (l *RootCheckLogic) RootCheck(in *dm.RootCheckReq) (*dm.Empty, error) {
	l.Infof("%s req=%+v", utils.FuncName(), in)
	if utils.Auth(l.svcCtx.Config.AuthWhite, in.Username, in.Password, in.Ip) {
		return &dm.Empty{}, nil
	}
	return &dm.Empty{}, errors.Permissions
}
