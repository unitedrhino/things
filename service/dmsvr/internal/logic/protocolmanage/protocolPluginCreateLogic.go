package protocolmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolPluginCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolPluginCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolPluginCreateLogic {
	return &ProtocolPluginCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 协议创建
func (l *ProtocolPluginCreateLogic) ProtocolPluginCreate(in *dm.ProtocolPlugin) (*dm.WithID, error) {
	po := utils.Copy[relationDB.DmProtocolPlugin](in)
	err := relationDB.NewProtocolPluginRepo(l.ctx).Insert(l.ctx, po)
	return &dm.WithID{Id: po.ID}, err
}
