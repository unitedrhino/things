package msg

import (
	"context"

	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyLogAggByDeviceIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 聚合属性历史记录,设备维度
func NewPropertyLogAggByDeviceIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyLogAggByDeviceIndexLogic {
	return &PropertyLogAggByDeviceIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyLogAggByDeviceIndexLogic) PropertyLogAggByDeviceIndex(req *types.DeviceMsgPropertyAggByDeviceIndexReq) (resp *types.DeviceMsgPropertyLogAggIndexResp, err error) {
	ret, err := l.svcCtx.DeviceMsg.PropertyLogAggByDeviceIndex(l.ctx, utils.Copy[dm.PropertyLogAggByDeviceIndexReq](req))
	return utils.Copy[types.DeviceMsgPropertyLogAggIndexResp](ret), err
}
