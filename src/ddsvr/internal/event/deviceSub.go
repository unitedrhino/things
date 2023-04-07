package event

import (
	"context"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/domain/deviceAuth"
	"github.com/i-Things/things/shared/traces"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/ddsvr/internal/repo/event/publish/pubInner"
	"github.com/i-Things/things/src/ddsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
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

// Msg 设备发布物模型消息的信息通过nats转发给内部服务
func (s *DeviceSubServer) Msg(topic string, payload []byte) error {
	pub, err := s.getDevPublish(topic, payload)
	if pub == nil {
		return err
	}
	ctx1, span := traces.StartSpan(s.ctx, topic, "")

	logx.Infof("mqtt.Msg trace:%s, spanID:%s topic:%s payload:%v",
		span.SpanContext().TraceID(), span.SpanContext().SpanID(), topic, string(payload))
	defer span.End()
	return s.svcCtx.PubInner.DevPubMsg(ctx1, pub)
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
		Topic:      topic,
		Timestamp:  time.Now().UnixMilli(),
		Payload:    payload,
		Handle:     strings.TrimPrefix(topicInfo.TopicHead, "$"),
		Type:       topicInfo.Types[0],
		ProductID:  topicInfo.ProductID,
		DeviceName: topicInfo.DeviceName,
	}, nil
}

func (s *DeviceSubServer) Connected(info *devices.DevConn) error {
	s.Infof("%s info:%v", utils.FuncName(), utils.Fmt(info))
	_, err := deviceAuth.GetLoginDevice(info.UserName)
	if err != nil { //只传送设备的消息
		return nil
	}
	return s.svcCtx.PubInner.PubConn(s.ctx, pubInner.Connect, info)
}
func (s *DeviceSubServer) Disconnected(info *devices.DevConn) error {
	s.Infof("%s info:%v", utils.FuncName(), utils.Fmt(info))
	_, err := deviceAuth.GetLoginDevice(info.UserName)
	if err != nil { //只传送设备的消息
		return nil
	}
	return s.svcCtx.PubInner.PubConn(s.ctx, pubInner.DisConnect, info)
}
