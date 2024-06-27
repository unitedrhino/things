package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductCategoryIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductCategoryIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductCategoryIndexLogic {
	return &ProductCategoryIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取产品信息列表
func (l *ProductCategoryIndexLogic) ProductCategoryIndex(in *dm.ProductCategoryIndexReq) (*dm.ProductCategoryIndexResp, error) {
	var (
		info []*dm.ProductCategory
		size int64
		err  error
		piDB = relationDB.NewProductCategoryRepo(l.ctx)
	)
	f := relationDB.ProductCategoryFilter{Name: in.Name, ParentID: in.ParentID, IDs: in.Ids}
	size, err = piDB.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}

	di, err := piDB.FindByFilter(l.ctx, f,
		logic.ToPageInfo(in.Page),
	)
	if err != nil {
		return nil, err
	}

	info = make([]*dm.ProductCategory, 0, len(di))
	for _, v := range di {
		info = append(info, logic.ToProductCategoryPb(l.ctx, l.svcCtx, v, nil))
	}
	return &dm.ProductCategoryIndexResp{List: info, Total: size}, nil
}
