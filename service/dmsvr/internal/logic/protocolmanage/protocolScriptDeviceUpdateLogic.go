package protocolmanagelogic

import (
	"context"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolScriptDeviceUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolScriptDeviceUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolScriptDeviceUpdateLogic {
	return &ProtocolScriptDeviceUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 协议更新
func (l *ProtocolScriptDeviceUpdateLogic) ProtocolScriptDeviceUpdate(in *dm.ProtocolScriptDevice) (*dm.Empty, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	old, err := relationDB.NewProtocolScriptDeviceRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.ProtocolScriptDeviceFilter{ID: in.Id})
	if err != nil {
		return &dm.Empty{}, err
	}
	if !ctxs.CanHandTenant(l.ctx, old.TenantCode) {
		return nil, errors.Permissions
	}
	if in.Priority != 0 {
		old.Priority = in.Priority
	}
	if in.Status != 0 {
		old.Status = in.Status
	}
	err = relationDB.NewProtocolScriptDeviceRepo(l.ctx).Update(l.ctx, old)
	return &dm.Empty{}, err
}
