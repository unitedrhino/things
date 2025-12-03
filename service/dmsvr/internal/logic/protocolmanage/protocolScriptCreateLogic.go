package protocolmanagelogic

import (
	"context"
	_ "embed"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/protocol"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolScriptCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolScriptCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolScriptCreateLogic {
	return &ProtocolScriptCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//go:embed script/before
var beforeScript string

//go:embed script/downAfter
var downAfterScript string

//go:embed script/upAfter
var upAfterScript string

// 协议创建
func (l *ProtocolScriptCreateLogic) ProtocolScriptCreate(in *dm.ProtocolScript) (*dm.WithID, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	if !ctxs.CanHandTenant(l.ctx, in.TenantCode) {
		return nil, errors.Permissions
	}
	po := utils.Copy[relationDB.DmProtocolScript](in)
	po.Script = beforeScript
	if in.TriggerTimer == protocol.TriggerTimerAfter {
		if in.TriggerDir == protocol.TriggerDirUp {
			po.Script = upAfterScript
		} else {
			po.Script = downAfterScript
		}
	}

	err := relationDB.NewProtocolScriptRepo(l.ctx).Insert(l.ctx, po)
	return &dm.WithID{Id: po.ID}, err
}
