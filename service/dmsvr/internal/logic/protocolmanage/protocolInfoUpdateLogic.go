package protocolmanagelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/eventBus"
	"github.com/i-Things/things/service/dmsvr/internal/domain/protocol"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolInfoUpdateLogic {
	return &ProtocolInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 协议更新
func (l *ProtocolInfoUpdateLogic) ProtocolInfoUpdate(in *dm.ProtocolInfo) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	old, err := relationDB.NewProtocolInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	newPo := logic.ToProtocolInfoPo(in)
	old.TransProtocol = newPo.TransProtocol
	old.Name = newPo.Name
	old.TransProtocol = newPo.TransProtocol
	old.Desc = newPo.Desc
	if newPo.ConfigFields != nil {
		old.ConfigFields = newPo.ConfigFields
	}
	if newPo.ConfigInfos != nil {
		old.ConfigInfos = newPo.ConfigInfos
	}
	old.Endpoints = newPo.Endpoints
	old.EtcdKey = newPo.EtcdKey
	if err := protocol.Check(old.ConfigFields, old.ConfigInfos); err != nil {
		return nil, err
	}
	err = relationDB.NewProtocolInfoRepo(l.ctx).Update(l.ctx, old)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.FastEvent.Publish(l.ctx, fmt.Sprintf(eventBus.DmProtocolInfoUpdate, old.Code), old.ConfigInfos.ToPubStu())
	if err != nil {
		return nil, err
	}
	return &dm.Empty{}, err
}
