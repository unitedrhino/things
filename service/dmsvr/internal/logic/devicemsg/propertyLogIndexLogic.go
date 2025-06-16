package devicemsglogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
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
		diDatas    []*dm.PropertyLogInfo
		dd         = l.svcCtx.SchemaManaRepo
		productIDs []string
		total      int64
		t          *schema.Model
		err        error
	)
	if in.Interval != 0 && in.ArgFunc == "" {
		return nil, errors.Parameter.AddMsg("填写了间隔就必须填写聚合函数")
	}
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	if len(in.DeviceNames) == 0 {
		if in.DeviceName == "" && !uc.IsAdmin {
			return nil, errors.Parameter.AddMsg("需要填写设备")
		}
		if in.DeviceName != "" {
			in.DeviceNames = append(in.DeviceNames, in.DeviceName)
		}
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
	if len(in.DeviceNames) == 1 && in.ProductCategoryID == 0 {
		t, err = l.svcCtx.DeviceSchemaRepo.GetData(l.ctx, devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceNames[0]})
		if err != nil {
			return nil, err
		}
	} else {
		if in.ProductID == "" && in.ProductCategoryID == 0 {
			return nil, errors.Parameter.AddMsg("请填写产品ID或品类ID")
		}
		if in.ProductID == "" {
			pis, err := relationDB.NewProductInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.ProductFilter{CategoryIDs: []int64{in.ProductCategoryID}}, nil)
			if err != nil {
				return nil, err
			}
			if len(pis) == 0 {
				return nil, errors.NotFind.AddMsg("未找到产品")
			}
			for _, p := range pis {
				productIDs = append(productIDs, p.ProductID)
			}
			t, err = l.svcCtx.ProductSchemaRepo.GetData(l.ctx, pis[0].ProductID)
			if err != nil {
				return nil, err
			}
		} else {
			t, err = l.svcCtx.ProductSchemaRepo.GetData(l.ctx, in.ProductID)
			if err != nil {
				return nil, err
			}
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
		Orders:    []def.OrderBy{},
	}
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
		Page:         page,
		ProductID:    in.ProductID,
		ProductIDs:   productIDs,
		DeviceNames:  in.DeviceNames,
		BelongGroup:  utils.CopyMap3[def.IDsInfo](in.BelongGroup),
		Order:        in.Order,
		DataID:       in.DataID,
		AreaIDs:      in.AreaIDs,
		AreaID:       in.AreaID,
		AreaIDPath:   in.AreaIDPath,
		Fill:         in.Fill,
		Interval:     in.Interval,
		IntervalUnit: def.TimeUnit(in.IntervalUnit),
		PartitionBy:  in.PartitionBy,
		NoFirstTs:    in.NoFirstTs,
		ArgFunc:      in.ArgFunc})
	if err != nil {
		l.Errorf("%s.GetPropertyDataByID err=%v", utils.FuncName(), err)
		return nil, err
	}
	for _, devData := range dds {
		if devData.TimeStamp.IsZero() && devData.Param == nil {
			continue
		}
		diData := dm.PropertyLogInfo{
			DeviceName:  devData.DeviceName,
			Timestamp:   devData.TimeStamp.UnixMilli(),
			DataID:      devData.Identifier,
			TenantCode:  string(devData.TenantCode),
			ProjectID:   int64(devData.ProjectID),
			AreaID:      int64(devData.AreaID),
			AreaIDPath:  string(devData.AreaIDPath),
			BelongGroup: utils.CopyMap2[dm.IDsInfo](devData.BelongGroup),
			//GroupIDs:     devData.GroupIDs,
			//GroupIDPaths: devData.GroupIDPaths,
		}
		var payload []byte
		if param, ok := devData.Param.(string); ok {
			payload = []byte(param)
		} else {
			payload = []byte(utils.ToString(devData.Param))
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
			AreaIDs:      in.AreaIDs,
			AreaID:       in.AreaID,
			AreaIDPath:   in.AreaIDPath,
			ProductID:    in.ProductID,
			DataID:       in.DataID,
			DeviceNames:  in.DeviceNames,
			Interval:     in.Interval,
			BelongGroup:  utils.CopyMap3[def.IDsInfo](in.BelongGroup),
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
