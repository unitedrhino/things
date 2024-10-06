package info

import (
	"context"

	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"gitee.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.ProtocolInfo) (*types.WithID, error) {
	ret, err := l.svcCtx.ProtocolM.ProtocolInfoCreate(l.ctx, ToInfoPb(req))
	if err != nil {
		return nil, err
	}
	return &types.WithID{ID: ret.Id}, nil
}
