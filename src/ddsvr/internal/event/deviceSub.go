package event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/src/ddsvr/ddExport"
	"github.com/i-Things/things/src/ddsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type DeviceSubServer struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewDeviceSubServer(svcCtx *svc.ServiceContext, ctx context.Context) *DeviceSubServer {
	return &DeviceSubServer{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

// Publish 设备发布的信息通过nats转发给内部服务
func (s *DeviceSubServer) Publish(topic string, payload []byte) error {
	s.Info("DeviceSubServer", "Publish", topic, string(payload))
	productId, deviceName, err := devices.GetDeviceInfo(topic)
	if err != nil {
		return err
	}
	pub := ddExport.DevPublish{
		Timestamp:  time.Now().UnixMilli(),
		Topic:      topic,
		Payload:    payload,
		ProductID:  productId,
		DeviceName: deviceName,
	}
	pubStr, _ := json.Marshal(pub)
	return s.svcCtx.InnerLink.Publish(s.ctx, fmt.Sprintf(ddExport.TopicDevPublish, productId, deviceName), pubStr)
}

func (s *DeviceSubServer) Connected(info *ddExport.DevConn) error {
	s.Info("Connected", info)
	str, _ := json.Marshal(info)
	return s.svcCtx.InnerLink.Publish(s.ctx, ddExport.TopicDevConnected, str)
}
func (s *DeviceSubServer) Disconnected(info *ddExport.DevConn) error {
	s.Info("Disconnected", info)
	str, _ := json.Marshal(info)
	return s.svcCtx.InnerLink.Publish(s.ctx, ddExport.TopicDevDisconnected, str)
}
