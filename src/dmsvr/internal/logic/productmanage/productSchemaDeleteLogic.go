package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/utils"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductSchemaDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductSchemaDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductSchemaDeleteLogic {
	return &ProductSchemaDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除产品
func (l *ProductSchemaDeleteLogic) ProductSchemaDelete(in *dm.ProductSchemaDeleteReq) (*dm.Response, error) {
	l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(in))

	return &dm.Response{}, nil
}
