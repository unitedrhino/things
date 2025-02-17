package protocolmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolScriptReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolScriptReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolScriptReadLogic {
	return &ProtocolScriptReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 协议详情
func (l *ProtocolScriptReadLogic) ProtocolScriptRead(in *dm.WithID) (*dm.ProtocolScript, error) {
	po, err := relationDB.NewProtocolScriptRepo(l.ctx).FindOne(l.ctx, in.Id)
	return utils.Copy[dm.ProtocolScript](po), err
}
