package productmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/stores"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

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

// 获取设备信息列表
func (l *ProductInfoIndexLogic) ProductInfoIndex(in *dm.ProductInfoIndexReq) (*dm.ProductInfoIndexResp, error) {
	var (
		info []*dm.ProductInfo
		size int64
		err  error
		piDB = relationDB.NewProductInfoRepo(l.ctx)
	)

	filter := relationDB.ProductFilter{
		SceneMode: in.SceneMode, DeviceType: in.DeviceType, DeviceTypes: in.DeviceTypes, ProductName: in.ProductName, ProtocolCode: in.ProtocolCode,
		Tags: in.Tags, ProductIDs: in.ProductIDs, WithProtocol: in.WithProtocol, WithCategory: in.WithCategory, ProtocolConf: in.ProtocolConf,
		Statuses: in.Statuses, Status: in.Status,
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
