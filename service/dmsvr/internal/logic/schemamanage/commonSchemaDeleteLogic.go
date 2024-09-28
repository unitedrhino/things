package schemamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
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
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	po, err := l.PsDB.FindOneByFilter(l.ctx, relationDB.CommonSchemaFilter{
		ID: in.Id,
	})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddMsg("标识符未找到")
		}
		return nil, err
	}
	count, err := relationDB.NewProductSchemaRepo(l.ctx).CountByFilter(l.ctx,
		relationDB.ProductSchemaFilter{Identifiers: []string{po.Identifier}, Tags: []int64{schema.TagRequired, schema.TagOptional}})
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.Parameter.AddMsgf("有%v个产品绑定该物模型,不允许删除", count)
	}
	if schema.AffordanceType(po.Type) == schema.AffordanceTypeProperty {
		if err := l.svcCtx.SchemaManaRepo.DeleteProperty(l.ctx, nil, "", po.Identifier); err != nil {
			l.Errorf("%s.DeleteProperty failure,err:%v", utils.FuncName(), err)
			return nil, errors.Database.AddDetail(err)
		}
	}
	err = relationDB.NewCommonSchemaRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.CommonSchemaFilter{ID: po.ID})
	if err != nil {
		return nil, err
	}
	return &dm.Empty{}, err

}
