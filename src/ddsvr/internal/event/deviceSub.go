package event

import (
	"context"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/traces"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/ddsvr/internal/repo/event/publish/pubInner"
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

// Thing 设备发布物模型消息的信息通过nats转发给内部服务
func (s *DeviceSubServer) Thing(topic string, payload []byte) error {
	pub, err := s.getDevPublish(topic, payload)
	if pub == nil {
		return err
	}
	ctx1, span := traces.StartSpan(s.ctx, topic, "")

	logx.Infof("mqtt.Thing trace:%s, spanID:%s topic:%s payload:%v",
		span.SpanContext().TraceID(), span.SpanContext().SpanID(), topic, string(payload))
	defer span.End()
	return s.svcCtx.PubInner.DevPubThing(ctx1, pub)
}

// Ota ota远程升级
func (s *DeviceSubServer) Ota(topic string, payload []byte) error {
	s.Infof("%s topic:%v payload:%v", utils.FuncName(), topic, string(payload))
	pub, err := s.getDevPublish(topic, payload)
	if pub == nil {
		return err
	}
	return s.svcCtx.PubInner.DevPubOta(s.ctx, pub)
}

// Config 设备远程配置
func (s *DeviceSubServer) Config(topic string, payload []byte) error {
	s.Infof("%s topic:%v payload:%v", utils.FuncName(), topic, string(payload))
	pub, err := s.getDevPublish(topic, payload)
	if pub == nil {
		return err
	}
	return s.svcCtx.PubInner.DevPubConfig(s.ctx, pub)
}

// Shadow 设备影子
func (s *DeviceSubServer) Shadow(topic string, payload []byte) error {
	s.Infof("%s topic:%v payload:%v", utils.FuncName(), topic, string(payload))
	pub, err := s.getDevPublish(topic, payload)
	if pub == nil {
		return err
	}
	return s.svcCtx.PubInner.DevPubShadow(s.ctx, pub)
}

// Log 设备调试日志
func (s *DeviceSubServer) SDKLog(topic string, payload []byte) error {
	s.Infof("%s topic:%v payload:%v", utils.FuncName(), topic, string(payload))
	pub, err := s.getDevPublish(topic, payload)
	if pub == nil {
		return err
	}
	return s.svcCtx.PubInner.DevPubSDKLog(s.ctx, pub)
}

// Log 设备调试日志
func (s *DeviceSubServer) Gateway(topic string, payload []byte) error {
	s.Infof("%s topic:%v payload:%v", utils.FuncName(), topic, string(payload))
	pub, err := s.getDevPublish(topic, payload)
	if pub == nil {
		return err
	}
	return s.svcCtx.PubInner.DevPubSDKLog(s.ctx, pub)
}

func (s *DeviceSubServer) getDevPublish(topic string, payload []byte) (*devices.DevPublish, error) {
	topicInfo, err := devices.GetTopicInfo(topic)
	if err != nil {
		return nil, err
	}
	if topicInfo.Direction == devices.Down {
		//服务器端下发的消息直接忽略
		return nil, nil
	}
	return &devices.DevPublish{
		Timestamp:  time.Now().UnixMilli(),
		Topic:      topic,
		Payload:    payload,
		ProductID:  topicInfo.ProductID,
		DeviceName: topicInfo.DeviceName,
	}, nil
}

func (s *DeviceSubServer) Connected(info *devices.DevConn) error {
	s.Infof("%s info:%v", utils.FuncName(), utils.Fmt(info))
	return s.svcCtx.PubInner.PubConn(s.ctx, pubInner.Connect, info)
}
func (s *DeviceSubServer) Disconnected(info *devices.DevConn) error {
	s.Infof("%s info:%v", utils.FuncName(), utils.Fmt(info))
	return s.svcCtx.PubInner.PubConn(s.ctx, pubInner.DisConnect, info)
}
