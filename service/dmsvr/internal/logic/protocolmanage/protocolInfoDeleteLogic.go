package protocolmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolInfoDeleteLogic {
	return &ProtocolInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 协议删除
func (l *ProtocolInfoDeleteLogic) ProtocolInfoDelete(in *dm.WithID) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	err := relationDB.NewProtocolInfoRepo(l.ctx).Delete(l.ctx, in.Id)
	return &dm.Empty{}, err
}
