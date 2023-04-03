package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductSchemaTslReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductSchemaTslReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductSchemaTslReadLogic {
	return &ProductSchemaTslReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取产品信息列表
func (l *ProductSchemaTslReadLogic) ProductSchemaTslRead(in *dm.ProductSchemaTslReadReq) (*dm.ProductSchemaTslReadResp, error) {
	l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(in))
	model, err := l.svcCtx.SchemaRepo.TslRead(l.ctx, in.ProductID)
	if err != nil {
		return nil, err
	}
	return &dm.ProductSchemaTslReadResp{Tsl: model.String()}, nil
}
