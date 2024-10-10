package info

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CanBindLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanBindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanBindLogic {
	return &CanBindLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanBindLogic) CanBind(req *types.DeviceInfoCanBindReq) error {
	_, err := l.svcCtx.DeviceM.DeviceInfoCanBind(l.ctx, utils.Copy[dm.DeviceInfoCanBindReq](req))
	return err
}
