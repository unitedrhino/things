package protocolmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolPluginReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolPluginReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolPluginReadLogic {
	return &ProtocolPluginReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 协议详情
func (l *ProtocolPluginReadLogic) ProtocolPluginRead(in *dm.WithID) (*dm.ProtocolPlugin, error) {
	po, err := relationDB.NewProtocolPluginRepo(l.ctx).FindOne(l.ctx, in.Id)
	return utils.Copy[dm.ProtocolPlugin](po), err
}
