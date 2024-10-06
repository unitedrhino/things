package info

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"gitee.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindLogic {
	return &BindLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindLogic) Bind(req *types.DeviceInfoBindReq) error {
	_, err := l.svcCtx.DeviceM.DeviceInfoBind(l.ctx, utils.Copy[dm.DeviceInfoBindReq](req))
	return err
}
