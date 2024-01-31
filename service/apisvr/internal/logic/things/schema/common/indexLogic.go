package common

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/apisvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

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

func (l *IndexLogic) Index(req *types.CommonSchemaIndexReq) (resp *types.CommonSchemaIndexResp, err error) {
	dmReq := &dm.CommonSchemaIndexReq{
		Page:        logic.ToDmPageRpc(req.Page),
		Type:        req.Type,
		Identifiers: req.Identifiers,
		ProductIDs:  req.ProductIDs,
	}
	dmResp, err := l.svcCtx.SchemaM.CommonSchemaIndex(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.CommonSchemaIndex req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	pis := make([]*types.CommonSchemaInfo, 0, len(dmResp.List))
	for _, v := range dmResp.List {
		pi := ToSchemaInfoTypes(v)
		pis = append(pis, pi)
	}
	return &types.CommonSchemaIndexResp{
		Total: dmResp.Total,
		List:  pis,
	}, nil
}
