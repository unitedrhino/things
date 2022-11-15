package dataUpdateEvent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgGateway"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgSdkLog"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
)

type DataUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataUpdateLogic {
	return &DataUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (d *DataUpdateLogic) ProductSchemaUpdate(info *events.DataUpdateInfo) error {
	d.Infof("%s DataUpdateInfo:%v", utils.FuncName(), info)
	return d.svcCtx.SchemaRepo.ClearCache(d.ctx, info.ProductID)
}

func (d *DataUpdateLogic) DeviceLogLevelUpdate(info *events.DataUpdateInfo) error {
	d.Infof("%s DataUpdateInfo:%v", utils.FuncName(), info)
	di, err := d.svcCtx.DeviceM.DeviceInfoRead(d.ctx, &dm.DeviceInfoReadReq{
		ProductID:  info.ProductID,
		DeviceName: info.DeviceName,
	})
	if err != nil {
		return err
	}
	tmpTopic := fmt.Sprintf("%s/down/%s/%s/%s", devices.TopicHeadLog, msgSdkLog.TypeUpdate, di.ProductID, di.DeviceName)
	resp := deviceMsg.NewRespCommonMsg(deviceMsg.GetStatus, "")
	resp.Data = map[string]any{"logLevel": di.LogLevel}

	er := d.svcCtx.PubDev.PublishToDev(d.ctx, tmpTopic, resp.AddStatus(errors.OK).Bytes())
	if er != nil {
		d.Errorf("%s.PublishToDev failure err:%v", utils.FuncName(), er)
	}
	return er
}

func (d *DataUpdateLogic) DeviceGatewayUpdate(info *events.GatewayUpdateInfo) error {
	tmpTopic := fmt.Sprintf("%s/down/%s/%s/%s", devices.TopicHeadGateway, msgSdkLog.TypeOperation,
		info.GatewayProductID, info.GatewayDeviceName)
	resp := &msgGateway.Msg{
		CommonMsg: deviceMsg.NewRespCommonMsg(deviceMsg.Change, "").AddStatus(errors.OK),
		Payload:   ToGatewayPayload(info.Status, info.Devices),
	}
	respBytes, _ := json.Marshal(resp)
	er := d.svcCtx.PubDev.PublishToDev(d.ctx, tmpTopic, respBytes)
	if er != nil {
		d.Errorf("%s.PublishToDev failure err:%v", utils.FuncName(), er)
	}
	return er
}
