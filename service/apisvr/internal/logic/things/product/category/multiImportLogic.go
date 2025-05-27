package category

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiImportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量导入产品品类
func NewMultiImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiImportLogic {
	return &MultiImportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiImportLogic) MultiImport(req *types.ProductCategoryImportReq) (resp *types.ImportResp, err error) {
	ret, err := l.svcCtx.ProductM.ProductCategoryMultiImport(l.ctx, utils.Copy[dm.ProductCategoryImportReq](req))

	return utils.Copy[types.ImportResp](ret), err
}
