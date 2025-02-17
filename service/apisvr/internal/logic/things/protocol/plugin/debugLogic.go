package plugin

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DebugLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 协议插件调试
func NewDebugLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DebugLogic {
	return &DebugLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DebugLogic) Debug(req *types.ProtocolPluginDebugReq) (resp *types.ProtocolPluginDebugResp, err error) {
	ret, err := l.svcCtx.ProtocolM.ProtocolPluginDebug(l.ctx, utils.Copy[dm.ProtocolPluginDebugReq](req))
	return utils.Copy[types.ProtocolPluginDebugResp](ret), err
}
