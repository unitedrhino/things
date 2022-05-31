package event

import (
	"context"
	"github.com/i-Things/things/shared/devices"
	tr "github.com/i-Things/things/shared/trace"
	"github.com/i-Things/things/src/ddsvr/internal/repo/event/innerLink"
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
	s.Infof("DeviceSubServer|Publish|topic:%v payload:%v", topic, string(payload))
	topicInfo, err := devices.GetTopicInfo(topic)
	if err != nil {
		return err
	}
	if topicInfo.Direction == devices.DOWN {
		//服务器端下发的消息直接忽略
		return nil
	}
	pub := devices.DevPublish{
		Timestamp:  time.Now().UnixMilli(),
		Topic:      topic,
		Payload:    payload,
		ProductID:  topicInfo.ProductID,
		DeviceName: topicInfo.DeviceName,
	}
	ctx1, span := tr.StartSpan(s.ctx, topic, "")

	logx.Infof("[mqtt.SubScribe]|-------------------trace:%s, spanid:%s|topic:%s",
		span.SpanContext().TraceID(), span.SpanContext().SpanID(), topic)
	defer span.End()

	return s.svcCtx.InnerLink.PubDevPublish(ctx1, pub)
}

func (s *DeviceSubServer) Connected(info *devices.DevConn) error {
	s.Info("Connected", info)
	return s.svcCtx.InnerLink.PubConn(s.ctx, innerLink.Connect, info)
}
func (s *DeviceSubServer) Disconnected(info *devices.DevConn) error {
	s.Info("Disconnected", info)
	return s.svcCtx.InnerLink.PubConn(s.ctx, innerLink.DisConnect, info)
}
