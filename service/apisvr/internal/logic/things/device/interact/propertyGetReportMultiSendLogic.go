package interact

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyGetReportMultiSendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量请求设备获取设备最新属性
func NewPropertyGetReportMultiSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyGetReportMultiSendLogic {
	return &PropertyGetReportMultiSendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyGetReportMultiSendLogic) PropertyGetReportMultiSend(req *types.DeviceInteractPropertyGetReportMultiSendReq) (resp *types.DeviceInteractPropertyGetReportMultiSendResp, err error) {
	dmResp, err := l.svcCtx.DeviceInteract.PropertyGetReportMultiSend(l.ctx, utils.Copy[dm.PropertyGetReportMultiSendReq](req))
	return utils.Copy[types.DeviceInteractPropertyGetReportMultiSendResp](dmResp), err
}
