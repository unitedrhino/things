package protocolmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolScriptDeviceReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolScriptDeviceReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolScriptDeviceReadLogic {
	return &ProtocolScriptDeviceReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 协议详情
func (l *ProtocolScriptDeviceReadLogic) ProtocolScriptDeviceRead(in *dm.WithID) (*dm.ProtocolScriptDevice, error) {
	po, err := relationDB.NewProtocolScriptDeviceRepo(l.ctx).FindOne(l.ctx, in.Id)
	return utils.Copy[dm.ProtocolScriptDevice](po), err
}
