package schemamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

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
	}
	if in.ProductCategoryID != 0 {
		var ProductCategoryIDs = []int64{in.ProductCategoryID}
		if in.ProductCategoryWithFather {
			pc, err := relationDB.NewProductCategoryRepo(l.ctx).FindOne(l.ctx, in.ProductCategoryID)
			if err != nil {
				return nil, err
			}
			ProductCategoryIDs = append(ProductCategoryIDs, utils.GetIDPath(pc.IDPath)...)
		}
		pcs, err := relationDB.NewProductCategorySchemaRepo(l.ctx).FindByFilter(l.ctx, relationDB.ProductCategorySchemaFilter{ProductCategoryIDs: ProductCategoryIDs}, nil)
		if err != nil {
			return nil, err
		}
		if len(pcs) == 0 {
			return &dm.CommonSchemaIndexResp{}, nil
		}
		ids := utils.ToSliceWithFunc(pcs, func(in *relationDB.DmProductCategorySchema) string {
			return in.Identifier
		})
		filter.Identifiers = append(filter.Identifiers, ids...)
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
		rst, err := relationDB.NewProductSchemaRepo(l.ctx).FindByFilter(l.ctx, relationDB.ProductSchemaFilter{ProductIDs: in.ProductIDs, Tags: []int64{schema.TagOptional, schema.TagRequired}}, nil)
		if err != nil {
			return nil, err
		}
		var identifyMap = map[string]int{}
		for _, v := range rst {
			identifyMap[v.Identifier]++
		}
		for k, v := range identifyMap {
			if v == len(in.ProductIDs) { //每个产品都有的物模型
				filter.Identifiers = append(filter.Identifiers, k)
			}
		}
		if len(filter.Identifiers) == 0 {
			return &dm.CommonSchemaIndexResp{}, nil
		}
	}

	schemas, err := l.PsDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page))
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
