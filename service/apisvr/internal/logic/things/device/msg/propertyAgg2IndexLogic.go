package msg

import (
	"context"

	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyAgg2IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 聚合属性历史记录,设备维度
func NewPropertyAgg2IndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyAgg2IndexLogic {
	return &PropertyAgg2IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyAgg2IndexLogic) PropertyAgg2Index(req *types.DeviceMsgPropertyAgg2IndexReq) (resp *types.DeviceMsgPropertyAggIndexResp, err error) {
	ret, err := l.svcCtx.DeviceMsg.PropertyAgg2Index(l.ctx, utils.Copy[dm.PropertyAgg2IndexReq](req))
	return utils.Copy[types.DeviceMsgPropertyAggIndexResp](ret), err
}
