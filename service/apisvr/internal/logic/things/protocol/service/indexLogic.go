package service

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

// 获取自定义协议服务器信息列表
func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.ProtocolServiceIndexReq) (resp *types.ProtocolServiceIndexResp, err error) {
	ret, err := l.svcCtx.ProtocolM.ProtocolServiceIndex(l.ctx, utils.Copy[dm.ProtocolServiceIndexReq](req))

	return &types.ProtocolServiceIndexResp{
		List:  utils.CopySlice[types.ProtocolService](ret.List),
		Total: ret.Total,
	}, nil
}
