package productmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductCategoryMultiExportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductCategoryMultiExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductCategoryMultiExportLogic {
	return &ProductCategoryMultiExportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductCategoryMultiExportLogic) ProductCategoryMultiExport(in *dm.ProductCategoryExportReq) (*dm.ProductCategoryExportResp, error) {
	pos, err := relationDB.NewProductCategoryRepo(l.ctx).FindByFilter(l.ctx, relationDB.ProductCategoryFilter{IDs: in.Ids, WithSchemas: true}, nil)
	if err != nil {
		return nil, err
	}
	return &dm.ProductCategoryExportResp{Categories: utils.MarshalNoErr(pos)}, nil
}
