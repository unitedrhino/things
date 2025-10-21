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

type ProtocolScriptDeviceDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolScriptDeviceDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolScriptDeviceDeleteLogic {
	return &ProtocolScriptDeviceDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 协议删除
func (l *ProtocolScriptDeviceDeleteLogic) ProtocolScriptDeviceDelete(in *dm.WithID) (*dm.Empty, error) {
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
	err = relationDB.NewProtocolScriptDeviceRepo(l.ctx).Delete(l.ctx, in.Id)

	return &dm.Empty{}, err
}
