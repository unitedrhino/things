package logic

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceTemplate"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeviceDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeviceDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceDataLogic {
	return &GetDeviceDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDeviceDataLogic) HandleDatas(in *dm.GetDeviceDataReq) (*dm.GetDeviceDataResp, error) {
	var (
		dmDatas []*dm.DeviceData
		total   int
	)

	tempInfo, err := l.svcCtx.ProductTemplate.FindOne(in.ProductID)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	temp, err := deviceTemplate.NewTemplate([]byte(tempInfo.Template))
	if err != nil {
		return nil, errors.System.AddDetail(err)
	}
	dd := l.svcCtx.DeviceDataRepo
	switch in.Method {
	case def.PROPERTY_METHOD:
		total = len(temp.Properties)
		for _, v := range temp.Properties {
			dds, err := dd.GetPropertyDataByID(l.ctx, in.ProductID, in.DeviceName, v.ID, def.PageInfo2{Limit: 1})
			if err != nil {
				l.Errorf("HandleData|GetPropertyDataByID|err=%v", err)
				return nil, errors.System
			}
			if len(dds) == 0 {
				continue
			}
			devData := dds[0]
			dmData := dm.DeviceData{
				Timestamp: devData.TimeStamp.UnixMilli(),
				DataID:    devData.ID,
			}
			var payload []byte
			payload, _ = json.Marshal(devData.Param)
			dmData.GetValue = string(payload)
			dmDatas = append(dmDatas, &dmData)
			l.Slowf("GetDeviceLogLogic|get data=%+v", dmData)
		}
	default:
		return nil, errors.NotRealize.AddDetailf("multi method not implemt:%v", in.Method)
	}
	return &dm.GetDeviceDataResp{
		Total: int64(total),
		List:  dmDatas,
	}, nil
}

func (l *GetDeviceDataLogic) HandleData(in *dm.GetDeviceDataReq) (*dm.GetDeviceDataResp, error) {
	dd := l.svcCtx.DeviceDataRepo
	var dmDatas []*dm.DeviceData
	switch in.Method {
	case def.PROPERTY_METHOD:
		dds, err := dd.GetPropertyDataByID(l.ctx, in.ProductID, in.DeviceName, in.DataID, def.PageInfo2{
			TimeStart: in.TimeStart,
			TimeEnd:   in.TimeEnd,
			Limit:     in.Limit,
		})
		if err != nil {
			l.Errorf("HandleData|GetPropertyDataByID|err=%v", err)
			return nil, errors.System
		}
		for _, devData := range dds {
			dmData := dm.DeviceData{
				Timestamp: devData.TimeStamp.UnixMilli(),
				DataID:    devData.ID,
			}
			var payload []byte
			payload, _ = json.Marshal(devData.Param)
			dmData.GetValue = string(payload)
			dmDatas = append(dmDatas, &dmData)
			l.Slowf("GetDeviceLogLogic|get data=%+v", dmData)
		}
	case def.EVENT_METHOD:
		dds, err := dd.GetEventDataByID(l.ctx, in.ProductID, in.DeviceName, in.DataID, def.PageInfo2{
			TimeStart: in.TimeStart,
			TimeEnd:   in.TimeEnd,
			Limit:     in.Limit,
		})
		if err != nil {
			l.Errorf("HandleData|GetPropertyDataByID|err=%v", err)
			return nil, errors.System
		}
		for _, devData := range dds {
			dmData := dm.DeviceData{
				Timestamp: devData.TimeStamp.UnixMilli(),
				Type:      devData.Type,
				DataID:    devData.ID,
			}
			var payload []byte
			payload, _ = json.Marshal(devData)
			dmData.GetValue = string(payload)
			dmDatas = append(dmDatas, &dmData)
			l.Slowf("GetDeviceLogLogic|get data=%+v", dmData)
		}
	}
	return &dm.GetDeviceDataResp{
		Total: int64(len(dmDatas)),
		List:  dmDatas,
	}, nil
}

func (l *GetDeviceDataLogic) GetDeviceData(in *dm.GetDeviceDataReq) (*dm.GetDeviceDataResp, error) {
	switch in.Method {
	case "property", "action", "event": //获取属性信息,获取操作信息,获取事件信息
		if in.DataID == "" {
			return l.HandleDatas(in)
		}
		return l.HandleData(in)
	case "status": //获取设备状态信息
	case "logs": //获取设备的调试日志
	default:
		return nil, errors.Method.AddDetail(in.Method)
	}
	return &dm.GetDeviceDataResp{}, nil
}
