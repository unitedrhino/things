package msg

import (
	"context"

	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyLogAggIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 聚合属性历史记录
func NewPropertyLogAggIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyLogAggIndexLogic {
	return &PropertyLogAggIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyLogAggIndexLogic) PropertyLogAggIndex(req *types.DeviceMsgPropertyLogAggIndexReq) (resp *types.DeviceMsgPropertyLogAggIndexResp, err error) {
	ret, err := l.svcCtx.DeviceMsg.PropertyLogAggIndex(l.ctx, utils.Copy[dm.PropertyAggIndexReq](req))
	return utils.Copy[types.DeviceMsgPropertyLogAggIndexResp](ret), err
}
