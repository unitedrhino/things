package protocolmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolScriptDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolScriptDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolScriptDeleteLogic {
	return &ProtocolScriptDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 协议删除
func (l *ProtocolScriptDeleteLogic) ProtocolScriptDelete(in *dm.WithID) (*dm.Empty, error) {
	err := relationDB.NewProtocolScriptRepo(l.ctx).Delete(l.ctx, in.Id)

	return &dm.Empty{}, err
}
