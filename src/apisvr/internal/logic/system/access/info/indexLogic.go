package info

import (
	"context"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

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

func (l *IndexLogic) Index(req *types.AccessIndexReq) (resp *types.AccessIndexResp, err error) {
	rst, err := l.svcCtx.AccessRpc.AccessInfoIndex(l.ctx, &sys.AccessInfoIndexReq{
		Page:       logic.ToSysPageRpc(req.Page),
		Group:      req.Group,
		Code:       req.Code,
		Name:       req.Name,
		IsNeedAuth: req.IsNeedAuth,
		WithApis:   req.WithApis,
	})
	if err != nil {
		return nil, err
	}
	return &types.AccessIndexResp{
		List:  ToAccessInfosTypes(rst.List),
		Total: rst.Total,
	}, nil
}
