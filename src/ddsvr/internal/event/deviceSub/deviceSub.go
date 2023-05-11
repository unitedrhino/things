package deviceSub

import (
	"context"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/domain/deviceAuth"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/ddsvr/internal/domain/custom"
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

	return s.svcCtx.PubInner.DevPubMsg(s.ctx, pub)
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
	finalPayload := payload
	handle := strings.TrimPrefix(topicInfo.TopicHead, "$")
	if len(topicInfo.Types) > 1 && topicInfo.Types[1] == custom.CustomType {
		//自定义协议
		f, err := s.svcCtx.Script.GetProtoFunc(s.ctx, topicInfo.ProductID, custom.ConvertTypeUp, handle, topicInfo.Types[0])
		if err != nil {
			s.Errorf("%s.GetProtoFunc topicInfo:%#v err:%v", utils.FuncName(), topicInfo, err)
			return nil, err
		}
		if f == nil {
			s.Errorf("%s.GetProtoFunc topicInfo:%#v transform func not find", utils.FuncName(), topicInfo)
			return nil, errors.Parameter.AddMsg("转换函数未找到")
		}
		finalPayload, err = f(payload)
		if err != nil {
			s.Errorf("%s.Transform topicInfo:%#v err:%v", utils.FuncName(), topicInfo, err)
			return nil, err
		}
		s.Infof("%s.transform success before:%#v after:%#v", utils.FuncName(), payload, finalPayload)
	}
	return &devices.DevPublish{
		Topic:      topic,
		Timestamp:  time.Now().UnixMilli(),
		Payload:    finalPayload,
		Handle:     handle,
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
