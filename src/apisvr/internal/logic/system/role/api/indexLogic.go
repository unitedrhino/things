package api

import (
	"context"
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

func (l *IndexLogic) Index(req *types.RoleApiIndexReq) (resp *types.RoleApiIndexResp, err error) {
	info, err := l.svcCtx.RoleRpc.RoleApiIndex(l.ctx, &sys.RoleApiIndexReq{
		Id:      req.ID,
		AppCode: req.AppCode,
	})
	if err != nil {
		return nil, err
	}
	return &types.RoleApiIndexResp{ApiIDs: info.ApiIDs}, nil
}
