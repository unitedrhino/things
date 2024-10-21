package protocolmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolServiceDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolServiceDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolServiceDeleteLogic {
	return &ProtocolServiceDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProtocolServiceDeleteLogic) ProtocolServiceDelete(in *dm.WithID) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	err := relationDB.NewProtocolServiceRepo(l.ctx).Delete(l.ctx, in.Id)

	return &dm.Empty{}, err
}
