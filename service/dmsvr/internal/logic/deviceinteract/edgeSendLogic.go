package deviceinteractlogic

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/domain/deviceMsg"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/event/deviceMsgEvent"
	"time"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type EdgeSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEdgeSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EdgeSendLogic {
	return &EdgeSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 提供给边缘端进行http访问
func (l *EdgeSendLogic) EdgeSend(in *dm.EdgeSendReq) (*dm.EdgeSendResp, error) {
	var msg = &deviceMsg.PublishMsg{
		Handle:       in.Handle,
		Type:         in.Type,
		Payload:      in.Payload,
		Timestamp:    time.Now().UnixMilli(),
		ProductID:    in.ProductID,
		DeviceName:   in.DeviceName,
		ProtocolCode: def.ProtocolCodeUnitedRhino,
	}
	var resp *deviceMsg.PublishMsg
	var err error
	switch in.Handle {
	case devices.Thing:
		resp, err = deviceMsgEvent.NewThingLogic(l.ctx, l.svcCtx).Handle(msg)
	case devices.Ota:
		resp, err = deviceMsgEvent.NewOtaLogic(l.ctx, l.svcCtx).Handle(msg)
	case devices.Config:
		resp, err = deviceMsgEvent.NewConfigLogic(l.ctx, l.svcCtx).Handle(msg)
	case devices.Log:
		resp, err = deviceMsgEvent.NewSDKLogLogic(l.ctx, l.svcCtx).Handle(msg)
	case devices.Shadow:
		resp, err = deviceMsgEvent.NewShadowLogic(l.ctx, l.svcCtx).Handle(msg)
	case devices.Gateway:
		resp, err = deviceMsgEvent.NewGatewayLogic(l.ctx, l.svcCtx).Handle(msg)
	case devices.Ext:
		resp, err = deviceMsgEvent.NewExtLogic(l.ctx, l.svcCtx).Handle(msg)
	}
	if err != nil {
		return nil, err
	}
	return &dm.EdgeSendResp{Payload: resp.Payload}, nil
}
