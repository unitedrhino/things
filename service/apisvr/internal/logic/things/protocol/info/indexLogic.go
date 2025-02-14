package info

import (
	"context"
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

func (l *IndexLogic) Index(req *types.ProtocolInfoIndexReq) (resp *types.ProtocolInfoIndexResp, err error) {
	ret, err := l.svcCtx.ProtocolM.ProtocolInfoIndex(l.ctx, utils.Copy[dm.ProtocolInfoIndexReq](req))

	return &types.ProtocolInfoIndexResp{
		List:  ToInfosTypes(ret.List),
		Total: ret.Total,
	}, nil
}
