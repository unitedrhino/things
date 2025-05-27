package info

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

// 批量导出产品
func NewMultiExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiExportLogic {
	return &MultiExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiExportLogic) MultiExport(req *types.ProductInfoExportReq) (resp *types.ProductInfoExportResp, err error) {
	ret, err := l.svcCtx.ProductM.ProductInfoMultiExport(l.ctx, utils.Copy[dm.ProductInfoExportReq](req))

	return utils.Copy[types.ProductInfoExportResp](ret), err
}
