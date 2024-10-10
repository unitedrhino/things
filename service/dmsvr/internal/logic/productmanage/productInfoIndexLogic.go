package productmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"time"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductInfoIndexLogic {
	return &ProductInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var (
	getProduct *caches.Cache[[]string, int64]
)

func init() {
	cache, err := caches.NewCache(caches.CacheConfig[[]string, int64]{
		GetData: func(ctx context.Context, key int64) (*[]string, error) {
			pis, err := relationDB.NewDeviceInfoRepo(ctx).FindProductIDsByFilter(ctx, relationDB.DeviceFilter{ProjectIDs: []int64{key}})
			if err != nil {
				return nil, err
			}
			return &pis, nil
		},
		ExpireTime: 10 * time.Minute,
	})
	logx.Must(err)
	getProduct = cache
}

// 获取设备信息列表
func (l *ProductInfoIndexLogic) ProductInfoIndex(in *dm.ProductInfoIndexReq) (*dm.ProductInfoIndexResp, error) {
	var (
		info []*dm.ProductInfo
		size int64
		err  error
		piDB = relationDB.NewProductInfoRepo(l.ctx)
	)

	filter := relationDB.ProductFilter{
		CategoryIDs: in.CategoryIDs,
		SceneMode:   in.SceneMode, SceneModes: in.SceneModes, DeviceType: in.DeviceType, DeviceTypes: in.DeviceTypes, ProductName: in.ProductName, ProtocolCode: in.ProtocolCode,
		Tags: in.Tags, ProductIDs: in.ProductIDs, WithProtocol: in.WithProtocol, WithCategory: in.WithCategory, ProtocolConf: in.ProtocolConf,
		Statuses: in.Statuses, Status: in.Status, NetType: in.NetType, AreaID: in.AreaID,
	}
	if !ctxs.IsTenantDefault(l.ctx) {
		filter.Status = devices.ProductStatusEnable
	}
	if in.ProjectID != 0 {
		uc := ctxs.GetUserCtxNoNil(l.ctx)
		if !uc.IsAdmin && uc.ProjectAuth[in.ProjectID] == nil { //如果没有项目的权限
			return nil, errors.Permissions.AddMsg("没有项目权限")
		}
		pis, err := getProduct.GetData(l.ctx, in.ProjectID)
		if err != nil {
			return nil, err
		}
		filter.ProductIDs = append(filter.ProductIDs, *pis...)
	}
	size, err = piDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	di, err := piDB.FindByFilter(l.ctx, filter,
		logic.ToPageInfo(in.Page).WithDefaultOrder(stores.OrderBy{Field: "created_time", Sort: stores.OrderDesc},
			stores.OrderBy{Field: "product_id", Sort: stores.OrderDesc}))
	if err != nil {
		return nil, err
	}

	info = make([]*dm.ProductInfo, 0, len(di))
	for _, v := range di {
		info = append(info, logic.ToProductInfo(l.ctx, l.svcCtx, v))
	}
	return &dm.ProductInfoIndexResp{List: info, Total: size}, nil
}
