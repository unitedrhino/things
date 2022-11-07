package dataUpdateEvent

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/service/deviceSend"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
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
	uuid, _ := uuid.GenerateUUID()
	tmpTopic := fmt.Sprintf("%s/down/update/%s/%s", devices.TopicHeadLog, di.ProductID, di.DeviceName)
	topic, payload := deviceSend.GenThingDeviceRespData(deviceSend.GetStatus, uuid, strings.Split(tmpTopic, "/"),
		errors.OK, map[string]any{"logLevel": di.LogLevel})
	er := d.svcCtx.PubDev.PublishToDev(d.ctx, topic, payload)
	if er != nil {
		d.Errorf("%s.PublishToDev failure err:%v", utils.FuncName(), er)
	}
	return er
}
