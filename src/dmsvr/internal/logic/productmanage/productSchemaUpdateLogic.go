package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events/topics"
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

func (l *ProductSchemaUpdateLogic) ruleCheck(in *dm.ProductSchemaUpdateReq) (*mysql.DmProductSchema, *mysql.DmProductSchema, error) {
	_, err := l.svcCtx.ProductInfo.FindOne(l.ctx, in.Info.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, nil, errors.Parameter.AddMsgf("产品id不存在:" + cast.ToString(in.Info.ProductID))
		}
		return nil, nil, errors.Database.AddDetail(err)
	}
	po, err := l.svcCtx.ProductSchema.FindOneByProductIDIdentifier(l.ctx, in.Info.ProductID, in.Info.Identifier)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, nil, errors.Parameter.AddMsgf("标识符不存在:" + in.Info.Identifier)
		}
		return nil, nil, errors.Database.AddDetail(err)
	}
	if po.Tag != in.Info.Tag {
		return nil, nil, errors.Parameter.AddMsg("功能标签不支持修改")
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
	if schema.AffordanceType(po.Type) == schema.AffordanceTypeProperty {
		if err := l.svcCtx.SchemaManaRepo.DeleteProperty(
			l.ctx, in.Info.ProductID, in.Info.Identifier); err != nil {
			l.Errorf("%s.DeleteProperty failure,err:%v", utils.FuncName(), err)
			return nil, errors.Database.AddDetail(err)
		}
	}
	if schema.AffordanceType(newPo.Type) == schema.AffordanceTypeProperty {
		if err := l.svcCtx.SchemaManaRepo.CreateProperty(
			l.ctx, mysql.ToPropertyDo(newPo), in.Info.ProductID); err != nil {
			l.Errorf("%s.CreateProperty failure,err:%v", utils.FuncName(), err)
			return nil, errors.Database.AddDetail(err)
		}
	}
	err = l.svcCtx.ProductSchema.Update(l.ctx, newPo)
	if err != nil {
		return nil, err
	}
	l.svcCtx.Bus.Publish(l.ctx, topics.DmProductSchemaUpdate, in.Info.ProductID)

	return &dm.Response{}, nil
}
