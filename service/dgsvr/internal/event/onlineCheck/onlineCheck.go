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
	logx.WithContext(o.ctx).Infof("online sync")
	var total int64 = 10000
	var limit int64 = 1000
	var page int64 = 1
	for page*limit < total {
		infos, to, err := o.svcCtx.MqttClient.GetOnlineClients(o.ctx, clients.GetOnlineClientsFilter{}, &def.PageInfo{
			Page: 1,
			Size: 200,
		})
		if err != nil {
			logx.WithContext(o.ctx).Error(err)
			return err
		}
		total = to
		page++
		var needOnlineDevices []*dm.DeviceOnlineMultiFix
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
			di, err := o.svcCtx.DeviceCache.GetData(o.ctx, devices.Core{
				ProductID:  dev.ProductID,
				DeviceName: dev.DeviceName,
			})
			if err != nil {
				continue
			}
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
		logx.WithContext(o.ctx).Infof("fixOnline %v", utils.Fmt(needOnlineDevices))
		_, err = o.svcCtx.DeviceM.DeviceOnlineMultiFix(o.ctx, &dm.DeviceOnlineMultiFixReq{Devices: needOnlineDevices})
		return err
	}
	return nil
}
