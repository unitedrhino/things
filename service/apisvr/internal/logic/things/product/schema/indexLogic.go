package schema

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.ProductSchemaIndexReq) (resp *types.ProductSchemaIndexResp, err error) {
	dmResp, err := l.svcCtx.ProductM.ProductSchemaIndex(l.ctx, utils.Copy[dm.ProductSchemaIndexReq](req))
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ProductSchemaIndex req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	pis := make([]*types.ProductSchemaInfo, 0, len(dmResp.List))
	for _, v := range dmResp.List {
		pi := ToSchemaInfoTypes(v)
		pis = append(pis, pi)
	}
	return &types.ProductSchemaIndexResp{
		Total: dmResp.Total,
		List:  pis,
	}, nil
}
