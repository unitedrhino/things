package plugin

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取协议插件列表
func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.ProtocolPluginIndexReq) (resp *types.ProtocolPluginIndexResp, err error) {
	ret, err := l.svcCtx.ProtocolM.ProtocolPluginIndex(l.ctx, utils.Copy[dm.ProtocolPluginIndexReq](req))
	return utils.Copy[types.ProtocolPluginIndexResp](ret), err
}
