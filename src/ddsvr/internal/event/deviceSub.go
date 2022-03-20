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
	logx.Logger
}

func NewDeviceSubServer(svcCtx *svc.ServiceContext, ctx context.Context) *DeviceSubServer {
	return &DeviceSubServer{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (s *DeviceSubServer) Publish(topic string, payload []byte) error {
	payloadStr := string(payload)
	s.Info("Publish", topic, payloadStr)
	pub := ddDef.DevPublish{
		Ts:      time.Now().UnixMilli(),
		Topic:   topic,
		Payload: payloadStr,
	}
	pubStr, _ := json.Marshal(pub)
	return s.svcCtx.InnerLink.Publish(ddDef.TopicDevPublish, pubStr)
}

func (s *DeviceSubServer) Login(info *ddDef.DevLogInOut) error {
	s.Info("Login", info)
	str, _ := json.Marshal(info)
	return s.svcCtx.InnerLink.Publish(ddDef.TopicDevLogin, str)
}
func (s *DeviceSubServer) Logout(info *ddDef.DevLogInOut) error {
	s.Info("Logout", info)
	str, _ := json.Marshal(info)
	return s.svcCtx.InnerLink.Publish(ddDef.TopicDevLogout, str)
}
