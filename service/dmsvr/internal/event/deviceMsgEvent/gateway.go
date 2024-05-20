package deviceMsgEvent

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/application"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/domain/deviceMsg/msgGateway"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/dmExport"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceLog"
	"github.com/i-Things/things/service/dmsvr/internal/repo/cache"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
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
	case msgGateway.TypeTopo:
		resp, err = l.HandleTopo(msg)
	case msgGateway.TypeStatus:
		resp, err = l.HandleStatus(msg)
	}
	respStr, _ := json.Marshal(resp)
	l.svcCtx.HubLogRepo.Insert(l.ctx, &deviceLog.Hub{
		ProductID:  msg.ProductID,
		Action:     "gateway",
		Timestamp:  time.Now(), // 记录当前时间
		DeviceName: msg.DeviceName,
		TraceID:    utils.TraceIdFromContext(l.ctx),
		RequestID:  l.dreq.MsgToken,
		Content:    string(msg.Payload),
		Topic:      msg.Topic,
		ResultCode: errors.Fmt(err).GetCode(),
	})
	return &deviceMsg.PublishMsg{
		Handle:       msg.Handle,
		Type:         msg.Type,
		Payload:      respStr,
		Timestamp:    time.Now().UnixMilli(),
		ProductID:    msg.ProductID,
		DeviceName:   msg.DeviceName,
		ProtocolCode: msg.ProtocolCode,
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
		_, err := l.svcCtx.DeviceCache.GetData(l.ctx, dmExport.GenDeviceInfoKey(v.ProductID, v.DeviceName))
		if err != nil {
			if errors.Cmp(err, errors.NotFind) {
				_, err = devicemanage.NewDeviceManageServer(l.svcCtx).DeviceInfoCreate(l.ctx, &dm.DeviceInfo{
					ProductID:  v.ProductID,
					DeviceName: v.DeviceName,
				})
			}
			if err != nil {
				l.Errorf("%s.DeviceM.DeviceInfoCreate productID:%v deviceName:%v err:%v",
					utils.FuncName(), v.ProductID, v.DeviceName, err)
				payload.Devices = append(payload.Devices, &msgGateway.Device{
					ProductID:  v.ProductID,
					DeviceName: v.DeviceName,
					Code:       errors.Fmt(err).GetCode(),
					Msg:        errors.Fmt(err).GetMsg(),
				})
				resp.AddStatus(err)
				continue
			}
		}
		di, err := l.svcCtx.DeviceCache.GetData(l.ctx, dmExport.GenDeviceInfoKey(v.ProductID, v.DeviceName))
		if err != nil {
			l.Errorf("%s.DeviceM.DeviceInfoRead productID:%v deviceName:%v err:%v",
				utils.FuncName(), v.ProductID, v.DeviceName, err)
			payload.Devices = append(payload.Devices, &msgGateway.Device{
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
				Code:       errors.Fmt(err).GetCode(),
				Msg:        errors.Fmt(err).GetMsg(),
			})
			resp.AddStatus(err)
			continue
		}
		c, err := relationDB.NewGatewayDeviceRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.GatewayDeviceFilter{
			SubDevice: &devices.Core{
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
			},
			SubDevices: nil,
		})
		if err == nil && !(c.GatewayProductID == msg.ProductID && c.GatewayDeviceName == msg.DeviceName) { //绑定了其他设备
			payload.Devices = append(payload.Devices, &msgGateway.Device{
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
				Code:       errors.Fmt(err).GetCode(),
				Msg:        errors.Fmt(err).GetMsg(),
			})
			resp.AddStatus(err)
			continue
		} else {
			if err != nil && !errors.Cmp(err, errors.NotFind) {
				resp.AddStatus(err)
				continue
			}
		}
		err = errors.OK
		payload.Devices = append(payload.Devices, &msgGateway.Device{
			ProductID:    v.ProductID,
			DeviceName:   v.DeviceName,
			DeviceSecret: di.GetSecret(),
			Code:         errors.Fmt(err).GetCode(),
			Msg:          errors.Fmt(err).GetMsg(),
		})
	}
	resp.Payload = &payload
	return resp, nil
}

func (l *GatewayLogic) HandleTopo(msg *deviceMsg.PublishMsg) (respMsg *msgGateway.Msg, err error) {
	l.Debugf("%s", utils.FuncName())
	var resp = msgGateway.Msg{
		CommonMsg: *deviceMsg.NewRespCommonMsg(l.ctx, l.dreq.Method, l.dreq.MsgToken),
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
				IsAuthSign: true,
				Gateway: &dm.DeviceCore{
					ProductID:  msg.ProductID,
					DeviceName: msg.DeviceName,
				},
				List: list,
			})
			if err != nil {
				resp.AddStatus(err)
				return &resp, err
			}
			resp.Payload = &msgGateway.GatewayPayload{Devices: l.dreq.Payload.Devices.GetCore()}
			return &resp, nil
		case deviceMsg.Unbind:
			_, err := devicemanage.NewDeviceManageServer(l.svcCtx).DeviceGatewayMultiDelete(l.ctx, &dm.DeviceGatewayMultiSaveReq{
				Gateway: &dm.DeviceCore{
					ProductID:  msg.ProductID,
					DeviceName: msg.DeviceName,
				},
				List: ToDmDevicesCore(l.dreq.Payload.Devices),
			})
			if err != nil {
				resp.AddStatus(err)
				return &resp, err
			}
			resp.Payload = &msgGateway.GatewayPayload{Devices: l.dreq.Payload.Devices.GetCore()}
		case deviceMsg.Found:
			//过滤已经入网的设备
			devs, err := relationDB.NewDeviceInfoRepo(l.ctx).FindByFilter(l.ctx,
				relationDB.DeviceFilter{Cores: l.dreq.Payload.Devices.GetDevCore()}, nil)
			if err != nil {
				resp.AddStatus(err)
				return &resp, err
			}
			var ca cache.GatewayCanBindStu
			ca.Gateway = devices.Core{
				ProductID:  msg.ProductID,
				DeviceName: msg.DeviceName,
			}
			ca.UpdatedTime = time.Now().Unix()
			for _, v := range devs {
				ca.SubDevices = append(ca.SubDevices, &devices.Core{
					ProductID:  v.ProductID,
					DeviceName: v.DeviceName,
				})
			}
			err = l.svcCtx.GatewayCanBind.Update(l.ctx, &ca)
			if err != nil {
				resp.AddStatus(err)
				return &resp, err
			}
			return &resp, nil
		case deviceMsg.GetTopo:
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
					Code:       errors.OK.Code,
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
		CommonMsg: *deviceMsg.NewRespCommonMsg(l.ctx, l.dreq.Method, l.dreq.MsgToken),
		Payload:   l.dreq.Payload,
	}
	resp.AddStatus(errors.OK)
	var (
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
	var subDevices []*devices.Core
	for _, v := range l.dreq.Payload.Devices {
		subDevices = append(subDevices, &devices.Core{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	gs, err := relationDB.NewGatewayDeviceRepo(l.ctx).CountByFilter(l.ctx, relationDB.GatewayDeviceFilter{SubDevices: subDevices, Gateway: &devices.Core{
		ProductID:  msg.ProductID,
		DeviceName: msg.DeviceName,
	}})
	if err != nil {
		resp.AddStatus(err)
		return &resp, err
	}
	if int(gs) != len(l.dreq.Payload.Devices) {
		err := errors.DeviceNotBound
		resp.AddStatus(err)
		return &resp, err
	}

	for _, v := range l.dreq.Payload.Devices {
		payload.Devices = append(payload.Devices, &msgGateway.Device{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
			Code:       errors.Fmt(err).GetCode(),
			Msg:        errors.Fmt(err).GetMsg(),
		})
	}
	if l.dreq.Method == deviceMsg.Online {
		for _, v := range l.dreq.Payload.Devices {
			di, err := relationDB.NewDeviceInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.DeviceFilter{Cores: []*devices.Core{{ProductID: v.ProductID, DeviceName: v.DeviceName}}})
			if err != nil {
				l.Error(err)
				continue
			}
			if di.IsOnline == def.True {
				l.Infof("already online:%#v", msg)
				continue
			}
			var updates = map[string]any{"is_online": def.True, "last_login": msg.Timestamp, "status": def.DeviceStatusOnline}
			if di.FirstLogin.Valid == false {
				updates["first_login"] = msg.Timestamp
			}
			err = relationDB.NewDeviceInfoRepo(l.ctx).UpdateWithField(l.ctx,
				relationDB.DeviceFilter{Cores: []*devices.Core{{ProductID: v.ProductID, DeviceName: v.DeviceName}}}, updates)
			if err != nil {
				l.Error(err)
			}
			err = l.svcCtx.DeviceCache.SetData(l.ctx, dmExport.GenDeviceInfoKey(v.ProductID, v.DeviceName), nil)
			if err != nil {
				l.Error(err)
			}
			//err = l.svcCtx.PubApp.DeviceStatusConnected(l.ctx, appMsg)
			//if err != nil {
			//	l.Errorf("%s.pubApp productID:%v deviceName:%v err:%v",
			//		utils.FuncName(), v.ProductID, v.DeviceName, err)
			//}
		}
	} else {
		var ( //这里是最后更新数据库状态的设备列表
			OffLineDevices []*devices.Core
		)
		for _, v := range l.dreq.Payload.Devices {
			OffLineDevices = append(OffLineDevices, &devices.Core{
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
			})
		}
		diDB := relationDB.NewDeviceInfoRepo(l.ctx)
		err = diDB.UpdateOfflineStatus(l.ctx, relationDB.DeviceFilter{Cores: OffLineDevices})
		if err != nil {
			l.Error(err)
		}
		for _, v := range OffLineDevices { //清除缓存
			err := l.svcCtx.DeviceCache.SetData(l.ctx, dmExport.GenDeviceInfoKey(v.ProductID, v.DeviceName), nil)
			if err != nil {
				l.Error(err)
			}
		}
	}
	return &resp, err
}
