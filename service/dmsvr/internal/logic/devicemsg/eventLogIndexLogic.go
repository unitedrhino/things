package devicemsglogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type EventLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEventLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EventLogIndexLogic {
	return &EventLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备数据信息
func (l *EventLogIndexLogic) EventLogIndex(in *dm.EventLogIndexReq) (*dm.EventLogIndexResp, error) {
	var (
		diDatas []*dm.EventLogInfo
		dd      = l.svcCtx.SchemaManaRepo
		total   int64
	)
	if len(in.DeviceNames) == 0 {
		return &dm.EventLogIndexResp{}, nil
	}
	for _, device := range in.DeviceNames {
		_, err := logic.SchemaAccess(l.ctx, l.svcCtx, def.AuthRead, devices.Core{
			ProductID:  in.ProductID,
			DeviceName: device,
		}, nil)
		if err != nil {
			return nil, err
		}
	}

	dds, err := dd.GetEventDataByFilter(l.ctx, msgThing.FilterOpt{
		Page: def.PageInfo2{
			TimeStart: in.TimeStart,
			TimeEnd:   in.TimeEnd,
			Page:      in.Page.GetPage(),
			Size:      in.Page.GetSize(),
		},
		ProductID:   in.ProductID,
		DeviceNames: in.DeviceNames,
		DataID:      in.DataID})
	if err != nil {
		l.Errorf("%s.GetEventDataByFilter err=%v", utils.FuncName(), err)
		return nil, errors.System.AddDetail(err)
	}
	for _, devData := range dds {
		diData := dm.EventLogInfo{
			Timestamp: devData.TimeStamp.UnixMilli(),
			Type:      devData.Type,
			DataID:    devData.Identifier,
		}
		var payload []byte
		payload, _ = json.Marshal(devData.Params)
		diData.Params = string(payload)
		diDatas = append(diDatas, &diData)
		l.Infof("%s get data=%+v", utils.FuncName(), diData)
	}

	total, err = dd.GetEventCountByFilter(l.ctx, msgThing.FilterOpt{
		Page: def.PageInfo2{
			TimeStart: in.TimeStart,
			TimeEnd:   in.TimeEnd,
			Page:      in.Page.GetPage(),
			Size:      in.Page.GetSize(),
		},
		ProductID:   in.ProductID,
		DeviceNames: in.DeviceNames,
		DataID:      in.DataID})
	if err != nil {
		l.Errorf("%s.GetEventCountByFilter err=%v", utils.FuncName(), err)
		return nil, errors.System.AddDetail(err)
	}

	return &dm.EventLogIndexResp{
		Total: total,
		List:  diDatas,
	}, nil
}
