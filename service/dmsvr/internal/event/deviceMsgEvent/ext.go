package deviceMsgEvent

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/domain/deviceMsg/msgExt"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type ExtLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	dreq msgExt.Req
}

func NewExtLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExtLogic {
	return &ExtLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ExtLogic) initMsg(msg *deviceMsg.PublishMsg) error {
	var err error
	err = utils.Unmarshal(msg.Payload, &l.dreq)
	if err != nil {
		return errors.Parameter.AddDetailf("payload unmarshal payload:%v err:%v", string(msg.Payload), err)
	}

	return nil
}

func (l *ExtLogic) DeviceResp(msg *deviceMsg.PublishMsg, err error, data any) *deviceMsg.PublishMsg {
	resp := &deviceMsg.CommonMsg{
		Method:   deviceMsg.GetRespMethod(l.dreq.Method),
		MsgToken: l.dreq.MsgToken,
		//Timestamp: time.Now().UnixMilli(),
		Data: data,
	}
	return &deviceMsg.PublishMsg{
		Handle:       msg.Handle,
		Type:         msg.Type,
		Payload:      resp.AddStatus(err).Bytes(),
		Timestamp:    time.Now().UnixMilli(),
		ProductID:    msg.ProductID,
		DeviceName:   msg.DeviceName,
		ProtocolCode: msg.ProtocolCode,
	}
}

// ntp时间返回
func (l *ExtLogic) HandleGetNtpReply(msg *deviceMsg.PublishMsg, req msgExt.Req) (respMsg *deviceMsg.PublishMsg, err error) {
	/*
		校准时间公式推导
		deviceSendTime标记为t1，t2、t3、t4类似,如下:
		deviceSendTime(t1) ->  serverRecvTime(t2)
		deviceRecvTime(t4) <-  serverSendTime (t3)
		消息链路上的延迟 delay = [（t4-t1）- (t3-t2)] / 2 //注意，不能以 t2-t1 表示链路上的延迟，因为t2和t1是不同设备上的时间

		设备端使用如下公式校准时间：
		t4 + offset = t3+delay = (t4+t3+t2-t1)/2
	*/
	resp := &msgExt.Resp{
		//CommonMsg:      deviceMsg.NewRespCommonMsg(deviceMsg.GetNtpReply, "").AddStatus(errors.OK),
		DeviceSendTime: req.Timestamp,
		ServerRecvTime: msg.Timestamp, //这里是云端dd的mqtt接收到消息转给nats之前打的时间戳
	}

	respBytes, _ := json.Marshal(resp)
	_resp := deviceMsg.PublishMsg{
		Handle:     devices.Ext,
		Type:       msgExt.TypeNtp,
		Payload:    respBytes,
		Timestamp:  time.Now().UnixMilli(),
		ProductID:  msg.ProductID,
		DeviceName: msg.DeviceName,
	}

	return l.DeviceResp(&_resp, errors.OK, resp), nil
}

// ntp请求处理
func (l *ExtLogic) HandleNtp(msg *deviceMsg.PublishMsg) (respMsg *deviceMsg.PublishMsg, err error) {
	l.Debugf("%s req:%v", utils.FuncName(), msg)
	switch l.dreq.Method { //操作方法
	case deviceMsg.GetNtp:
		//if l.dreq.Code != errors.OK.Code { //如果不成功,则记录日志即可
		//	return nil, errors.DeviceResp.AddMsg(l.dreq.Msg).AddDetail(msg.Payload)
		//}
		respMsg, err = l.HandleGetNtpReply(msg, l.dreq)
		return respMsg, err
	default:
		return nil, errors.Method.AddMsg(l.dreq.Method)
	}
}

// Handle for topics.DeviceUpExtAll
func (l *ExtLogic) Handle(msg *deviceMsg.PublishMsg) (respMsg *deviceMsg.PublishMsg, err error) {
	l.Infof("%s req=%v", utils.FuncName(), msg)

	err = l.initMsg(msg)
	if err != nil {
		return nil, err
	}

	var action = devices.Ext
	respMsg, err = func() (respMsg *deviceMsg.PublishMsg, err error) {
		action = msg.Type
		switch msg.Type { //操作类型 从topic中提取 物模型下就是   property属性 event事件 action行为
		case msgExt.TypeNtp: //设备上报的 属性或信息
			return l.HandleNtp(msg)
		default:
			action = devices.Ext
			return nil, errors.Parameter.AddDetailf("ext types is err:%v", msg.Type)
		}
	}()
	if l.dreq.NoAsk() { //如果不需要回复
		respMsg = nil
	}
	hub := &deviceLog.Hub{
		ProductID:   msg.ProductID,
		Action:      action,
		Timestamp:   time.Now(), // 操作时间
		DeviceName:  msg.DeviceName,
		TraceID:     utils.TraceIdFromContext(l.ctx),
		RequestID:   l.dreq.MsgToken,
		Content:     string(msg.Payload),
		Topic:       msg.Topic,
		ResultCode:  errors.Fmt(err).GetCode(),
		RespPayload: respMsg.GetPayload(),
	}
	_ = l.svcCtx.HubLogRepo.Insert(l.ctx, hub)
	l.svcCtx.UserSubscribe.Publish(l.ctx, def.UserSubscribeDevicePublish, hub.ToApp(), map[string]any{
		"productID":  msg.ProductID,
		"deviceName": msg.DeviceName,
	})
	return
}
