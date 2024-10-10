package schemamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/domain/schema"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommonSchemaIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PsDB *relationDB.CommonSchemaRepo
}

func NewCommonSchemaIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommonSchemaIndexLogic {
	return &CommonSchemaIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PsDB:   relationDB.NewCommonSchemaRepo(ctx),
	}
}

// 获取产品信息列表
func (l *CommonSchemaIndexLogic) CommonSchemaIndex(in *dm.CommonSchemaIndexReq) (*dm.CommonSchemaIndexResp, error) {
	filter := relationDB.CommonSchemaFilter{
		Type:              in.Type,
		Types:             in.Types,
		Name:              in.Name,
		Identifiers:       in.Identifiers,
		IsCanSceneLinkage: in.IsCanSceneLinkage,
		FuncGroup:         in.FuncGroup,
		UserPerm:          in.UserPerm,
		PropertyMode:      in.PropertyMode,
		ControlMode:       in.ControlMode,
		ProductSceneMode:  in.ProductSceneMode,
	}
	if in.ProductCategoryID != 0 {
		var ProductCategoryIDs = []int64{in.ProductCategoryID}
		if in.ProductCategoryWithFather {
			ProductCategoryIDs = append(ProductCategoryIDs, def.RootNode)
			pc, err := relationDB.NewProductCategoryRepo(l.ctx).FindOne(l.ctx, in.ProductCategoryID)
			if err != nil {
				return nil, err
			}
			ProductCategoryIDs = append(ProductCategoryIDs, utils.GetIDPath(pc.IDPath)...)
		}
		filter.ProductCategoryIDs = ProductCategoryIDs
	}
	if in.AreaID != 0 {
		cols, err := relationDB.NewDeviceInfoRepo(l.ctx).FindProductIDsByFilter(l.ctx, relationDB.DeviceFilter{AreaIDs: []int64{in.AreaID}})
		if err != nil {
			return nil, err
		}
		if len(cols) == 0 {
			return &dm.CommonSchemaIndexResp{}, nil
		}
		in.ProductIDs = append(in.ProductIDs, cols...)
	}

	if in.GroupID != 0 {
		cols, err := relationDB.NewDeviceInfoRepo(l.ctx).FindProductIDsByFilter(l.ctx, relationDB.DeviceFilter{GroupID: in.GroupID})
		if err != nil {
			return nil, err
		}
		if len(cols) == 0 {
			return &dm.CommonSchemaIndexResp{}, nil
		}
		in.ProductIDs = append(in.ProductIDs, cols...)
	}
	if len(in.GroupIDs) != 0 {
		cols, err := relationDB.NewDeviceInfoRepo(l.ctx).FindProductIDsByFilter(l.ctx, relationDB.DeviceFilter{GroupIDs: in.GroupIDs})
		if err != nil {
			return nil, err
		}
		if len(cols) == 0 {
			return &dm.CommonSchemaIndexResp{}, nil
		}
		in.ProductIDs = append(in.ProductIDs, cols...)
	}

	if len(in.ProductIDs) != 0 {
		if len(in.ProductIDs) == 1 { //直接返回该设备的物模型
			f := relationDB.ProductSchemaFilter{
				ProductIDs:        in.ProductIDs,
				Type:              in.Type,
				Types:             in.Types,
				Name:              in.Name,
				Identifiers:       in.Identifiers,
				IsCanSceneLinkage: in.IsCanSceneLinkage,
				FuncGroup:         in.FuncGroup,
				UserPerm:          in.UserPerm,
				ControlMode:       in.ControlMode,
				ProductSceneMode:  in.ProductSceneMode,
				PropertyMode:      in.PropertyMode}
			schemas, err := relationDB.NewProductSchemaRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page).WithDefaultOrder(stores.OrderBy{
				Field: "order",
				Sort:  stores.OrderAsc,
			}))
			if err != nil {
				return nil, err
			}
			total, err := relationDB.NewProductSchemaRepo(l.ctx).CountByFilter(l.ctx, f)
			if err != nil {
				return nil, err
			}
			return &dm.CommonSchemaIndexResp{List: utils.CopySlice[dm.CommonSchemaInfo](schemas), Total: total}, nil
		}
		var productIDs = map[string]struct{}{}
		rst, err := relationDB.NewProductSchemaRepo(l.ctx).FindByFilter(l.ctx, relationDB.ProductSchemaFilter{
			ProductIDs: in.ProductIDs, Tags: []int64{schema.TagOptional, schema.TagRequired}, ProductSceneMode: in.ProductSceneMode, ControlMode: in.ControlMode}, nil)
		if err != nil {
			return nil, err
		}
		var identifyMap = map[string]int{}
		for _, v := range rst {
			identifyMap[v.Identifier]++
			productIDs[v.ProductID] = struct{}{}
		}
		for k, v := range identifyMap {
			if v == len(productIDs) { //每个产品都有的物模型
				filter.Identifiers = append(filter.Identifiers, k)
			}
		}
		if len(filter.Identifiers) == 0 {
			return &dm.CommonSchemaIndexResp{}, nil
		}
	}

	schemas, err := l.PsDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page).WithDefaultOrder(stores.OrderBy{
		Field: "order",
		Sort:  stores.OrderAsc,
	}))
	if err != nil {
		return nil, err
	}
	total, err := l.PsDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	list := make([]*dm.CommonSchemaInfo, 0, len(schemas))
	for _, s := range schemas {
		list = append(list, ToCommonSchemaRpc(s))
	}
	return &dm.CommonSchemaIndexResp{List: list, Total: total}, nil
}
