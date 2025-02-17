package protocolmanagelogic

import (
	"context"
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
	err := relationDB.NewProtocolScriptDeviceRepo(l.ctx).Delete(l.ctx, in.Id)

	return &dm.Empty{}, err
}
