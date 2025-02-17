package script

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

func (l *IndexLogic) Index(req *types.ProtocolScriptIndexReq) (resp *types.ProtocolScriptIndexResp, err error) {
	ret, err := l.svcCtx.ProtocolM.ProtocolScriptIndex(l.ctx, utils.Copy[dm.ProtocolScriptIndexReq](req))
	return utils.Copy[types.ProtocolScriptIndexResp](ret), err
}
