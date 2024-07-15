package deviceSub

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/deviceAuth"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/domain/deviceMsg/msgExt"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/sdk/service/protocol"
	"github.com/i-Things/things/service/dgsvr/internal/domain/custom"
	"github.com/i-Things/things/service/dgsvr/internal/event/innerSub"
	deviceauthlogic "github.com/i-Things/things/service/dgsvr/internal/logic/deviceauth"
	"github.com/i-Things/things/service/dgsvr/internal/repo/event/publish/pubInner"
	"github.com/i-Things/things/service/dgsvr/internal/svc"
	"github.com/i-Things/things/service/dgsvr/pb/dg"
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
func (s *DeviceSubServer) ExtMsg(topic string, payload []byte) error {
	pub, err := s.getDevPublish(topic, payload)
	if pub == nil || err != nil {
		return err
	}
	if pub.Type != msgExt.TypeRegister {
		return s.svcCtx.PubInner.DevPubMsg(s.ctx, pub)
	}
	var (
		resp = devices.InnerPublish{
			Handle:    pub.Handle,
			ProductID: pub.ProductID,
			Type:      pub.Type,
			//Payload:pub.Payload,
			DeviceName:   pub.DeviceName,
			ProtocolCode: pub.ProtocolCode,
		}
		respPayload = msgExt.RespRegister{}
	)
	var req msgExt.RegisterReq
	err = json.Unmarshal(payload, &req)
	if err != nil {
		return err
	}
	respPayload.CommonMsg = &deviceMsg.CommonMsg{
		MsgToken:  req.MsgToken,
		Timestamp: time.Now().UnixMilli(),
		Code:      errors.OK.GetCode(),
		Msg:       errors.OK.GetMsg(),
	}
	ret, err := deviceauthlogic.NewDeviceRegisterLogic(s.ctx, s.svcCtx).DeviceRegister(&dg.DeviceRegisterReq{
		ProductID:  pub.ProductID,
		DeviceName: pub.DeviceName,
		Nonce:      req.Payload.Nonce,
		Timestamp:  req.Payload.Timestamp,
		Signature:  req.Payload.Signature,
	})
	if err != nil {
		s.Error(err)
		e := errors.Fmt(err)
		respPayload.Code = e.Code
		respPayload.Msg = e.GetMsg()
	} else {
		respPayload.Data = msgExt.RespData{
			Len:     ret.Len,
			Payload: ret.Payload,
		}
	}
	resp.Payload, err = json.Marshal(respPayload)
	if err != nil {
		return err
	}
	err = innerSub.NewInnerSubServer(s.svcCtx, s.ctx).PublishToDev(&resp)
	return err
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
	if len(topicInfo.Types) > 1 && topicInfo.Types[1] == custom.CustomType { //协议转换脚本
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
	} else if handle == "custom" { //自定义topic
		f, err := s.svcCtx.Script.GetTransFormFunc(s.ctx, topicInfo.ProductID, topicInfo.Direction)
		if err != nil {
			s.Errorf("%s.GetTransFormFunc topicInfo:%#v err:%v", utils.FuncName(), topicInfo, err)
			return nil, err
		}
		if f == nil {
			s.Errorf("%s.GetTransFormFunc topicInfo:%#v transform func not find", utils.FuncName(), topicInfo)
			return nil, errors.Parameter.AddMsg("自定义topic转换函数未找到")
		}
		ret, err := f(topic, payload)
		if err != nil {
			s.Errorf("%s.TransFormFunc topic:%v payload:%v ret err:%v", utils.FuncName(), topic, string(payload), err)
			return nil, err
		}
		if ret.Topic == "" || len(ret.PayLoad) == 0 {
			s.Errorf("%s.TransFormFunc getEmpty topic:%v payload:%v  ret.Topic:%v,ret.PayLoad:%v err:%v",
				utils.FuncName(), topic, string(payload), ret.Topic, string(ret.PayLoad), err)
			return nil, err
		}
		ti, err := devices.GetTopicInfo(ret.Topic)
		if err != nil {
			return nil, err
		}
		if ti.TopicHead == "$custom" { //禁止循环嵌套
			return nil, errors.Parameter.AddMsg("禁止topic循环嵌套")
		}
		return s.getDevPublish(ret.Topic, ret.PayLoad)
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
	newDo, err := deviceAuth.GetLoginDevice(info.UserName)
	if err != nil { //只传送设备的消息
		return nil
	}
	dev := devices.Core{
		ProductID:  newDo.ProductID,
		DeviceName: newDo.DeviceName,
	}
	err = caches.GetStore().Hset(protocol.DeviceMqttDevice, protocol.GenDeviceTopicKey(dev), utils.MarshalNoErr(info))
	if err != nil {
		logx.Error(err)
	}
	err = caches.GetStore().Hset(protocol.DeviceMqttClientID, info.ClientID, utils.MarshalNoErr(dev))
	if err != nil {
		logx.Error(err)
	}
	return s.svcCtx.PubInner.PubConn(s.ctx, pubInner.Connect, info)
}

func (s *DeviceSubServer) Disconnected(info *devices.DevConn) error {
	s.Infof("%s info:%v", utils.FuncName(), utils.Fmt(info))
	_, err := deviceAuth.GetLoginDevice(info.UserName)
	if err != nil { //只传送设备的消息
		return nil
	}
	if info.Reason == "takeovered" || info.Reason == "takenover" || info.Reason == "discard" || info.Reason == "discarded" { //连接还在的时候被别人顶了,忽略这种下线
		logx.Errorf("忽略的下线状态:%v", utils.Fmt(info))
		return nil
	}
	dev := devices.Core{
		ProductID:  info.ProductID,
		DeviceName: info.DeviceName,
	}
	protocol.DeleteDeviceActivity(s.ctx, dev)
	return s.svcCtx.PubInner.PubConn(s.ctx, pubInner.DisConnect, info)
}
