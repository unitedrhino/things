package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductSchemaCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDb *relationDB.ProductInfoRepo
}

func NewProductSchemaCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductSchemaCreateLogic {
	return &ProductSchemaCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDb:   relationDB.NewProductInfoRepo(ctx),
	}
}

func (l *ProductSchemaCreateLogic) ruleCheck(in *dm.ProductSchemaCreateReq) (*mysql.DmProductSchema, error) {
	_, err := l.PiDb.FindOneByFilter(l.ctx, relationDB.ProductFilter{ProductIDs: []string{in.Info.ProductID}})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddMsg("找不到该产品:" + in.Info.ProductID)
		}
		return nil, err
	}
	_, err = l.svcCtx.ProductSchema.FindOneByProductIDIdentifier(l.ctx, in.Info.ProductID, in.Info.Identifier)
	if err != nil {
		if err == mysql.ErrNotFound {
			po := ToProductSchemaPo(in.Info)
			if po.Name == "" {
				return nil, errors.Parameter.AddMsg("功能名称不能为空")
			}
			if po.Required == 0 {
				po.Required = def.False
			}
			if err := CheckAffordance(po); err != nil {
				return nil, err
			}
			return po, nil
		}
		return nil, errors.Database.AddDetail(err)
	}
	return nil, errors.Parameter.AddMsg("标识符在该产品中已经被使用:" + in.Info.Identifier)
}

// 新增产品
func (l *ProductSchemaCreateLogic) ProductSchemaCreate(in *dm.ProductSchemaCreateReq) (*dm.Response, error) {
	l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(in))
	po, err := l.ruleCheck(in)
	if err != nil {
		l.Errorf("%s.ruleCheck err:%v", err)
		return nil, err
	}

	if schema.AffordanceType(po.Type) == schema.AffordanceTypeProperty {
		if err := l.svcCtx.SchemaManaRepo.CreateProperty(l.ctx, mysql.ToPropertyDo(po), in.Info.ProductID); err != nil {
			l.Errorf("%s.CreateProperty failure,err:%v", utils.FuncName(), err)
			return nil, errors.Database.AddDetail(err)
		}
	}
	_, err = l.svcCtx.ProductSchema.Insert(l.ctx, po)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.DataUpdate.ProductSchemaUpdate(l.ctx, &events.DeviceUpdateInfo{ProductID: in.Info.ProductID})
	if err != nil {
		return nil, err
	}
	return &dm.Response{}, nil
}
