package script

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量导出协议脚本
func NewMultiExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiExportLogic {
	return &MultiExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiExportLogic) MultiExport(req *types.ProtocolScriptExportReq) (resp *types.ProtocolScriptExportResp, err error) {
	ret, err := l.svcCtx.ProtocolM.ProtocolScriptMultiExport(l.ctx, utils.Copy[dm.ProtocolScriptExportReq](req))

	return utils.Copy[types.ProtocolScriptExportResp](ret), err
}
