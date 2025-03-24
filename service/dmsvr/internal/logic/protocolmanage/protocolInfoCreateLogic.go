package protocolmanagelogic

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/protocol"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/topics"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

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
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}

	po := logic.ToProtocolInfoPo(in)
	if err := protocol.Check(po.ConfigFields, po.ConfigInfos); err != nil {
		return nil, err
	}
	old, err := relationDB.NewProtocolInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.ProtocolInfoFilter{Code: in.Code})
	if err == nil {
		old.Name = in.Name
		old.Desc = in.Desc
		old.IsEnableSyncProduct = in.IsEnableSyncProduct
		old.IsEnableSyncDevice = in.IsEnableSyncDevice
		old.ProductFields = po.ProductFields
		if old.ProductFields == nil {
			old.ProductFields = protocol.ConfigFields{}
		}
		old.ConfigFields = po.ConfigFields
		if old.ConfigFields == nil {
			old.ConfigFields = protocol.ConfigFields{}
		}
		old.Endpoints = po.Endpoints
		if old.Endpoints == nil {
			old.Endpoints = []string{}
		}
		old.EtcdKey = po.EtcdKey
		err := relationDB.NewProtocolInfoRepo(l.ctx).Update(l.ctx, old)
		return &dm.WithID{Id: old.ID}, err
	} else if !errors.NotFind.Eq(err) {
		return nil, err
	}
	err = relationDB.NewProtocolInfoRepo(l.ctx).Insert(l.ctx, po)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.FastEvent.Publish(l.ctx, fmt.Sprintf(topics.DmProtocolInfoUpdate, po.Code), po.ConfigInfos.ToPubStu())
	if err != nil {
		return nil, err
	}
	return &dm.WithID{Id: po.ID}, err
}
