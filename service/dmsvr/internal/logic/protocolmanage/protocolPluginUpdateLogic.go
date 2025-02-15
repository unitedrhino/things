package protocolmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolPluginUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolPluginUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolPluginUpdateLogic {
	return &ProtocolPluginUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 协议更新
func (l *ProtocolPluginUpdateLogic) ProtocolPluginUpdate(in *dm.ProtocolPlugin) (*dm.Empty, error) {
	old, err := relationDB.NewProtocolPluginRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return &dm.Empty{}, err
	}
	if in.Name != "" {
		old.Name = in.Name
	}
	if in.Desc != nil {
		old.Desc = in.Desc.GetValue()
	}
	if in.Priority != 0 {
		old.Priority = in.Priority
	}
	if in.Script != "" && in.Script != old.Script {
		old.Script = in.Script
	}
	if in.Status != 0 {
		old.Status = in.Status
	}
	err = relationDB.NewProtocolPluginRepo(l.ctx).Update(l.ctx, old)
	return &dm.Empty{}, err
}
