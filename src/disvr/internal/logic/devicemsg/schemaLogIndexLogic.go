package devicemsglogic

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"

	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/disvr/pb/di"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSchemaLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaLogIndexLogic {
	return &SchemaLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备数据信息
func (l *SchemaLogIndexLogic) SchemaLogIndex(in *di.SchemaLogIndexReq) (*di.SchemaIndexResp, error) {
	var (
		diDatas []*di.SchemaIndex
		dd      = l.svcCtx.SchemaMsgRepo
		total   int64
	)
	if in.Interval != 0 && in.ArgFunc == "" {
		return nil, errors.Parameter.AddMsg("填写了间隔就必须填写聚合函数")
	}
	switch in.Method {
	case devices.PropertyMethod, "":
		dds, err := dd.GetPropertyDataByID(l.ctx, deviceMsg.FilterOpt{
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
			diData := di.SchemaIndex{
				Timestamp: devData.TimeStamp.UnixMilli(),
				DataID:    devData.ID,
			}
			var payload []byte
			if param, ok := devData.Param.(string); ok {
				payload = []byte(param)
			} else {
				payload, _ = json.Marshal(devData.Param)
			}
			diData.GetValue = string(payload)
			diDatas = append(diDatas, &diData)
		}
		if in.ArgFunc == "" && in.Interval == 0 {
			total, err = dd.GetPropertyCountByID(l.ctx, deviceMsg.FilterOpt{
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
	case devices.EventMethod:
		dds, err := dd.GetEventDataByID(l.ctx, deviceMsg.FilterOpt{
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
			diData := di.SchemaIndex{
				Timestamp: devData.TimeStamp.UnixMilli(),
				Type:      devData.Type,
				DataID:    devData.ID,
			}
			var payload []byte
			payload, _ = json.Marshal(devData)
			diData.GetValue = string(payload)
			diDatas = append(diDatas, &diData)
			l.Infof("GetDeviceLogLogic|get data=%+v", diData)
		}
		if in.ArgFunc == "" && in.Interval == 0 {
			total, err = dd.GetEventCountByID(l.ctx, deviceMsg.FilterOpt{
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
	return &di.SchemaIndexResp{
		Total: total,
		List:  diDatas,
	}, nil
}
