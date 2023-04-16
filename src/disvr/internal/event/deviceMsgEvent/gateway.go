package deviceMsgEvent

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/domain/application"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgGateway"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgHubLog"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"
)

type GatewayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	dreq   msgGateway.Msg
	topics []string
}

func NewGatewayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayLogic {
	return &GatewayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *GatewayLogic) initMsg(msg *deviceMsg.PublishMsg) (err error) {
	err = utils.Unmarshal(msg.Payload, &l.dreq)
	if err != nil {
		return errors.Parameter.AddDetailf("payload unmarshal payload:%v err:%v", string(msg.Payload), err)
	}
	l.topics = strings.Split(msg.Topic, "/")
	if len(l.topics) < 5 || l.topics[1] != "up" {
		return errors.Parameter.AddDetail("initMsg topic is err:" + msg.Topic)
	}
	return nil
}

func (l *GatewayLogic) Handle(msg *deviceMsg.PublishMsg) (respMsg *deviceMsg.PublishMsg, err error) {
	l.Infof("%s req=%+v", utils.FuncName(), msg)
	err = l.initMsg(msg)
	if err != nil {
		return nil, err
	}
	var (
		resp *msgGateway.Msg
	)

	switch l.topics[2] {
	case msgGateway.TypeOperation:
		resp, err = l.HandleOperation(msg)
	case msgGateway.TypeStatus:
		resp, err = l.HandleStatus(msg)
	}
	respStr, _ := json.Marshal(resp)
	l.svcCtx.HubLogRepo.Insert(l.ctx, &msgHubLog.HubLog{
		ProductID:  msg.ProductID,
		Action:     "gateway",
		Timestamp:  time.Now(), // 记录当前时间
		DeviceName: msg.DeviceName,
		TranceID:   utils.TraceIdFromContext(l.ctx),
		RequestID:  l.dreq.ClientToken,
		Content:    string(msg.Payload),
		Topic:      msg.Topic,
		ResultType: errors.Fmt(err).GetCode(),
	})
	return &deviceMsg.PublishMsg{
		Topic:      deviceMsg.GenRespTopic(msg.Topic),
		Payload:    respStr,
		Timestamp:  time.Now(),
		ProductID:  msg.ProductID,
		DeviceName: msg.DeviceName,
	}, nil
}

func (l *GatewayLogic) HandleOperation(msg *deviceMsg.PublishMsg) (respMsg *msgGateway.Msg, err error) {
	l.Debugf("%s", utils.FuncName())
	var resp = msgGateway.Msg{
		CommonMsg: deviceMsg.NewRespCommonMsg(l.dreq.Method, l.dreq.ClientToken),
	}
	resp.AddStatus(errors.OK)
	switch l.dreq.Method {
	case deviceMsg.Bind:
		_, err := l.svcCtx.DeviceM.DeviceGatewayMultiCreate(l.ctx, &dm.DeviceGatewayMultiCreateReq{
			GatewayProductID:  msg.ProductID,
			GatewayDeviceName: msg.DeviceName,
			List:              ToDmDevicesCore(l.dreq.Payload.Devices),
		})
		if err != nil {
			resp.AddStatus(err)
			return &resp, err
		}
		resp.Payload = &msgGateway.GatewayPayload{Devices: l.dreq.Payload.Devices}
	case deviceMsg.Unbind:
		_, err := l.svcCtx.DeviceM.DeviceGatewayMultiDelete(l.ctx, &dm.DeviceGatewayMultiDeleteReq{
			GatewayProductID:  msg.ProductID,
			GatewayDeviceName: msg.DeviceName,
			List:              ToDmDevicesCore(l.dreq.Payload.Devices),
		})
		if err != nil {
			resp.AddStatus(err)
			return &resp, err
		}
		resp.Payload = &msgGateway.GatewayPayload{Devices: l.dreq.Payload.Devices}
	case deviceMsg.DescribeSubDevices:
		deviceList, err := l.svcCtx.DeviceM.DeviceGatewayIndex(l.ctx, &dm.DeviceGatewayIndexReq{
			GatewayProductID:  msg.ProductID,
			GatewayDeviceName: msg.DeviceName,
		})
		if err != nil {
			resp.AddStatus(err)
			return &resp, err
		}
		var payload msgGateway.GatewayPayload
		for _, device := range deviceList.List {
			payload.Devices = append(payload.Devices, &msgGateway.Device{
				ProductID:  device.ProductID,
				DeviceName: device.DeviceName,
				Result:     errors.OK.Code,
			})
		}
		resp.Payload = &payload
	default:
		return nil, errors.Method.AddMsg(l.dreq.Method)

	}
	return &resp, err
}

func (l *GatewayLogic) HandleStatus(msg *deviceMsg.PublishMsg) (respMsg *msgGateway.Msg, err error) {
	l.Debugf("%s", utils.FuncName())
	var resp = msgGateway.Msg{
		CommonMsg: deviceMsg.NewRespCommonMsg(l.dreq.Method, l.dreq.ClientToken),
	}
	resp.AddStatus(errors.OK)
	var (
		isOnline   = int64(def.False)
		payload    msgGateway.GatewayPayload
		appConnMsg = application.ConnectMsg{
			Device: devices.Core{
				ProductID:  msg.ProductID,
				DeviceName: msg.DeviceName,
			},
			Timestamp: msg.Timestamp.UnixMilli(),
		}
	)

	switch l.dreq.Method {
	case deviceMsg.Online:
		isOnline = def.True
		err = l.svcCtx.PubApp.DeviceStatusConnected(l.ctx, appConnMsg)
		if err != nil {
			l.Errorf("%s.DeviceStatusConnected productID:%v deviceName:%v err:%v",
				utils.FuncName(), msg.ProductID, msg.DeviceName, err)
		}
	case deviceMsg.Offline:
		err = l.svcCtx.PubApp.DeviceStatusDisConnected(l.ctx, appConnMsg)
		if err != nil {
			l.Errorf("%s.DeviceStatusDisConnected productID:%v deviceName:%v err:%v",
				utils.FuncName(), msg.ProductID, msg.DeviceName, err)
		}
	default:
		err := errors.Parameter.AddDetailf("not support method :%s", l.dreq.Method)
		resp.AddStatus(err)
		return &resp, err
	}
	for _, v := range l.dreq.Payload.Devices {
		//更新对应设备的online状态
		_, err := l.svcCtx.DeviceM.DeviceInfoUpdate(l.ctx, &dm.DeviceInfo{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
			IsOnline:   isOnline,
		})
		if err != nil {
			l.Errorf("%s.LogRepo.DeviceInfoUpdate productID:%v deviceName:%v err:%v",
				utils.FuncName(), v.ProductID, v.DeviceName, err)
		}
		payload.Devices = append(payload.Devices, &msgGateway.Device{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
			Result:     errors.Fmt(err).GetCode(),
		})
	}
	return &resp, err
}
