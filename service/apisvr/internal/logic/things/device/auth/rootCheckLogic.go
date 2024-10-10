package auth

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RootCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRootCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RootCheckLogic {
	return &RootCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RootCheckLogic) RootCheck(req *types.DeviceAuthRootCheckReq) error {
	l.Infof("%s req=%v", utils.FuncName(), req)
	_, err := l.svcCtx.DeviceM.RootCheck(l.ctx, &dm.RootCheckReq{
		Username:    req.Username,
		Password:    req.Password,
		ClientID:    req.ClientID,
		Ip:          req.Ip,
		Certificate: req.Certificate,
	})
	return err
}
