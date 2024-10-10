package productmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

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
		ProductIDs:   []string{in.ProductID},
		WithCategory: in.WithCategory,
		WithProtocol: in.WithProtocol,
	})
	if err != nil {
		return nil, err
	}
	return logic.ToProductInfo(l.ctx, l.svcCtx, pi), nil
}
