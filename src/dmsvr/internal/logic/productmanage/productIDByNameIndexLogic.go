package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductIDByNameIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductIDByNameIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductIDByNameIndexLogic {
	return &ProductIDByNameIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通过产品名称获取产品ID列表
func (l *ProductIDByNameIndexLogic) ProductIDByNameIndex(in *dm.ProductIDByNameIndexReq) (*dm.ProductIDByNameIndexResp, error) {
	ids, err := l.svcCtx.ProductInfo.FindIDsByNames(l.ctx, in.ProductNames)
	if err != nil {
		return nil, err
	}
	return &dm.ProductIDByNameIndexResp{ProductIDs: ids}, nil
}
