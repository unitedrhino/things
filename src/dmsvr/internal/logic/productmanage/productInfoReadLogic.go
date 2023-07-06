package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductInfoReadLogic {
	return &ProductInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备信息详情
func (l *ProductInfoReadLogic) ProductInfoRead(in *dm.ProductInfoReadReq) (*dm.ProductInfo, error) {
	pi, err := relationDB.NewProductInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.ProductFilter{
		ProductIDs: []string{in.ProductID},
	})
	if err != nil {
		return nil, err
	}
	return ToProductInfo(l.ctx, pi, l.svcCtx), nil
}
