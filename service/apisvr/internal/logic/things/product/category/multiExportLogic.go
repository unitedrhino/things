package category

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

// 批量导出产品品类
func NewMultiExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiExportLogic {
	return &MultiExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiExportLogic) MultiExport(req *types.ProductCategoryExportReq) (resp *types.ProductCategoryExportResp, err error) {
	ret, err := l.svcCtx.ProductM.ProductCategoryMultiExport(l.ctx, utils.Copy[dm.ProductCategoryExportReq](req))

	return utils.Copy[types.ProductCategoryExportResp](ret), err
}
