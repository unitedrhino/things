package devicemsglogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"

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
		if in.DeviceName == "" {
			return nil, errors.Parameter.AddMsg("需要填写设备")
		}
		in.DeviceNames = append(in.DeviceNames, in.DeviceName)
	}
	if len(in.DeviceNames) == 0 {
		return nil, errors.Parameter.AddMsg("需要填写设备")
	}
	for _, dev := range in.DeviceNames {
		_, err := logic.SchemaAccess(l.ctx, l.svcCtx, def.AuthRead, devices.Core{
			ProductID:  in.ProductID,
			DeviceName: dev,
		}, nil)
		if err != nil {
			return nil, err
		}
	}
	page := def.PageInfo2{
		TimeStart: in.TimeStart,
		TimeEnd:   in.TimeEnd,
		Page:      in.Page.GetPage(),
		Size:      in.Page.GetSize(),
	}
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	if !uc.IsAdmin {
		di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
		})
		if err != nil {
			return nil, err
		}
		if di.LastBind*1000 > page.TimeStart {
			page.TimeStart = di.LastBind * 1000
		}
	}
	dds, err := dd.GetEventDataByFilter(l.ctx, msgThing.FilterOpt{
		Page:        page,
		ProductID:   in.ProductID,
		DeviceNames: in.DeviceNames,
		DataID:      in.DataID})
	if err != nil {
		l.Errorf("%s.GetEventDataByFilter err=%v", utils.FuncName(), err)
		return nil, errors.System.AddDetail(err)
	}
	for _, devData := range dds {
		diData := dm.EventLogInfo{
			Timestamp:  devData.TimeStamp.UnixMilli(),
			Type:       devData.Type,
			DataID:     devData.Identifier,
			DeviceName: devData.DeviceName,
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
