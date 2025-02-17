package protocolmanagelogic

import (
	"context"
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
	po := utils.Copy[relationDB.DmProtocolScriptDevice](in)
	err := relationDB.NewProtocolScriptDeviceRepo(l.ctx).Insert(l.ctx, po)
	return &dm.WithID{Id: po.ID}, err
}
