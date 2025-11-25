package protocolmanagelogic

import (
	"context"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolScriptDeviceCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolScriptDeviceCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolScriptDeviceCreateLogic {
	return &ProtocolScriptDeviceCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 协议创建
func (l *ProtocolScriptDeviceCreateLogic) ProtocolScriptDeviceCreate(in *dm.ProtocolScriptDevice) (*dm.WithID, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	s, err := relationDB.NewProtocolScriptRepo(l.ctx).FindOne(l.ctx, in.ScriptID)
	if err != nil {
		return nil, err
	}
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	if !uc.IsRoot() && in.TenantCode != "" && in.TenantCode != uc.TenantCode {
		return nil, errors.Permissions.AddMsg("普通租户只能绑定自己租户下的")
	}
	if in.TenantCode == "" && ctxs.IsRoot(l.ctx) == nil {
		in.TenantCode = def.TenantCodeCommon
	} else if in.TenantCode != "" {
		if !ctxs.CanHandTenant(l.ctx, s.TenantCode) {
			return nil, errors.Permissions
		}
	}
	po := utils.Copy[relationDB.DmProtocolScriptDevice](in)

	err = relationDB.NewProtocolScriptDeviceRepo(l.ctx).Insert(l.ctx, po)
	return &dm.WithID{Id: po.ID}, err
}
