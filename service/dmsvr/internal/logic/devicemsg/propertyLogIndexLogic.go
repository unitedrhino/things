package devicemsglogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/share/domain/schema"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
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
	for _, device := range in.DeviceNames {
		_, err := logic.SchemaAccess(l.ctx, l.svcCtx, def.AuthRead, devices.Core{
			ProductID:  in.ProductID,
			DeviceName: device,
		}, nil)
		if err != nil {
			return nil, err
		}
	}
	var (
		diDatas []*dm.PropertyLogInfo
		dd      = l.svcCtx.SchemaManaRepo
		total   int64
	)
	if in.Interval != 0 && in.ArgFunc == "" {
		return nil, errors.Parameter.AddMsg("填写了间隔就必须填写聚合函数")
	}
	t, err := l.svcCtx.SchemaRepo.GetData(l.ctx, in.ProductID)
	if err != nil {
		return nil, err
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
	dds, err := dd.GetPropertyDataByID(l.ctx, p, msgThing.FilterOpt{
		Page: def.PageInfo2{
			TimeStart: in.TimeStart,
			TimeEnd:   in.TimeEnd,
			Page:      in.Page.GetPage(),
			Size:      in.Page.GetSize(),
		},
		ProductID:   in.ProductID,
		DeviceNames: in.DeviceNames,
		Order:       in.Order,
		DataID:      in.DataID,
		Fill:        in.Fill,
		Interval:    in.Interval,
		ArgFunc:     in.ArgFunc})
	if err != nil {
		l.Errorf("%s.GetPropertyDataByID err=%v", utils.FuncName(), err)
		return nil, err
	}
	for _, devData := range dds {
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
			ProductID:   in.ProductID,
			DeviceNames: in.DeviceNames,
			DataID:      in.DataID,
			Interval:    in.Interval,
			ArgFunc:     in.ArgFunc})
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
