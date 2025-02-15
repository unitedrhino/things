package protocolmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolPluginDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolPluginDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolPluginDeleteLogic {
	return &ProtocolPluginDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 协议删除
func (l *ProtocolPluginDeleteLogic) ProtocolPluginDelete(in *dm.WithID) (*dm.Empty, error) {
	err := relationDB.NewProtocolPluginRepo(l.ctx).Delete(l.ctx, in.Id)

	return &dm.Empty{}, err
}
