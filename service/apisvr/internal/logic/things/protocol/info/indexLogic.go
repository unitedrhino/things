package info

import (
	"context"
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

func (l *IndexLogic) Index(req *types.ProtocolInfoIndexReq) (resp *types.ProtocolInfoIndexResp, err error) {
	ret, err := l.svcCtx.ProtocolM.ProtocolInfoIndex(l.ctx, &dm.ProtocolInfoIndexReq{
		Page:          logic.ToDmPageRpc(req.Page),
		Name:          req.Name,
		Code:          req.Code,
		TransProtocol: req.TransProtocol,
	})

	return &types.ProtocolInfoIndexResp{
		List:  ToInfosTypes(ret.List),
		Total: ret.Total,
	}, nil
}
