package logic

import (
	"context"
	"encoding/json"
	"github.com/go-things/things/shared/def"
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/src/dmsvr/dm"
	"github.com/go-things/things/src/dmsvr/internal/svc"
	"github.com/tal-tech/go-zero/core/logx"
)

type GetDeviceLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeviceLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceLogLogic {
	return &GetDeviceLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDeviceLogLogic) HandleData(in *dm.GetDeviceLogReq) (*dm.GetDeviceLogResp, error) {
	dd := l.svcCtx.DeviceData(l.ctx)
	var dmDatas []*dm.DeviceData
	switch in.Method {
	case def.PROPERTY_METHOD:
		dds, err := dd.GetPropertyDataWithID(in.ProductID, in.DeviceName, in.DataID, in.TimeStart, in.TimeEnd, in.Limit)
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
		dds, err := dd.GetEventDataWithID(in.ProductID, in.DeviceName, in.DataID, in.TimeStart, in.TimeEnd, in.Limit)
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
	return &dm.GetDeviceLogResp{
		Total: int64(len(dmDatas)),
		Data:  dmDatas,
	}, nil
}

func (l *GetDeviceLogLogic) GetDeviceLog(in *dm.GetDeviceLogReq) (*dm.GetDeviceLogResp, error) {
	switch in.Method {
	case "property", "action", "event": //获取属性信息,获取操作信息,获取事件信息
		return l.HandleData(in)
	case "status": //获取设备状态信息
	case "logs": //获取设备的调试日志
	default:
		return nil, errors.Method.AddDetail(in.Method)
	}
	return &dm.GetDeviceLogResp{}, nil
}
