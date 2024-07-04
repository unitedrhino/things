package devicemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/sysExport"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/application"
	"gitee.com/i-Things/share/domain/deviceAuth"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/sdk/service/protocol"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceLog"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceStatus"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"time"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

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
				Device: &devices.Core{
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
		if msg.Device != nil {
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
		err = svcCtx.StatusRepo.Insert(ctx, &deviceLog.Status{
			ProductID:  ld.ProductID,
			Status:     status,
			Timestamp:  msg.Timestamp, // 操作时间
			DeviceName: ld.DeviceName,
		})
		if err != nil {
			log.Errorf("%s.HubLogRepo.insert productID:%v deviceName:%v err:%v",
				utils.FuncName(), ld.ProductID, ld.DeviceName, err)
		}
		appMsg := application.ConnectMsg{
			Device: devices.Core{
				ProductID:  ld.ProductID,
				DeviceName: ld.DeviceName,
			},
			Status:    status,
			Timestamp: msg.Timestamp.UnixMilli(),
		}
		utils.Go(ctx, func() {
			err = svcCtx.WebHook.Publish(svcCtx.WithDeviceTenant(ctx, appMsg.Device), sysExport.CodeDmDeviceConn, appMsg)
			if err != nil {
				log.Error(err)
			}
			err = svcCtx.UserSubscribe.Publish(ctx, def.UserSubscribeDeviceConn, appMsg, map[string]any{
				"productID":  ld.ProductID,
				"deviceName": ld.DeviceName,
			}, map[string]any{
				"productID": ld.ProductID,
			}, map[string]any{})
			if err != nil {
				log.Error(err)
			}
		})
		dev := devices.Core{
			ProductID:  ld.ProductID,
			DeviceName: ld.DeviceName,
		}
		if status == def.ConnectedStatus {
			di, err := relationDB.NewDeviceInfoRepo(ctx).FindOneByFilter(ctx, relationDB.DeviceFilter{Cores: []*devices.Core{&dev}})
			if err != nil {
				log.Error(err)
				return
			}
			if di.IsOnline == def.True {
				log.Infof("already online:%#v", msg)
				return
			}
			var updates = map[string]any{"is_online": def.True, "last_login": msg.Timestamp, "status": def.DeviceStatusOnline}
			if di.FirstLogin.Valid == false {
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

			err = svcCtx.PubApp.DeviceStatusConnected(ctx, appMsg)
			if err != nil {
				log.Errorf("%s.pubApp productID:%v deviceName:%v err:%v",
					utils.FuncName(), ld.ProductID, ld.DeviceName, err)
			}
			protocol.UpdateDeviceActivity(ctx, dev)
		} else {
			di, err := svcCtx.DeviceCache.GetData(ctx, dev)
			if err != nil {
				log.Error(err)
			} else if di.DeviceType == def.DeviceTypeGateway { //如果是网关类型下线,则需要把子设备全部下线
				subDevs, err := relationDB.NewGatewayDeviceRepo(ctx).FindByFilter(ctx,
					relationDB.GatewayDeviceFilter{Gateway: &dev}, nil)
				if err != nil {
					log.Error(err)
				} else {
					for _, v := range subDevs {
						subDeviceInsert = append(subDeviceInsert, &deviceStatus.ConnectMsg{Action: msg.Action,
							Device: &devices.Core{ProductID: v.ProductID, DeviceName: v.DeviceName}})
					}
				}
			}
			OffLineDevices = append(OffLineDevices, &dev)
			err = svcCtx.PubApp.DeviceStatusDisConnected(ctx, appMsg)
			if err != nil {
				log.Errorf("%s.pubApp productID:%v deviceName:%v err:%v",
					utils.FuncName(), ld.ProductID, ld.DeviceName, err)
			}
			protocol.DeleteDeviceActivity(ctx, dev)
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
