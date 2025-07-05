package msg

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyAggIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 聚合属性历史记录
func NewPropertyAggIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyAggIndexLogic {
	return &PropertyAggIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyAggIndexLogic) PropertyAggIndex(req *types.DeviceMsgPropertyAggIndexReq) (resp *types.DeviceMsgPropertyAggIndexResp, err error) {
	l.svcCtx.DeviceMsg.PropertyAggIndex(l.ctx, utils.Copy[dm.PropertyAggIndexReq](req))
	return
}
