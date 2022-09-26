package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductSchemaUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductSchemaUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductSchemaUpdateLogic {
	return &ProductSchemaUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductSchemaUpdateLogic) ModifyProductSchema(in *dm.ProductSchemaUpdateReq, oldT *schema.Model) (*dm.Response, error) {
	l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(in))

	//l.Infof("%s ProductID:%v", utils.FuncName(), in.Info.ProductID)
	//newT, err := schema.ValidateWithFmt([]byte(in.Info.Schema))
	//if err != nil {
	//	return nil, err
	//}
	//err = schema.CheckModify(oldT, newT)
	//if err != nil {
	//	return nil, err
	//}
	//if err := l.svcCtx.SchemaManaRepo.UpdateProduct(l.ctx, oldT, newT, in.Info.ProductID); err != nil {
	//	l.Errorf("%s UpdateProduct failure,err:%v", utils.FuncName(), err)
	//	return nil, errors.Database.AddDetail(err)
	//}
	//err = l.svcCtx.SchemaRepo.Update(l.ctx, in.Info.ProductID, newT)
	//if err != nil {
	//	l.Errorf("%s.Update err=%+v", utils.FuncName(), err)
	//	return nil, errors.System.AddDetail(err)
	//}
	//err = l.svcCtx.DataUpdate.ProductSchemaUpdate(l.ctx, &events.DataUpdateInfo{ProductID: in.Info.ProductID})
	//if err != nil {
	//	return nil, err
	//}
	return &dm.Response{}, nil
}

func (l *ProductSchemaUpdateLogic) AddProductSchema(in *dm.ProductSchemaUpdateReq) (*dm.Response, error) {
	l.Infof("%s ProductID:%v", utils.FuncName(), in.Info.ProductID)
	//_, err := l.svcCtx.ProductInfo.FindOne(l.ctx, in.Info.ProductID)
	//if err != nil {
	//	if err == mysql.ErrNotFound {
	//		return nil, errors.Parameter.AddDetail("not find ProductID id:" + cast.ToString(in.Info.ProductID))
	//	}
	//	return nil, errors.Database.AddDetail(err)
	//}
	//t, err := schema.ValidateWithFmt([]byte(in.Info.Schema))
	//if err != nil {
	//	return nil, err
	//}
	//if err := l.svcCtx.HubLogRepo.InitProduct(
	//	l.ctx, in.Info.ProductID); err != nil {
	//	l.Errorf("%s.DeviceLogRepo.InitProduct failure,err:%v", utils.FuncName(), err)
	//	return nil, errors.Database.AddDetail(err)
	//}
	//if err := l.svcCtx.SchemaManaRepo.InitProduct(l.ctx, t, in.Info.ProductID); err != nil {
	//	l.Errorf("%s.SchemaManaRepo.InitProduct failure,err:%v", utils.FuncName(), err)
	//	return nil, errors.Database.AddDetail(err)
	//}
	//err = l.svcCtx.SchemaRepo.Insert(l.ctx, in.Info.ProductID, t)
	//if err != nil {
	//	return nil, err
	//}
	//err = l.svcCtx.DataUpdate.ProductSchemaUpdate(l.ctx, &events.DataUpdateInfo{ProductID: in.Info.ProductID})
	//if err != nil {
	//	return nil, err
	//}
	return &dm.Response{}, nil

	//return &dm.Response{}, err
}

func (l *ProductSchemaUpdateLogic) ruleCheck(in *dm.ProductSchemaUpdateReq) (*mysql.ProductSchema, *mysql.ProductSchema, error) {
	_, err := l.svcCtx.ProductInfo.FindOne(l.ctx, in.Info.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, nil, errors.Parameter.AddDetail("not find ProductID id:" + cast.ToString(in.Info.ProductID))
		}
		return nil, nil, errors.Database.AddDetail(err)
	}
	po, err := l.svcCtx.ProductSchema.FindOneByProductIDIdentifier(l.ctx, in.Info.ProductID, in.Info.Identifier)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, nil, nil
		}
		return nil, nil, errors.Database.AddDetail(err)
	}
	if po.Tag != in.Info.Tag {
		return nil, nil, errors.Parameter.AddMsg("功能标签不支持修改")
	}
	if po.Type != in.Info.Type {
		return nil, nil, errors.Parameter.AddMsg("功能类型不支持修改,请删除后新增")
	}
	newPo := ToProductSchemaPo(in.Info)
	newPo.Id = po.Id
	newPo.Tag = po.Tag
	if in.Info.Affordance == nil {
		newPo.Affordance = po.Affordance
	}
	if in.Info.Name == nil {
		newPo.Name = po.Name
	}
	if in.Info.Desc == nil {
		newPo.Desc = po.Desc
	}
	if in.Info.Required == 0 {
		newPo.Required = po.Required
	}
	if err := CheckAffordance(newPo); err != nil {
		return nil, nil, err
	}
	return po, newPo, nil
}

// 更新产品物模型
func (l *ProductSchemaUpdateLogic) ProductSchemaUpdate(in *dm.ProductSchemaUpdateReq) (*dm.Response, error) {
	po, newPo, err := l.ruleCheck(in)
	if err != nil {
		return nil, err
	}
	if schema.AffordanceType(newPo.Type) == schema.AffordanceTypeProperty {
		if err := l.svcCtx.SchemaManaRepo.UpdateProperty(l.ctx,
			mysql.ToPropertyDo(po), mysql.ToPropertyDo(newPo), in.Info.ProductID); err != nil {
			l.Errorf("%s.SchemaManaRepo.UpdateProperty failure,err:%v", utils.FuncName(), err)
			return nil, errors.Database.AddDetail(err)
		}
	}
	err = l.svcCtx.ProductSchema.Update(l.ctx, newPo)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.DataUpdate.ProductSchemaUpdate(l.ctx, &events.DataUpdateInfo{ProductID: in.Info.ProductID})
	if err != nil {
		return nil, err
	}
	return &dm.Response{}, nil
}
