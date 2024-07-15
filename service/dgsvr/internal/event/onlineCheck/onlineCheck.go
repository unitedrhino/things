package onlineCheck

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/clients"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/sdk/service/protocol"
	"github.com/i-Things/things/service/dgsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
)

type CheckEvent struct {
	svcCtx *svc.ServiceContext
	logx.Logger
	ctx context.Context
}

func NewOnlineCheckEvent(svcCtx *svc.ServiceContext, ctx context.Context) *CheckEvent {
	return &CheckEvent{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (o *CheckEvent) Check() error {
	logx.WithContext(o.ctx).Infof("online_sync")
	var total int64 = 10000
	var limit int64 = 1000
	var page int64 = 1
	devs, err := protocol.GetActivityDevices(o.ctx)
	if err != nil {
		logx.WithContext(o.ctx).Error(err)
		devs = map[devices.Core]struct{}{}
	}
	var needOnlineDevices []*dm.DeviceOnlineMultiFix
	for page*limit < total {
		infos, to, err := o.svcCtx.MqttClient.GetOnlineClients(o.ctx, clients.GetOnlineClientsFilter{}, &clients.PageInfo{
			Page: page,
			Size: limit,
		})
		if err != nil {
			logx.WithContext(o.ctx).Error(err)
			return err
		}
		o.Infof("GetOnlineClients total:%v infos: %v ", total, utils.Fmt(infos))
		total = to
		page++
		for _, info := range infos {
			devStr, err := caches.GetStore().HgetCtx(o.ctx, protocol.DeviceMqttClientID, info.ClientID)
			if err != nil {
				continue
			}
			var dev devices.DevConn
			err = json.Unmarshal([]byte(devStr), &dev)
			if err != nil {
				continue
			}
			c := devices.Core{
				ProductID:  dev.ProductID,
				DeviceName: dev.DeviceName,
			}
			di, err := o.svcCtx.DeviceCache.GetData(o.ctx, c)
			if err != nil {
				continue
			}
			delete(devs, c)
			if di.IsOnline != def.True {
				needOnlineDevices = append(needOnlineDevices, &dm.DeviceOnlineMultiFix{
					Device: &dm.DeviceCore{
						ProductID:  di.ProductID,
						DeviceName: di.DeviceName,
					},
					IsOnline:  def.True,
					ConnectAt: info.Timestamp,
				})
			}
		}

	}

	if len(devs) > 0 { //如果全部过滤完了这里还有在线的,同时在emq上是离线的,那么需要下线该设备
		logx.WithContext(o.ctx).Infof("fixOffLine %v", utils.Fmt(devs))
		for dev := range devs {
			di, err := o.svcCtx.DeviceCache.GetData(o.ctx, dev)
			if err != nil || di.DeviceType == def.DeviceTypeSubset {
				continue
			}
			if di.IsOnline == def.True {
				needOnlineDevices = append(needOnlineDevices, &dm.DeviceOnlineMultiFix{
					Device: &dm.DeviceCore{
						ProductID:  di.ProductID,
						DeviceName: di.DeviceName,
					},
					IsOnline:  def.False,
					ConnectAt: 0,
				})
			} else {
				protocol.DeleteDeviceActivity(o.ctx, dev)
			}
		}
	}
	logx.WithContext(o.ctx).Infof("fixOnline %v", utils.Fmt(needOnlineDevices))
	if len(needOnlineDevices) > 0 {
		_, err = o.svcCtx.DeviceM.DeviceOnlineMultiFix(o.ctx, &dm.DeviceOnlineMultiFixReq{Devices: needOnlineDevices})
	}
	return err
}
