package deviceMsgEvent

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/domain/deviceAuth"
	"gitee.com/unitedrhino/share/domain/deviceMsg"
	"gitee.com/unitedrhino/share/domain/deviceMsg/msgGateway"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceStatus"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/product"
	devicemanagelogic "gitee.com/unitedrhino/things/service/dmsvr/internal/logic/devicemanage"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/cache"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	devicemanage "gitee.com/unitedrhino/things/service/dmsvr/internal/server/devicemanage"
	productmanage "gitee.com/unitedrhino/things/service/dmsvr/internal/server/productmanage"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
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

func (l *GatewayLogic) DeviceResp(msg *deviceMsg.PublishMsg, err error, data any) *deviceMsg.PublishMsg {
	if !errors.Cmp(err, errors.OK) {
		l.Errorf("%s.DeviceResp err:%v, msg:%v", utils.FuncName(), err, msg)
	}
	resp := &deviceMsg.CommonMsg{
		Method:   deviceMsg.GetRespMethod(l.dreq.Method),
		MsgToken: l.dreq.MsgToken,
		//Timestamp: time.Now().UnixMilli(),
		Data: data,
	}
	if msg.ProtocolCode == "" {
		msg.ProtocolCode = def.ProtocolCodeUnitedRhino
	}
	return &deviceMsg.PublishMsg{
		Handle:       msg.Handle,
		Type:         msg.Type,
		Payload:      resp.AddStatus(err, l.dreq.NeedRetMsg()).Bytes(),
		Timestamp:    time.Now().UnixMilli(),
		ProductID:    msg.ProductID,
		DeviceName:   msg.DeviceName,
		ProtocolCode: msg.ProtocolCode,
	}
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
	case msgGateway.TypeThing:
		resp, err = l.HandleThing(msg)
	}
	respStr, _ := json.Marshal(resp)
	hub := &deviceLog.Hub{
		ProductID:   msg.ProductID,
		Action:      deviceLog.ActionTypeGateway,
		Timestamp:   time.Now(), // 记录当前时间
		DeviceName:  msg.DeviceName,
		TraceID:     utils.TraceIdFromContext(l.ctx),
		RequestID:   l.dreq.MsgToken,
		Content:     string(msg.Payload),
		Topic:       msg.Topic,
		ResultCode:  errors.Fmt(err).GetCode(),
		RespPayload: respMsg.GetPayload(),
	}
	l.svcCtx.HubLogRepo.Insert(l.ctx, hub)
	l.svcCtx.UserSubscribe.Publish(l.ctx, def.UserSubscribeDevicePublish, hub.ToApp(), map[string]any{
		"productID":  msg.ProductID,
		"deviceName": msg.DeviceName,
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
		resp.AddStatus(er, l.dreq.NeedRetMsg())
		return resp, er
	}
	{ //参数检查
		if len(pis.List) != len(pds) {
			er := errors.Parameter.AddMsg("有产品id不正确,请查验")
			resp.AddStatus(er, l.dreq.NeedRetMsg())
			return resp, er
		}
		for _, pi := range pis.List {
			if pi.AutoRegister != def.AutoRegAuto {
				er := errors.Parameter.AddMsgf("产品:%s 未打开自动注册", pi.ProductName)
				resp.AddStatus(er, l.dreq.NeedRetMsg())
				return resp, er
			}
		}
	}
	for _, v := range l.dreq.Payload.Devices {
		di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
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
				resp.AddStatus(err, l.dreq.NeedRetMsg())
				continue
			}
		}
		if di == nil {
			di, err = l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
			})
			if err != nil {
				resp.AddStatus(err, l.dreq.NeedRetMsg())
				continue
			}
		}
		c, err := relationDB.NewGatewayDeviceRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.GatewayDeviceFilter{
			SubDevice: &devices.Core{
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
			},
		})
		if err == nil && !(c.GatewayProductID == msg.ProductID && c.GatewayDeviceName == msg.DeviceName) { //绑定了其他设备
			err = errors.DeviceBound
			payload.Devices = append(payload.Devices, &msgGateway.Device{
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
				Code:       errors.Fmt(err).GetCode(),
				Msg:        errors.Fmt(err).GetMsg(),
			})
			resp.AddStatus(err, l.dreq.NeedRetMsg())
			continue
		} else {
			if err != nil && !errors.Cmp(err, errors.NotFind) {
				resp.AddStatus(err, l.dreq.NeedRetMsg())
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
	resp.AddStatus(errors.OK, l.dreq.NeedRetMsg())
	rsp, err := func() (respMsg *msgGateway.Msg, err error) {
		switch l.dreq.Method {
		case deviceMsg.Register:
			return l.HandleRegister(msg, &resp)
		case deviceMsg.Bind:
			list, err := ToDmDevicesBind(l.dreq.Payload.Devices)
			if err != nil {
				resp.AddStatus(err, l.dreq.NeedRetMsg())
				return &resp, err
			}
			_, err = devicemanage.NewDeviceManageServer(l.svcCtx).DeviceGatewayMultiCreate(l.ctx, &dm.DeviceGatewayMultiCreateReq{
				IsAuthSign:  true,
				IsNotNotify: true,
				Gateway: &dm.DeviceCore{
					ProductID:  msg.ProductID,
					DeviceName: msg.DeviceName,
				},
				List: list,
			})
			if err != nil {
				resp.AddStatus(err, l.dreq.NeedRetMsg())
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
				IsNotNotify: true,
				List:        ToDmDevicesCore(l.dreq.Payload.Devices),
			})
			if err != nil {
				resp.AddStatus(err, l.dreq.NeedRetMsg())
				return &resp, err
			}
			resp.Payload = &msgGateway.GatewayPayload{Devices: l.dreq.Payload.Devices.GetCore()}
			return &resp, nil
		case deviceMsg.Found:
			var devs []*devices.Core
			devs, err = devicemanagelogic.FilterCanBindSubDevices(l.ctx, l.svcCtx, &devices.Core{
				ProductID:  msg.ProductID,
				DeviceName: msg.DeviceName,
			}, l.dreq.Payload.Devices.GetDevCore(), devicemanagelogic.CheckDeviceType)
			if err != nil {
				resp.AddStatus(err, l.dreq.NeedRetMsg())
				return &resp, err
			}
			var ca = cache.GatewayCanBindStu{
				Gateway: devices.Core{
					ProductID:  msg.ProductID,
					DeviceName: msg.DeviceName,
				},
				SubDevices:  devs,
				UpdatedTime: time.Now().Unix(),
			}
			err = l.svcCtx.GatewayCanBind.Update(l.ctx, &ca)
			if err != nil {
				resp.AddStatus(err, l.dreq.NeedRetMsg())
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
				resp.AddStatus(err, l.dreq.NeedRetMsg())
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
		case deviceMsg.GetFound:
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

var (
	ActionMap = map[string]string{
		deviceMsg.Online:  devices.ActionConnected,
		deviceMsg.Offline: devices.ActionDisconnected,
	}
)

func (l *GatewayLogic) HandleThing(msg *deviceMsg.PublishMsg) (respMsg *msgGateway.Msg, err error) {
	l.Infof("%s req=%+v", utils.FuncName(), msg)
	err = l.initMsg(msg)
	if err != nil {
		return nil, err
	}
	var resp = &msgGateway.Msg{
		CommonMsg: *deviceMsg.NewRespCommonMsg(l.ctx, l.dreq.Method, l.dreq.MsgToken),
	}
	resp.AddStatus(errors.OK, l.dreq.NeedRetMsg())
	p := l.dreq.Payload
	if p != nil && p.ProductID != "" && p.DeviceName != "" && p.DeviceName != msg.DeviceName && p.ProductID != msg.ProductID {
		gs, err := relationDB.NewGatewayDeviceRepo(l.ctx).CountByFilter(l.ctx,
			relationDB.GatewayDeviceFilter{SubDevices: []*devices.Core{{ProductID: p.ProductID, DeviceName: p.DeviceName}},
				Gateway: &devices.Core{
					ProductID:  msg.ProductID,
					DeviceName: msg.DeviceName,
				}})
		if err != nil {
			resp.AddStatus(err, l.dreq.NeedRetMsg())
			return resp, err
		}
		if gs == 0 {
			err := errors.DeviceNotBound
			resp.AddStatus(err, l.dreq.NeedRetMsg())
			return resp, err
		}
	}
	switch l.dreq.Method {
	case deviceMsg.CreateSchema:
		resp, err = l.HandleCreateSchema(msg, resp)
	case deviceMsg.DeleteSchema:
		resp, err = l.HandleDeleteSchema(msg, resp)
	case deviceMsg.GetSchema:
		resp, err = l.HandlePropertyGetSchema(msg, resp)
	}
	if l.dreq.NoAsk() { //如果不需要回复
		resp = nil
	}
	return resp, err
}

func (l *GatewayLogic) HandlePropertyGetSchema(msg *deviceMsg.PublishMsg, resp *msgGateway.Msg) (respMsg *msgGateway.Msg, err error) {
	var (
		payload msgGateway.GatewayPayload
	)
	if l.dreq.Payload == nil || l.dreq.Payload.ProductID == "" { //如果没有传产品,则会返回设备物模型
		s, err := l.svcCtx.DeviceSchemaRepo.GetData(l.ctx, devices.Core{ProductID: msg.ProductID, DeviceName: msg.DeviceName})
		if err != nil {
			resp.AddStatus(err, l.dreq.NeedRetMsg())
			return resp, err
		}
		payload.Schema = s.ToSimple()
		payload.ProductID = msg.ProductID
		payload.DeviceName = msg.DeviceName
		resp.Payload = &payload
		return resp, nil
	}
	if l.dreq.Payload.DeviceName == "" {
		s, err := l.svcCtx.ProductSchemaRepo.GetData(l.ctx, devices.Core{ProductID: l.dreq.Payload.ProductID})
		if err != nil {
			resp.AddStatus(err, l.dreq.NeedRetMsg())
			return resp, err
		}
		payload.Schema = s.ToSimple()
		payload.ProductID = l.dreq.Payload.ProductID
		resp.Payload = &payload
		return resp, nil
	}
	s, err := l.svcCtx.DeviceSchemaRepo.GetData(l.ctx, devices.Core{ProductID: l.dreq.Payload.ProductID, DeviceName: l.dreq.Payload.DeviceName})
	if err != nil {
		resp.AddStatus(err, l.dreq.NeedRetMsg())
		return resp, err
	}
	payload.ProductID = l.dreq.Payload.ProductID
	payload.DeviceName = l.dreq.Payload.DeviceName
	payload.Schema = s.ToSimple()
	resp.Payload = &payload
	return resp, nil
}

func (l *GatewayLogic) HandleCreateSchema(msg *deviceMsg.PublishMsg, resp *msgGateway.Msg) (respMsg *msgGateway.Msg, err error) {
	if l.dreq.Payload == nil || l.dreq.Payload.Schema == nil {
		er := errors.Parameter.AddMsg("需要填写schema")
		resp.AddStatus(er, l.dreq.NeedRetMsg())
		return resp, er
	}
	if l.dreq.Payload.ProductID == "" || l.dreq.Payload.DeviceName == "" {
		l.dreq.Payload.ProductID = msg.ProductID
		l.dreq.Payload.DeviceName = msg.DeviceName
	}
	pi, err := l.svcCtx.ProductCache.GetData(l.ctx, msg.ProductID)
	if err != nil {
		resp.AddStatus(err, l.dreq.NeedRetMsg())
		return resp, err
	}
	if pi.DeviceSchemaMode < product.DeviceSchemaModeAutoCreate {
		er := errors.Permissions.AddMsg("产品未开启设备自动创建")
		resp.AddStatus(er, l.dreq.NeedRetMsg())
		return resp, er
	}
	m := l.dreq.Payload.Schema.ToModel()
	err = m.ValidateWithFmt()
	if err != nil {
		resp.AddStatus(err, l.dreq.NeedRetMsg())
		return resp, err
	}
	err = relationDB.NewDeviceSchemaRepo(l.ctx).MultiInsert2(l.ctx, msg.ProductID, msg.DeviceName, m)
	if err != nil {
		resp.AddStatus(err, l.dreq.NeedRetMsg())
		return resp, err
	}
	l.svcCtx.DeviceSchemaRepo.SetData(l.ctx, devices.Core{
		ProductID:  msg.ProductID,
		DeviceName: msg.DeviceName,
	}, nil)
	return resp, nil
}
func (l *GatewayLogic) HandleDeleteSchema(msg *deviceMsg.PublishMsg, resp *msgGateway.Msg) (respMsg *msgGateway.Msg, err error) {
	if l.dreq.Payload == nil || len(l.dreq.Payload.Identifiers) == 0 {
		er := errors.Parameter.AddMsg("需要填写identifiers")
		resp.AddStatus(er, l.dreq.NeedRetMsg())
		return resp, er
	}
	pi, err := l.svcCtx.ProductCache.GetData(l.ctx, msg.ProductID)
	if err != nil {
		resp.AddStatus(err, l.dreq.NeedRetMsg())
		return resp, err
	}
	if pi.DeviceSchemaMode < product.DeviceSchemaModeAutoCreate {
		er := errors.Permissions.AddMsg("产品未开启设备自动创建")
		resp.AddStatus(er, l.dreq.NeedRetMsg())
		return resp, er
	}
	_, err = devicemanagelogic.NewDeviceSchemaMultiDeleteLogic(l.ctx, l.svcCtx).DeviceSchemaMultiDelete(&dm.DeviceSchemaMultiDeleteReq{
		ProductID:   msg.ProductID,
		DeviceName:  msg.DeviceName,
		Identifiers: l.dreq.Payload.Identifiers,
	})
	if err != nil {
		resp.AddStatus(err, l.dreq.NeedRetMsg())
		return resp, err
	}
	l.svcCtx.DeviceSchemaRepo.SetData(l.ctx, devices.Core{
		ProductID:  msg.ProductID,
		DeviceName: msg.DeviceName,
	}, nil)
	return resp, nil
}

func (l *GatewayLogic) HandleStatus(msg *deviceMsg.PublishMsg) (respMsg *msgGateway.Msg, err error) {
	l.Debugf("%s", utils.FuncName())

	var resp = msgGateway.Msg{
		CommonMsg: *deviceMsg.NewRespCommonMsg(l.ctx, l.dreq.Method, l.dreq.MsgToken),
		Payload:   l.dreq.Payload,
	}
	resp.AddStatus(errors.OK, l.dreq.NeedRetMsg())
	if !utils.SliceIn(l.dreq.Method, deviceMsg.Offline, deviceMsg.Online) {
		err = errors.Parameter.AddMsg("method not support")
		resp.AddStatus(err, l.dreq.NeedRetMsg())
		return &resp, err
	}
	var (
		payload msgGateway.GatewayPayload
	)
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
		resp.AddStatus(err, l.dreq.NeedRetMsg())
		return &resp, err
	}
	if int(gs) != len(l.dreq.Payload.Devices) {
		err := errors.DeviceNotBound
		resp.AddStatus(err, l.dreq.NeedRetMsg())
		return &resp, err
	}

	for _, v := range l.dreq.Payload.Devices {
		payload.Devices = append(payload.Devices, &msgGateway.Device{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
		//更新在线状态
		err := devicemanagelogic.HandleOnlineFix(l.ctx, l.svcCtx, &deviceStatus.ConnectMsg{
			ClientID:  deviceAuth.GenClientID(v.ProductID, v.DeviceName),
			Timestamp: l.dreq.GetTimeStamp(),
			Action:    ActionMap[l.dreq.Method],
			Reason:    "gateway report",
		})
		if err != nil {
			l.Error(err)
		}
	}
	return &resp, err
}
