package productmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/def"
	"github.com/i-Things/things/src/dmsvr/internal/logic"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"

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
		piDB = relationDB.NewProductInfoRepo(l.ctx)
	)

	filter := relationDB.ProductFilter{
		DeviceType: in.DeviceType, ProductName: in.ProductName,
		Tags: in.Tags, ProductIDs: in.ProductIDs, WithProtocolInfo: true,
	}
	size, err = piDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	di, err := piDB.FindByFilter(l.ctx, filter,
		logic.ToPageInfoWithDefault(in.Page, &def.PageInfo{
			Page: 1, Size: 20,
			Orders: []def.OrderBy{{"created_time", def.OrderDesc}, {"product_id", def.OrderDesc}},
		}),
	)
	if err != nil {
		return nil, err
	}

	info = make([]*dm.ProductInfo, 0, len(di))
	for _, v := range di {
		info = append(info, ToProductInfo(l.ctx, v, l.svcCtx))
	}
	return &dm.ProductInfoIndexResp{List: info, Total: size}, nil
}
