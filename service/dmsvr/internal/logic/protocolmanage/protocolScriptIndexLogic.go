package protocolmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolScriptIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolScriptIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolScriptIndexLogic {
	return &ProtocolScriptIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 协议列表
func (l *ProtocolScriptIndexLogic) ProtocolScriptIndex(in *dm.ProtocolScriptIndexReq) (*dm.ProtocolScriptIndexResp, error) {
	var (
		size int64
		err  error
		piDB = relationDB.NewProtocolScriptRepo(l.ctx)
	)

	filter := relationDB.ProtocolScriptFilter{
		Name:         in.Name,
		Status:       in.Status,
		TriggerDir:   in.TriggerDir,
		TriggerTimer: in.TriggerTimer,
	}
	size, err = piDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	di, err := piDB.FindByFilter(l.ctx, filter,
		logic.ToPageInfo(in.Page),
	)
	if err != nil {
		return nil, err
	}

	info := utils.CopySlice[dm.ProtocolScript](di)
	return &dm.ProtocolScriptIndexResp{List: info, Total: size}, nil
}
