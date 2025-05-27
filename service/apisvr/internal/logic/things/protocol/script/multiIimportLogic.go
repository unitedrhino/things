package script

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiIimportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量导入协议脚本
func NewMultiIimportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiIimportLogic {
	return &MultiIimportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiIimportLogic) MultiIimport(req *types.ProtocolScriptImportReq) (resp *types.ImportResp, err error) {
	ret, err := l.svcCtx.ProtocolM.ProtocolScriptMultiImport(l.ctx, utils.Copy[dm.ProtocolScriptImportReq](req))

	return utils.Copy[types.ImportResp](ret), err
}
