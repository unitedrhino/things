package schemamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"gorm.io/gorm"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommonSchemaDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PsDB *relationDB.CommonSchemaRepo
}

func NewCommonSchemaDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommonSchemaDeleteLogic {
	return &CommonSchemaDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PsDB:   relationDB.NewCommonSchemaRepo(ctx),
	}
}

// 删除产品
func (l *CommonSchemaDeleteLogic) CommonSchemaDelete(in *dm.WithID) (*dm.Empty, error) {
	po, err := l.PsDB.FindOneByFilter(l.ctx, relationDB.CommonSchemaFilter{
		ID: in.Id,
	})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddMsg("标识符未找到")
		}
		return nil, err
	}
	if schema.AffordanceType(po.Type) == schema.AffordanceTypeProperty {
		if err := l.svcCtx.SchemaManaRepo.DeleteProperty(l.ctx, nil, "", po.Identifier); err != nil {
			l.Errorf("%s.DeleteProperty failure,err:%v", utils.FuncName(), err)
			return nil, errors.Database.AddDetail(err)
		}
	}
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err = relationDB.NewCommonSchemaRepo(tx).DeleteByFilter(l.ctx, relationDB.CommonSchemaFilter{ID: po.ID})
		if err != nil {
			return err
		}
		err = relationDB.NewProductSchemaRepo(tx).DeleteByFilter(l.ctx, relationDB.ProductSchemaFilter{Identifiers: []string{po.Identifier}})
		if err != nil {
			return err
		}
		return nil
	})

	return &dm.Empty{}, err

}
