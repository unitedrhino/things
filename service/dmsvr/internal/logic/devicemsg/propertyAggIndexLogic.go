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
	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyAggIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPropertyAggIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyAggIndexLogic {
	return &PropertyAggIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PropertyAggIndexLogic) PropertyAggIndex(in *dm.PropertyAggIndexReq) (*dm.PropertyAggIndexResp, error) {
	var (
		//diDatas    []*dm.PropertyLogInfo
		dd         = l.svcCtx.SchemaManaRepo
		productIDs []string
		t          *schema.Model
		err        error
	)
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	if !uc.IsAdmin {
		return nil, errors.Permissions.AddMsg("只允许管理员操作")
	}
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

	dds, err := dd.GetPropertyAgg(l.ctx, t, msgThing.FilterAggOpt{
		TimeStart: in.TimeStart,
		TimeEnd:   in.TimeEnd,
		Filter: msgThing.Filter{
			ProductID:    in.ProductID,
			ProductIDs:   productIDs,
			DeviceNames:  in.DeviceNames,
			BelongGroup:  utils.CopyMap3[def.IDsInfo](in.BelongGroup),
			AreaIDs:      in.AreaIDs,
			AreaID:       in.AreaID,
			AreaIDPath:   in.AreaIDPath,
			Interval:     in.Interval,
			IntervalUnit: def.TimeUnit(in.IntervalUnit),
			PartitionBy:  in.PartitionBy,
		}, Aggs: utils.CopySlice3[msgThing.PropertyAgg](in.Aggs),
	})
	if err != nil {
		l.Errorf("%s.GetPropertyDataByID err=%v", utils.FuncName(), err)
		return nil, err
	}
	var diDatas []*dm.PropertyAggResp
	for _, devData := range dds {
		diData := dm.PropertyAggResp{
			DeviceName:  devData.DeviceName,
			TenantCode:  string(devData.TenantCode),
			ProjectID:   int64(devData.ProjectID),
			AreaID:      int64(devData.AreaID),
			AreaIDPath:  string(devData.AreaIDPath),
			BelongGroup: utils.CopyMap2[dm.IDsInfo](devData.BelongGroup),
		}
		for _, v1 := range devData.Values {
			var dv = dm.PropertyAggRespDetail{DataID: v1.Identifier, TimeWindow: v1.TsWindow.UnixMilli(), Values: map[string]*dm.PropertyAggRespDataDetail{}}
			for k2, v2 := range v1.Values {
				dv2 := dm.PropertyAggRespDataDetail{Timestamp: v2.TimeStamp.UnixMilli()}
				if dv2.Timestamp < 0 {
					dv2.Timestamp = 0
				}
				var payload []byte
				if param, ok := v2.Param.(string); ok {
					payload = []byte(param)
				} else {
					payload = []byte(utils.ToString(v2.Param))
				}
				dv2.Value = string(payload)
				dv.Values[k2] = &dv2
			}
			diData.Values = append(diData.Values, &dv)
		}
		diDatas = append(diDatas, &diData)
	}
	return &dm.PropertyAggIndexResp{List: diDatas}, nil
}
