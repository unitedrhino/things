package protocolmanagelogic

import (
	"context"
	_ "embed"
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

//go:embed script/after
var afterScript string

//go:embed script/before
var beforeScript string

// 协议创建
func (l *ProtocolScriptCreateLogic) ProtocolScriptCreate(in *dm.ProtocolScript) (*dm.WithID, error) {
	po := utils.Copy[relationDB.DmProtocolScript](in)
	switch in.TriggerTimer {
	case protocol.TriggerTimerBefore:
		po.Script = beforeScript
	case protocol.TriggerTimerAfter:
		po.Script = afterScript
	}
	err := relationDB.NewProtocolScriptRepo(l.ctx).Insert(l.ctx, po)
	return &dm.WithID{Id: po.ID}, err
}
