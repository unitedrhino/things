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

type ProtocolPluginIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolPluginIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolPluginIndexLogic {
	return &ProtocolPluginIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 协议列表
func (l *ProtocolPluginIndexLogic) ProtocolPluginIndex(in *dm.ProtocolPluginIndexReq) (*dm.ProtocolPluginIndexResp, error) {
	var (
		size int64
		err  error
		piDB = relationDB.NewProtocolPluginRepo(l.ctx)
	)

	filter := relationDB.ProtocolPluginFilter{
		Name:         in.Name,
		TriggerSrc:   in.TriggerSrc,
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

	info := utils.CopySlice[dm.ProtocolPlugin](di)
	return &dm.ProtocolPluginIndexResp{List: info, Total: size}, nil
}
