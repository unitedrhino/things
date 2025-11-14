package msg

import (
	"context"

	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyLatestAggIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 聚合属性最新值
func NewPropertyLatestAggIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyLatestAggIndexLogic {
	return &PropertyLatestAggIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyLatestAggIndexLogic) PropertyLatestAggIndex(req *types.DeviceMsgPropertyLatestAggIndexReq) (resp *types.DeviceMsgPropertyLatestAggIndexResp, err error) {
	ret, err := l.svcCtx.DeviceMsg.PropertyLatestAggIndex(l.ctx, utils.Copy[dm.PropertyLatestAggIndexReq](req))
	if err != nil {
		return nil, err
	}
	return utils.Copy[types.DeviceMsgPropertyLatestAggIndexResp](ret), nil
}
