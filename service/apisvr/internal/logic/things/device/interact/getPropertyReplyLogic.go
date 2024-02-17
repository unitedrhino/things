package interact

import (
	"context"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPropertyReplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPropertyReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPropertyReplyLogic {
	return &GetPropertyReplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPropertyReplyLogic) GetPropertyReply(req *types.DeviceInteractGetPropertyReplyReq) (resp *types.DeviceInteractGetPropertyReplyResp, err error) {
	return NewGetPropertyLatestReplyLogic(l.ctx, l.svcCtx).GetPropertyLatestReply(req)
}
