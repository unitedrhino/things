package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/sysExport"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/sdk/protocol"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceStatus"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/application"
	"gitee.com/unitedrhino/things/share/domain/deviceAuth"
	"github.com/spf13/cast"
	"time"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceOnlineMultiFixLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceOnlineMultiFixLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceOnlineMultiFixLogic {
	return &DeviceOnlineMultiFixLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceOnlineMultiFixLogic) DeviceOnlineMultiFix(in *dm.DeviceOnlineMultiFixReq) (*dm.Empty, error) {
	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
		var insertList []*deviceStatus.ConnectMsg
		for _, device := range in.Devices {
			ld := device.Device
			t := time.Now()
			if device.ConnectAt != 0 {
				t = time.UnixMilli(device.ConnectAt)
			}
			action := devices.ActionConnected
			if device.IsOnline == def.False {
				action = devices.ActionDisconnected
			}
			insertList = append(insertList, &deviceStatus.ConnectMsg{
				Device: devices.Core{
					ProductID:  ld.ProductID,
					DeviceName: ld.DeviceName,
				},
				Timestamp: t,
				Action:    action,
				Reason:    "在线状态修复",
			})
		}
		err := HandleOnlineFix(ctx, l.svcCtx, insertList...)
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
	})

	return &dm.Empty{}, nil
}

func HandleOnlineFix(ctx context.Context, svcCtx *svc.ServiceContext, insertList ...*deviceStatus.ConnectMsg) (err error) {
	ctx = ctxs.WithRoot(ctx)
	var ( //这里是最后更新数据库状态的设备列表
		OffLineDevices  []*devices.Core
		subDeviceInsert []*deviceStatus.ConnectMsg
	)

	var log = logx.WithContext(ctx)

	handleMsg := func(msg *deviceStatus.ConnectMsg) {
		status := int64(def.ConnectedStatus)
		if msg.Action == devices.ActionDisconnected {
			status = def.DisConnectedStatus
		}
		var ld *deviceAuth.LoginDevice
		if msg.Device.DeviceName != "" {
			ld = &deviceAuth.LoginDevice{
				ProductID:  msg.Device.ProductID,
				DeviceName: msg.Device.DeviceName,
			}
		} else {
			ld, err = deviceAuth.GetClientIDInfo(msg.ClientID)
			if err != nil {
				log.Error(err)
				return
			}
		}
		appMsg := application.ConnectMsg{
			Device: devices.Core{
				ProductID:  ld.ProductID,
				DeviceName: ld.DeviceName,
			},
			Status:    status,
			Timestamp: msg.Timestamp.UnixMilli(),
		}

		dev := devices.Core{
			ProductID:  ld.ProductID,
			DeviceName: ld.DeviceName,
		}
		di, err := svcCtx.DeviceCache.GetData(ctx, dev)
		if err != nil {
			log.Error(err)
			return
		}
		push := func(appMsg application.ConnectMsg, di *dm.DeviceInfo) {
			err = svcCtx.StatusRepo.Insert(ctx, &deviceLog.Status{
				ProductID:  appMsg.Device.ProductID,
				Status:     status,
				Timestamp:  msg.Timestamp, // 操作时间
				DeviceName: appMsg.Device.DeviceName,
			})
			if err != nil {
				log.Errorf("%s.HubLogRepo.insert productID:%v deviceName:%v err:%v",
					utils.FuncName(), ld.ProductID, ld.DeviceName, err)
			}

			err = svcCtx.PubApp.DeviceStatusDisConnected(ctx, appMsg)
			if err != nil {
				log.Errorf("%s.pubApp productID:%v deviceName:%v err:%v",
					utils.FuncName(), ld.ProductID, ld.DeviceName, err)
			}
			err = svcCtx.WebHook.Publish(svcCtx.WithDeviceTenant(ctx, appMsg.Device), func() string {
				if status == def.ConnectedStatus {
					return sysExport.CodeDmDeviceConn
				}
				return sysExport.CodeDmDeviceDisConn
			}(), appMsg)
			if err != nil {
				log.Error(err)
			}
			if di == nil {
				di, err = svcCtx.DeviceCache.GetData(ctx, dev)
				if err != nil {
					log.Error(err)
					return
				}
			}
			err = svcCtx.UserSubscribe.Publish(ctx, def.UserSubscribeDeviceConn, appMsg, map[string]any{
				"productID":  appMsg.Device.ProductID,
				"deviceName": appMsg.Device.DeviceName,
			}, map[string]any{
				"projectID": di.ProjectID,
			}, map[string]any{
				"projectID": cast.ToString(di.ProjectID),
				"areaID":    cast.ToString(di.AreaID),
			})
			if err != nil {
				log.Error(err)
			}
		}
		if status == def.ConnectedStatus {
			var updates = map[string]any{"is_online": def.True, "last_login": msg.Timestamp, "status": def.DeviceStatusOnline, "last_ip": msg.Address}
			if di.FirstLogin == 0 {
				updates["first_login"] = msg.Timestamp
			}
			err = relationDB.NewDeviceInfoRepo(ctx).UpdateWithField(ctx,
				relationDB.DeviceFilter{Cores: []*devices.Core{&dev}}, updates)
			if err != nil {
				log.Error(err)
			}
			err = svcCtx.DeviceCache.SetData(ctx, dev, nil)
			if err != nil {
				log.Error(err)
			}
			if di.IsOnline != def.True {
				push(appMsg, di)
			}
			protocol.UpdateDeviceActivity(ctx, dev)
		} else {
			if !utils.SliceIn(msg.Reason, "takeovered", "takenover", "discard", "discarded") { //连接还在的时候被别人顶了,忽略这种下线
				if di.DeviceType == def.DeviceTypeGateway { //如果是网关类型下线,则需要把子设备全部下线
					subDevs, err := relationDB.NewGatewayDeviceRepo(ctx).FindByFilter(ctx,
						relationDB.GatewayDeviceFilter{Gateway: &dev}, nil)
					if err != nil {
						log.Error(err)
					} else {
						for _, v := range subDevs {
							app := appMsg
							app.Device = devices.Core{ProductID: v.ProductID, DeviceName: v.DeviceName}
							push(appMsg, nil)
							subDeviceInsert = append(subDeviceInsert, &deviceStatus.ConnectMsg{Action: msg.Action,
								Device: devices.Core{ProductID: v.ProductID, DeviceName: v.DeviceName}})
						}
					}
				}
				OffLineDevices = append(OffLineDevices, &dev)
				if di.IsOnline == def.True {
					push(appMsg, di)
				}
				protocol.DeleteDeviceActivity(ctx, dev)
			}
		}

	}
	for _, msg := range insertList {
		handleMsg(msg)
	}
	if len(subDeviceInsert) != 0 {
		//子设备下线
		log.Infof("子设备下线: %v", utils.Fmt(subDeviceInsert))
		for _, msg := range subDeviceInsert {
			handleMsg(msg)
		}
	}
	diDB := relationDB.NewDeviceInfoRepo(ctx)
	if len(OffLineDevices) > 0 {
		err = diDB.UpdateOfflineStatus(ctx, relationDB.DeviceFilter{Cores: OffLineDevices})
		if err != nil {
			log.Error(err)
		}
		for _, v := range OffLineDevices { //清除缓存
			err := svcCtx.DeviceCache.SetData(ctx, *v, nil)
			if err != nil {
				log.Error(err)
			}
		}
	}
	return nil
}
