package common

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

// 批量导入通用物模型
func NewMultiIimportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiIimportLogic {
	return &MultiIimportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiIimportLogic) MultiIimport(req *types.CommonSchemaImportReq) (resp *types.ImportResp, err error) {
	ret, err := l.svcCtx.SchemaM.CommonSchemaMultiImport(l.ctx, utils.Copy[dm.CommonSchemaImportReq](req))
	return utils.Copy[types.ImportResp](ret), err
}
