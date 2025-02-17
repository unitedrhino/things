package protocolmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolScriptUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolScriptUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolScriptUpdateLogic {
	return &ProtocolScriptUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 协议更新
func (l *ProtocolScriptUpdateLogic) ProtocolScriptUpdate(in *dm.ProtocolScript) (*dm.Empty, error) {
	old, err := relationDB.NewProtocolScriptRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return &dm.Empty{}, err
	}
	if in.Name != "" {
		old.Name = in.Name
	}
	if in.Desc != nil {
		old.Desc = in.Desc.GetValue()
	}
	if in.Priority != 0 {
		old.Priority = in.Priority
	}
	if in.Script != "" && in.Script != old.Script {
		old.Script = in.Script
	}
	if in.Status != 0 {
		old.Status = in.Status
	}
	err = relationDB.NewProtocolScriptRepo(l.ctx).Update(l.ctx, old)
	return &dm.Empty{}, err
}
