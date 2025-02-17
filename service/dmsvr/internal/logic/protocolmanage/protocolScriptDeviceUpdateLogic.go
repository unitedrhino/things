package protocolmanagelogic

import (
	"context"
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
	old, err := relationDB.NewProtocolScriptDeviceRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return &dm.Empty{}, err
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
