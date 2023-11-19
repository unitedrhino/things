package dataUpdateEvent

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg/msgGateway"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg/msgRemoteConfig"
	devicemanage "github.com/i-Things/things/src/dmsvr/internal/server/devicemanage"
	remoteconfig "github.com/i-Things/things/src/dmsvr/internal/server/remoteconfig"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type DataUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDataUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataUpdateLogic {
	return &DataUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (d *DataUpdateLogic) ProductSchemaUpdate(info *events.DeviceUpdateInfo) error {
	d.Infof("%s DeviceUpdateInfo:%v", utils.FuncName(), info)
	return d.svcCtx.SchemaRepo.ClearCache(d.ctx, info.ProductID)
}

func (d *DataUpdateLogic) DeviceGatewayUpdate(info *events.GatewayUpdateInfo) error {
	resp := &msgGateway.Msg{
		CommonMsg: deviceMsg.NewRespCommonMsg(d.ctx, deviceMsg.Change, "").AddStatus(errors.OK),
		Payload:   ToGatewayPayload(info.Status, info.Devices),
	}
	respBytes, _ := json.Marshal(resp)
	msg := deviceMsg.PublishMsg{
		Handle:     devices.Gateway,
		Type:       msgGateway.TypeOperation,
		Payload:    respBytes,
		Timestamp:  time.Now().UnixMilli(),
		ProductID:  info.GatewayProductID,
		DeviceName: info.GatewayDeviceName,
	}
	er := d.svcCtx.PubDev.PublishToDev(d.ctx, &msg)
	if er != nil {
		d.Errorf("%s.PublishToDev failure err:%v", utils.FuncName(), er)
	}
	return er
}

func (d *DataUpdateLogic) DeviceRemoteConfigUpdate(info *events.DeviceUpdateInfo) error {
	d.Infof("%s DeviceRemoteConfigUpdate:%v", utils.FuncName(), info)

	//1. 根据产品id获取配置json
	respConfig, err := remoteconfig.NewRemoteConfigServer(d.svcCtx).RemoteConfigLastRead(d.ctx, &dm.RemoteConfigLastReadReq{
		ProductID: info.ProductID,
	})
	if err != nil {
		d.Errorf("%s.RemoteConfigLastRead failure err:%v", utils.FuncName(), err)
		return err
	}

	//2. 根据产品id获取产品下的所有设备信息
	respDevices, err := devicemanage.NewDeviceManageServer(d.svcCtx).DeviceInfoIndex(d.ctx, &dm.DeviceInfoIndexReq{
		ProductID: info.ProductID,
	})
	if err != nil {
		d.Errorf("%s.RemoteConfigLastRead failure err:%v", utils.FuncName(), err)
		return err
	}

	//3. for循环所有设备发送消息给设备
	for _, v := range respDevices.List {
		resp := &msgRemoteConfig.RemoteConfigMsg{
			Method:  "push",
			Code:    0,
			Payload: respConfig.Info.Content,
		}
		respBytes, _ := json.Marshal(resp)
		msg := deviceMsg.PublishMsg{
			Handle:     devices.Config,
			Type:       msgRemoteConfig.TypePush,
			Payload:    respBytes,
			Timestamp:  time.Now().UnixMilli(),
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		}
		er := d.svcCtx.PubDev.PublishToDev(d.ctx, &msg)
		if er != nil {
			d.Errorf("%s.PublishToDev failure err:%v", utils.FuncName(), er)
		}
	}

	return err
}
