package info

import (
	"context"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"gitee.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbindLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnbindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindLogic {
	return &UnbindLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnbindLogic) Unbind(req *types.DeviceCore) error {
	_, err := l.svcCtx.DeviceM.DeviceInfoUnbind(l.ctx, &dm.DeviceCore{DeviceName: req.DeviceName, ProductID: req.ProductID})

	return err
}
