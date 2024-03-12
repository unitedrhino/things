package productmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	commonschemalogic "github.com/i-Things/things/service/dmsvr/internal/logic/schemamanage"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductSchemaUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
	PsDB *relationDB.ProductSchemaRepo
}

func NewProductSchemaUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductSchemaUpdateLogic {
	return &ProductSchemaUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		PsDB:   relationDB.NewProductSchemaRepo(ctx),
	}
}

func (l *ProductSchemaUpdateLogic) ruleCheck(in *dm.ProductSchemaUpdateReq) (*relationDB.DmProductSchema, *relationDB.DmProductSchema, error) {
	_, err := l.PiDB.FindOneByFilter(l.ctx, relationDB.ProductFilter{ProductIDs: []string{in.Info.ProductID}})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, nil, errors.Parameter.AddMsgf("产品id不存在:" + cast.ToString(in.Info.ProductID))
		}
		return nil, nil, errors.Database.AddDetail(err)
	}
	po, err := l.PsDB.FindOneByFilter(l.ctx, relationDB.ProductSchemaFilter{
		ProductID: in.Info.ProductID, Identifiers: []string{in.Info.Identifier},
	})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, nil, errors.Parameter.AddMsgf("标识符不存在:" + in.Info.Identifier)
		}
		return nil, nil, err
	}
	if po.Tag != in.Info.Tag {
		return nil, nil, errors.Parameter.AddMsg("功能标签不支持修改")
	}
	newPo := logic.ToProductSchemaPo(in.Info)
	newPo.ID = po.ID
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
	if in.Info.IsCanSceneLinkage == 0 {
		newPo.IsCanSceneLinkage = po.IsCanSceneLinkage
	}
	if in.Info.IsShareAuthPerm == 0 {
		newPo.IsShareAuthPerm = po.IsShareAuthPerm
	}

	if in.Info.IsHistory == 0 {
		newPo.IsHistory = po.IsHistory
	}

	if in.Info.Order == 0 {
		newPo.Order = po.Order
	}

	if in.Info.ExtendConfig == "" {
		newPo.ExtendConfig = po.ExtendConfig
		if newPo.ExtendConfig == "" {
			newPo.ExtendConfig = "{}"
		}
	}
	if err := commonschemalogic.CheckAffordance(&newPo.DmSchemaCore); err != nil {
		return nil, nil, err
	}
	return po, newPo, nil
}

// 更新产品物模型
func (l *ProductSchemaUpdateLogic) ProductSchemaUpdate(in *dm.ProductSchemaUpdateReq) (*dm.Empty, error) {
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
			l.ctx, relationDB.ToPropertyDo(&newPo.DmSchemaCore), in.Info.ProductID); err != nil {
			l.Errorf("%s.CreateProperty failure,err:%v", utils.FuncName(), err)
			return nil, errors.Database.AddDetail(err)
		}
	}
	err = l.PsDB.Update(l.ctx, newPo)
	if err != nil {
		return nil, err
	}
	//清除缓存
	err = l.svcCtx.SchemaRepo.SetData(l.ctx, in.Info.ProductID, nil)
	if err != nil {
		return nil, err
	}
	return &dm.Empty{}, nil
}
