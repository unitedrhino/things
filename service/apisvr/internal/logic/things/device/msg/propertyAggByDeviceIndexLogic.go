package msg

import (
	"context"

	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyAggByDeviceIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 聚合属性历史记录,设备维度
func NewPropertyAggByDeviceIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyAggByDeviceIndexLogic {
	return &PropertyAggByDeviceIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyAggByDeviceIndexLogic) PropertyAggByDeviceIndex(req *types.DeviceMsgPropertyAggByDeviceIndexReq) (resp *types.DeviceMsgPropertyAggIndexResp, err error) {
	ret, err := l.svcCtx.DeviceMsg.PropertyAggByDeviceIndex(l.ctx, utils.Copy[dm.PropertyAggByDeviceIndexReq](req))
	return utils.Copy[types.DeviceMsgPropertyAggIndexResp](ret), err
}
