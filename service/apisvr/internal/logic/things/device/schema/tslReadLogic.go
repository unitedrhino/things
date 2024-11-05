package schema

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TslReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取产品物模型tsl
func NewTslReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TslReadLogic {
	return &TslReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TslReadLogic) TslRead(req *types.DeviceSchemaTslReadReq) (resp *types.DeviceSchemaTslReadResp, err error) {
	ret, err := l.svcCtx.DeviceM.DeviceSchemaTslRead(l.ctx, utils.Copy[dm.DeviceSchemaTslReadReq](req))

	return utils.Copy[types.DeviceSchemaTslReadResp](ret), err
}
