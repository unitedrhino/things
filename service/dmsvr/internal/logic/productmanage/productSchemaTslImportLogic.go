package productmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/spf13/cast"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductSchemaTslImportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
}

func NewProductSchemaTslImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductSchemaTslImportLogic {
	return &ProductSchemaTslImportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
	}
}

// 删除产品
func (l *ProductSchemaTslImportLogic) ProductSchemaTslImport(in *dm.ProductSchemaTslImportReq) (*dm.Empty, error) {
	l.Infof("%s req:%v", utils.FuncName(), in)
	_, err := l.PiDB.FindOneByFilter(l.ctx, relationDB.ProductFilter{ProductIDs: []string{in.ProductID}})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddDetail("not find ProductID id:" + cast.ToString(in.ProductID))
		}
		return nil, err
	}
	t, err := schema.ValidateWithFmt([]byte(in.Tsl))
	if err != nil {
		return nil, err
	}
	{ //更新td物模型表

		db := relationDB.NewProductSchemaRepo(l.ctx)
		dbSchemas, err := db.FindByFilter(l.ctx, relationDB.ProductSchemaFilter{ProductID: in.ProductID}, nil)
		if err != nil {
			return nil, err
		}
		oldT := relationDB.ToSchemaDo(in.ProductID, dbSchemas)
		if err := l.svcCtx.SchemaManaRepo.DeleteProduct(l.ctx, oldT, in.ProductID); err != nil {
			l.Errorf("%s.SchemaManaRepo.InitProduct failure,err:%v", utils.FuncName(), err)
			return nil, err
		}
		if err := l.svcCtx.SchemaManaRepo.InitProduct(l.ctx, t, in.ProductID); err != nil {
			l.Errorf("%s.SchemaManaRepo.InitProduct failure,err:%v", utils.FuncName(), err)
			return nil, err
		}
		if err := db.MultiUpdate(l.ctx, in.ProductID, t); err != nil {
			l.Errorf("%s.db.MultiUpdate failure,err:%v", utils.FuncName(), err)
			return nil, err
		}
	}
	err = l.svcCtx.SchemaRepo.SetData(l.ctx, in.ProductID, t)
	if err != nil {
		return nil, err
	}
	return &dm.Empty{}, nil
}
