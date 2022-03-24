package event

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/ddsvr/ddDef"
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

func (s *DeviceSubServer) Publish(topic string, payload []byte) error {
	s.Info("Publish", topic, string(payload))
	pub := ddDef.DevPublish{
		Timestamp: time.Now().UnixMilli(),
		Topic:     topic,
		Payload:   payload,
	}
	pubStr, _ := json.Marshal(pub)
	return s.svcCtx.InnerLink.Publish(s.ctx, ddDef.TopicDevPublish, pubStr)
}

func (s *DeviceSubServer) Connected(info *ddDef.DevConn) error {
	s.Info("Connected", info)
	str, _ := json.Marshal(info)
	return s.svcCtx.InnerLink.Publish(s.ctx, ddDef.TopicDevConnected, str)
}
func (s *DeviceSubServer) Disconnected(info *ddDef.DevConn) error {
	s.Info("Disconnected", info)
	str, _ := json.Marshal(info)
	return s.svcCtx.InnerLink.Publish(s.ctx, ddDef.TopicDevDisconnected, str)
}
