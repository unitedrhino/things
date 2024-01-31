package deviceMsgEvent

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/devices"
	"gitee.com/i-Things/core/shared/domain/application"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg/msgGateway"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg/msgHubLog"
	devicemanage "github.com/i-Things/things/service/dmsvr/internal/server/devicemanage"
	productmanage "github.com/i-Things/things/service/dmsvr/internal/server/productmanage"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type GatewayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	dreq msgGateway.Msg
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

	switch msg.Type {
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
		RequestID:  l.dreq.MsgToken,
		Content:    string(msg.Payload),
		Topic:      msg.Topic,
		ResultType: errors.Fmt(err).GetCode(),
	})
	return &deviceMsg.PublishMsg{
		Handle:     msg.Handle,
		Type:       msg.Type,
		Payload:    respStr,
		Timestamp:  time.Now().UnixMilli(),
		ProductID:  msg.ProductID,
		DeviceName: msg.DeviceName,
	}, nil
}

func (l *GatewayLogic) HandleRegister(msg *deviceMsg.PublishMsg, resp *msgGateway.Msg) (respMsg *msgGateway.Msg, err error) {
	var (
		payload msgGateway.GatewayPayload
	)

	pds := l.dreq.Payload.Devices.GetProductIDs()
	pis, err := productmanage.NewProductManageServer(l.svcCtx).ProductInfoIndex(l.ctx, &dm.ProductInfoIndexReq{
		ProductIDs: pds,
	})
	if err != nil {
		er := errors.Fmt(err)
		resp.AddStatus(er)
		return resp, er
	}
	{ //参数检查
		if len(pis.List) != len(pds) {
			er := errors.Parameter.AddMsg("有产品id不正确,请查验")
			resp.AddStatus(er)
			return resp, er
		}
		for _, pi := range pis.List {
			if pi.AutoRegister != def.AutoRegAuto {
				er := errors.Parameter.AddMsgf("产品:%s 未打开自动注册", pi.ProductName)
				resp.AddStatus(er)
				return resp, er
			}
		}
	}
	for _, v := range l.dreq.Payload.Devices {
		//更新对应设备的online状态
		_, err := devicemanage.NewDeviceManageServer(l.svcCtx).DeviceInfoCreate(l.ctx, &dm.DeviceInfo{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
		if err != nil {
			l.Errorf("%s.DeviceM.DeviceInfoCreate productID:%v deviceName:%v err:%v",
				utils.FuncName(), v.ProductID, v.DeviceName, err)
			payload.Devices = append(payload.Devices, &msgGateway.Device{
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
				Result:     errors.Fmt(err).GetCode(),
				Status:     errors.Fmt(err).GetMsg(),
			})
			continue
		}
		di, err := devicemanage.NewDeviceManageServer(l.svcCtx).DeviceInfoRead(l.ctx, &dm.DeviceInfoReadReq{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
		if err != nil {
			l.Errorf("%s.DeviceM.DeviceInfoRead productID:%v deviceName:%v err:%v",
				utils.FuncName(), v.ProductID, v.DeviceName, err)
		}
		payload.Devices = append(payload.Devices, &msgGateway.Device{
			ProductID:    v.ProductID,
			DeviceName:   v.DeviceName,
			DeviceSecret: di.GetSecret(),
			Result:       errors.Fmt(err).GetCode(),
			Status:       errors.Fmt(err).GetMsg(),
		})
	}
	resp.Payload = &payload
	return resp, nil
}

func (l *GatewayLogic) HandleOperation(msg *deviceMsg.PublishMsg) (respMsg *msgGateway.Msg, err error) {
	l.Debugf("%s", utils.FuncName())
	var resp = msgGateway.Msg{
		CommonMsg: deviceMsg.NewRespCommonMsg(l.ctx, l.dreq.Method, l.dreq.MsgToken),
	}
	resp.AddStatus(errors.OK)
	rsp, err := func() (respMsg *msgGateway.Msg, err error) {
		switch l.dreq.Method {
		case deviceMsg.Register:
			return l.HandleRegister(msg, &resp)
		case deviceMsg.Bind:
			list, err := ToDmDevicesBind(l.dreq.Payload.Devices)
			if err != nil {
				resp.AddStatus(err)
				return &resp, err
			}
			_, err = devicemanage.NewDeviceManageServer(l.svcCtx).DeviceGatewayMultiCreate(l.ctx, &dm.DeviceGatewayMultiCreateReq{
				IsAuthSign:        true,
				GatewayProductID:  msg.ProductID,
				GatewayDeviceName: msg.DeviceName,
				List:              list,
			})
			if err != nil {
				resp.AddStatus(err)
				return &resp, err
			}
			resp.Payload = &msgGateway.GatewayPayload{Devices: l.dreq.Payload.Devices.GetCore()}
			return &resp, nil
		case deviceMsg.Unbind:
			_, err := devicemanage.NewDeviceManageServer(l.svcCtx).DeviceGatewayMultiDelete(l.ctx, &dm.DeviceGatewayMultiDeleteReq{
				GatewayProductID:  msg.ProductID,
				GatewayDeviceName: msg.DeviceName,
				List:              ToDmDevicesCore(l.dreq.Payload.Devices),
			})
			if err != nil {
				resp.AddStatus(err)
				return &resp, err
			}
			resp.Payload = &msgGateway.GatewayPayload{Devices: l.dreq.Payload.Devices.GetCore()}
		case deviceMsg.DescribeSubDevices:
			deviceList, err := devicemanage.NewDeviceManageServer(l.svcCtx).DeviceGatewayIndex(l.ctx, &dm.DeviceGatewayIndexReq{
				Gateway: &dm.DeviceCore{
					ProductID:  msg.ProductID,
					DeviceName: msg.DeviceName,
				}})
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
			return &resp, err
		default:
			return nil, errors.Method.AddMsg(l.dreq.Method)
		}
		return nil, errors.Parameter.AddDetailf("gateway types is err:%v", msg.Type)
	}()
	if l.dreq.NoAsk() { //如果不需要回复
		rsp = nil
	}
	return rsp, err
}

func (l *GatewayLogic) HandleStatus(msg *deviceMsg.PublishMsg) (respMsg *msgGateway.Msg, err error) {
	l.Debugf("%s", utils.FuncName())
	var resp = msgGateway.Msg{
		CommonMsg: deviceMsg.NewRespCommonMsg(l.ctx, l.dreq.Method, l.dreq.MsgToken),
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
			Timestamp: msg.Timestamp,
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
		_, err := devicemanage.NewDeviceManageServer(l.svcCtx).DeviceInfoUpdate(l.ctx, &dm.DeviceInfo{
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
			Status:     errors.Fmt(err).GetMsg(),
		})
	}
	return &resp, err
}
