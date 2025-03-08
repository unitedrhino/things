package onlineCheck

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/sdk/protocol"
	"gitee.com/unitedrhino/things/service/dgsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/clients"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/protocols"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/atomic"
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

var isRun atomic.Bool

func (o *CheckEvent) Check() error {
	logx.WithContext(o.ctx).Infof("online_sync")
	if !isRun.CompareAndSwap(false, true) {
		logx.WithContext(o.ctx).Infof("online_sync other run")
		return nil
	}
	defer isRun.Store(false)
	var total int64 = 10000
	var limit int64 = 500
	var page int64 = 0
	devs, err := protocol.GetActivityDevices(o.ctx)
	if err != nil {
		logx.WithContext(o.ctx).Error(err)
		devs = map[devices.Core]struct{}{}
	}
	var needOnlineDevices []*dm.DeviceOnlineMultiFix
	for page*limit < total {
		page++
		infos, to, err := o.svcCtx.MqttClient.GetOnlineClients(o.ctx, clients.GetOnlineClientsFilter{}, &clients.PageInfo{
			Page: page,
			Size: limit,
		})
		if err != nil {
			logx.WithContext(o.ctx).Error(err)
			return err
		}
		o.Infof("GetOnlineClients page:%v total:%v", page, total)
		total = to
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
			pi, err := o.svcCtx.ProductCache.GetData(o.ctx, di.ProductID)
			if err != nil {
				continue
			}
			if pi.Protocol != nil && pi.Protocol.TransProtocol != protocols.ProtocolMqtt {
				delete(devs, c)
				continue
			}
			if pi.SubProtocol != nil && pi.SubProtocol.TransProtocol != protocols.ProtocolMqtt {
				delete(devs, c)
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
			}
		}
	}
	logx.WithContext(o.ctx).Infof("fixOnline %v", utils.Fmt(needOnlineDevices))
	if len(needOnlineDevices) > 0 {
		_, err = o.svcCtx.DeviceM.DeviceOnlineMultiFix(o.ctx, &dm.DeviceOnlineMultiFixReq{Devices: needOnlineDevices})
	}
	return err
}
