package product

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSchemaUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaUpdateLogic {
	return &SchemaUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SchemaUpdateLogic) SchemaUpdate(req *types.ProductSchemaUpdateReq) error {
	dmReq := &dm.ProductSchemaUpdateReq{
		Info: &dm.ProductSchema{
			ProductID: req.ProductID, //产品id 只读
			Schema:    req.Schema,
		},
	}

	_, err := l.svcCtx.DmRpc.ProductSchemaUpdate(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.ManageProductTemplate|req=%v|err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
