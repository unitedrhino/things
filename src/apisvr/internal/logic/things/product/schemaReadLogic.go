package product

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSchemaReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaReadLogic {
	return &SchemaReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SchemaReadLogic) SchemaRead(req *types.ProductSchemaReadReq) (resp *types.ProductSchemaReadResp, err error) {
	dmReq := &dm.ProductSchemaReadReq{
		ProductID: req.ProductID, //产品id
	}
	dmResp, err := l.svcCtx.DmRpc.ProductSchemaRead(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	respProductSchema := productSchemaToApi(dmResp)
	return &types.ProductSchemaReadResp{ProductSchema: respProductSchema}, nil
}
