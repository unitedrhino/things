package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

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
	)
	filter := mysql.ProductFilter{
		DeviceType: in.DeviceType, ProductName: in.ProductName, ProductIDs: in.ProductIDs}
	size, err = l.svcCtx.DmDB.GetProductsCountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	di, err := l.svcCtx.DmDB.FindProductsByFilter(l.ctx, filter, def.PageInfo{Size: in.Page.Size, Page: in.Page.Page})
	if err != nil {
		return nil, err
	}
	info = make([]*dm.ProductInfo, 0, len(di))
	for _, v := range di {
		info = append(info, ToProductInfo(v))
	}
	return &dm.ProductInfoIndexResp{List: info, Total: size}, nil
}
