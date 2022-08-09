package logic

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/service/deviceData"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataSchemaLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDataSchemaLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataSchemaLogIndexLogic {
	return &DataSchemaLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备数据信息
func (l *DataSchemaLogIndexLogic) DataSchemaLogIndex(in *dm.DataSchemaLogIndexReq) (*dm.DataSchemaIndexResp, error) {
	var (
		dmDatas []*dm.DataSchemaIndex
		dd      = l.svcCtx.DeviceDataRepo
		total   int64
	)
	if in.Interval != 0 && in.ArgFunc == "" {
		return nil, errors.Parameter.AddMsg("填写了间隔就必须填写聚合函数")
	}
	switch in.Method {
	case def.PROPERTY_METHOD, "":
		dds, err := dd.GetPropertyDataByID(l.ctx, deviceData.FilterOpt{
			Page: def.PageInfo2{
				TimeStart: in.TimeStart,
				TimeEnd:   in.TimeEnd,
				Page:      in.Page.GetPage(),
				Size:      in.Page.GetSize(),
			},
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
			DataID:     in.DataID,
			Interval:   in.Interval,
			ArgFunc:    in.ArgFunc})
		if err != nil {
			l.Errorf("HandleData|GetPropertyDataByID|err=%v", err)
			return nil, err
		}
		for _, devData := range dds {
			dmData := dm.DataSchemaIndex{
				Timestamp: devData.TimeStamp.UnixMilli(),
				DataID:    devData.ID,
			}
			var payload []byte
			if param, ok := devData.Param.(string); ok {
				payload = []byte(param)
			} else {
				payload, _ = json.Marshal(devData.Param)
			}
			dmData.GetValue = string(payload)
			dmDatas = append(dmDatas, &dmData)
		}
		if in.ArgFunc == "" && in.Interval == 0 {
			total, err = dd.GetPropertyCountByID(l.ctx, deviceData.FilterOpt{
				Page: def.PageInfo2{
					TimeStart: in.TimeStart,
					TimeEnd:   in.TimeEnd,
					Page:      in.Page.GetPage(),
					Size:      in.Page.GetSize(),
				},
				ProductID:  in.ProductID,
				DeviceName: in.DeviceName,
				DataID:     in.DataID,
				Interval:   in.Interval,
				ArgFunc:    in.ArgFunc})
			if err != nil {
				l.Errorf("HandleData|GetPropertyCountByID|err=%v", err)
				return nil, err
			}
		}
	case def.EVENT_METHOD:
		dds, err := dd.GetEventDataByID(l.ctx, deviceData.FilterOpt{
			Page: def.PageInfo2{
				TimeStart: in.TimeStart,
				TimeEnd:   in.TimeEnd,
				Page:      in.Page.GetPage(),
				Size:      in.Page.GetSize(),
			},
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
			DataID:     in.DataID,
			Interval:   in.Interval,
			ArgFunc:    in.ArgFunc})
		if err != nil {
			l.Errorf("HandleData|GetPropertyDataByID|err=%v", err)
			return nil, errors.System.AddDetail(err)
		}
		for _, devData := range dds {
			dmData := dm.DataSchemaIndex{
				Timestamp: devData.TimeStamp.UnixMilli(),
				Type:      devData.Type,
				DataID:    devData.ID,
			}
			var payload []byte
			payload, _ = json.Marshal(devData)
			dmData.GetValue = string(payload)
			dmDatas = append(dmDatas, &dmData)
			l.Infof("GetDeviceLogLogic|get data=%+v", dmData)
		}
		if in.ArgFunc == "" && in.Interval == 0 {
			total, err = dd.GetEventCountByID(l.ctx, deviceData.FilterOpt{
				Page: def.PageInfo2{
					TimeStart: in.TimeStart,
					TimeEnd:   in.TimeEnd,
					Page:      in.Page.GetPage(),
					Size:      in.Page.GetSize(),
				},
				ProductID:  in.ProductID,
				DeviceName: in.DeviceName,
				DataID:     in.DataID,
				Interval:   in.Interval,
				ArgFunc:    in.ArgFunc})
			if err != nil {
				l.Errorf("HandleData|GetEventCountByID|err=%v", err)
				return nil, errors.System.AddDetail(err)
			}
		}

	}
	return &dm.DataSchemaIndexResp{
		Total: total,
		List:  dmDatas,
	}, nil
}
