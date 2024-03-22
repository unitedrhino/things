package protocolmanagelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/eventBus"
	"github.com/i-Things/things/service/dmsvr/internal/domain/protocol"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolInfoCreateLogic {
	return &ProtocolInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 协议创建
func (l *ProtocolInfoCreateLogic) ProtocolInfoCreate(in *dm.ProtocolInfo) (*dm.WithID, error) {
	po := logic.ToProtocolInfoPo(in)
	if err := protocol.Check(po.ConfigFields, po.ConfigInfos); err != nil {
		return nil, err
	}
	err := relationDB.NewProtocolInfoRepo(l.ctx).Insert(l.ctx, po)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.FastEvent.Publish(l.ctx, fmt.Sprintf(eventBus.DmProtocolInfoUpdate, po.Code), po.ConfigInfos.ToPubStu())
	if err != nil {
		return nil, err
	}
	return &dm.WithID{Id: po.ID}, err
}
