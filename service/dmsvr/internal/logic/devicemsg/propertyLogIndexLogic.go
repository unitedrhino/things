package devicemsglogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/share/domain/schema"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPropertyLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyLogIndexLogic {
	return &PropertyLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备数据信息
func (l *PropertyLogIndexLogic) PropertyLogIndex(in *dm.PropertyLogIndexReq) (*dm.PropertyLogIndexResp, error) {
	var (
		diDatas []*dm.PropertyLogInfo
		dd      = l.svcCtx.SchemaManaRepo
		total   int64
		t       *schema.Model
		err     error
	)
	if in.Interval != 0 && in.ArgFunc == "" {
		return nil, errors.Parameter.AddMsg("填写了间隔就必须填写聚合函数")
	}
	if len(in.DeviceNames) == 0 {
		if in.DeviceName == "" {
			return nil, errors.Parameter.AddMsg("需要填写设备")
		}
		in.DeviceNames = append(in.DeviceNames, in.DeviceName)
	}
	if len(in.DeviceNames) == 0 {
		return nil, errors.Parameter.AddMsg("需要填写设备")
	}
	if len(in.DeviceNames) == 1 {
		t, err = l.svcCtx.DeviceSchemaRepo.GetData(l.ctx, devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceNames[0]})
		if err != nil {
			return nil, err
		}
	} else {
		t, err = l.svcCtx.ProductSchemaRepo.GetData(l.ctx, in.ProductID)
		if err != nil {
			return nil, err
		}
	}

	p, ok := t.Property[in.DataID]
	if !ok {
		id, _, ok := schema.GetArray(in.DataID)
		if ok {
			p, ok = t.Property[id]
			if !ok {
				return nil, errors.Parameter.AddMsg("标识符未找到")
			}
		} else {
			return nil, errors.Parameter.AddMsg("标识符未找到")
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
		var lastBind int64
		for _, d := range in.DeviceNames {
			di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{ProductID: in.ProductID, DeviceName: d})
			if err != nil {
				return nil, err
			}
			if di.LastBind > lastBind {
				lastBind = di.LastBind
			}
		}
		if lastBind*1000 > page.TimeStart {
			page.TimeStart = lastBind * 1000
		}
	}
	dds, err := dd.GetPropertyDataByID(l.ctx, p, msgThing.FilterOpt{
		Page:        page,
		ProductID:   in.ProductID,
		DeviceNames: in.DeviceNames,
		//DeviceNames: in.DeviceNames,
		Order:        in.Order,
		DataID:       in.DataID,
		Fill:         in.Fill,
		Interval:     in.Interval,
		IntervalUnit: def.TimeUnit(in.IntervalUnit),
		ArgFunc:      in.ArgFunc})
	if err != nil {
		l.Errorf("%s.GetPropertyDataByID err=%v", utils.FuncName(), err)
		return nil, err
	}
	for _, devData := range dds {
		if devData.TimeStamp.IsZero() {
			continue
		}
		if in.Interval != 0 { //如果走了聚合函数,则需要将时间戳取整
			devData.TimeStamp = devData.TimeStamp.Truncate(def.TimeUnit(in.IntervalUnit).ToDuration(in.Interval))
		}
		diData := dm.PropertyLogInfo{
			Timestamp: devData.TimeStamp.UnixMilli(),
			DataID:    devData.Identifier,
		}
		var payload []byte
		if param, ok := devData.Param.(string); ok {
			payload = []byte(param)
		} else {
			payload, _ = json.Marshal(devData.Param)
		}
		diData.Value = string(payload)
		v, err := p.Define.FmtValue(string(payload))
		if err == nil {
			diData.Value = cast.ToString(v)
		}
		diDatas = append(diDatas, &diData)
	}
	if in.ArgFunc == "" && in.Interval == 0 {
		total, err = dd.GetPropertyCountByID(l.ctx, p, msgThing.FilterOpt{
			Page: def.PageInfo2{
				TimeStart: in.TimeStart,
				TimeEnd:   in.TimeEnd,
				Page:      in.Page.GetPage(),
				Size:      in.Page.GetSize(),
			},
			ProductID:    in.ProductID,
			DataID:       in.DataID,
			DeviceNames:  in.DeviceNames,
			Interval:     in.Interval,
			IntervalUnit: def.TimeUnit(in.IntervalUnit),
			ArgFunc:      in.ArgFunc})
		if err != nil {
			l.Errorf("%s.GetPropertyCountByID err=%v", utils.FuncName(), err)
			return nil, err
		}
	}

	return &dm.PropertyLogIndexResp{
		Total: total,
		List:  diDatas,
	}, nil
}
