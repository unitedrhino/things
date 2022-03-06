package logic

import (
	"context"
	"encoding/json"
	"github.com/go-things/things/shared/def"
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/src/dmsvr/dm"
	"github.com/go-things/things/src/dmsvr/internal/svc"
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

func (l *GetDeviceDataLogic) HandleData(in *dm.GetDeviceDataReq) (*dm.GetDeviceDataResp, error) {
	dd := l.svcCtx.DeviceData(l.ctx)
	var dmDatas []*dm.DeviceData
	switch in.Method {
	case def.PROPERTY_METHOD:
		dds, err := dd.GetPropertyDataWithID(in.ProductID, in.DeviceName, in.DataID, def.PageInfo2{
			TimeStart: in.TimeStart,
			TimeEnd:   in.TimeEnd,
			Limit:     in.Limit,
		})
		if err != nil {
			l.Errorf("HandleData|GetPropertyDataWithID|err=%v", err)
			return nil, errors.System
		}
		for _, devData := range dds {
			dmData := dm.DeviceData{
				Timestamp: devData.TimeStamp.UnixMilli(),
				Method:    in.Method,
				DataID:    in.DataID,
			}
			var payload []byte
			payload, _ = json.Marshal(devData.Param)
			dmData.Payload = string(payload)
			dmDatas = append(dmDatas, &dmData)
			l.Slowf("GetDeviceLogLogic|get data=%+v", dmData)
		}
	case def.EVENT_METHOD:
		dds, err := dd.GetEventDataWithID(in.ProductID, in.DeviceName, in.DataID, def.PageInfo2{
			TimeStart: in.TimeStart,
			TimeEnd:   in.TimeEnd,
			Limit:     in.Limit,
		})
		if err != nil {
			l.Errorf("HandleData|GetPropertyDataWithID|err=%v", err)
			return nil, errors.System
		}
		for _, devData := range dds {
			dmData := dm.DeviceData{
				Timestamp: devData.TimeStamp.UnixMilli(),
				Method:    in.Method,
				DataID:    in.DataID,
			}
			var payload []byte
			payload, _ = json.Marshal(devData)
			dmData.Payload = string(payload)
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
		return l.HandleData(in)
	case "status": //获取设备状态信息
	case "logs": //获取设备的调试日志
	default:
		return nil, errors.Method.AddDetail(in.Method)
	}
	return &dm.GetDeviceDataResp{}, nil
}
