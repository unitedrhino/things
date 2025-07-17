package productmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	commonschemalogic "gitee.com/unitedrhino/things/service/dmsvr/internal/logic/schemamanage"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/domain/schema"
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

func (l *ProductSchemaUpdateLogic) ruleCheck(in *dm.ProductSchemaUpdateReq) (*relationDB.DmSchemaInfo, *relationDB.DmSchemaInfo, error) {
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
	if in.Info.Tag != 0 && po.Tag != in.Info.Tag {
		return nil, nil, errors.Parameter.AddMsg("功能标签不支持修改")
	}
	newPo := logic.ToProductSchemaPo(in.Info)
	newPo.ID = po.ID
	if in.Info.Affordance != nil && po.Tag == schema.TagCustom {
		po.Affordance = newPo.Affordance
	}
	if in.Info.Name != nil {
		po.Name = newPo.Name
	}
	if in.Info.Desc != nil {
		po.Desc = newPo.Desc
	}
	if in.Info.Required != 0 {
		po.Required = newPo.Required
	}
	if in.Info.IsCanSceneLinkage != 0 && po.Tag == schema.TagCustom {
		po.IsCanSceneLinkage = newPo.IsCanSceneLinkage
	}
	if in.Info.FuncGroup != 0 && po.Tag == schema.TagCustom {
		po.FuncGroup = newPo.FuncGroup
	}
	if in.Info.ControlMode != 0 && po.Tag == schema.TagCustom {
		po.ControlMode = newPo.ControlMode
	}
	if in.Info.UserPerm != 0 && po.Tag == schema.TagCustom {
		po.UserPerm = newPo.UserPerm
	}

	if in.Info.IsHistory != 0 && po.Tag == schema.TagCustom {
		po.IsHistory = newPo.IsHistory
	}

	if in.Info.Order != 0 {
		po.Order = newPo.Order
	}
	if in.Info.IsPassword != 0 {
		po.IsPassword = newPo.IsPassword
	}
	if in.Info.ExtendConfig != "" {
		po.ExtendConfig = newPo.ExtendConfig
	}
	if err := commonschemalogic.CheckAffordance(po.Identifier, &newPo.DmSchemaCore); err != nil {
		return nil, nil, err
	}
	return po, newPo, nil
}

// 更新产品物模型
func (l *ProductSchemaUpdateLogic) ProductSchemaUpdate(in *dm.ProductSchemaUpdateReq) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	po, newPo, err := l.ruleCheck(in)
	if err != nil {
		return nil, err
	}
	if schema.AffordanceType(po.Type) == schema.AffordanceTypeProperty {
		if err := l.svcCtx.SchemaManaRepo.DeleteProperty(
			l.ctx, relationDB.ToPropertyDo(po.Identifier, &po.DmSchemaCore), in.Info.ProductID, in.Info.Identifier); err != nil {
			l.Errorf("%s.DeleteProperty failure,err:%v", utils.FuncName(), err)
			return nil, errors.Database.AddDetail(err)
		}
	}
	if schema.AffordanceType(newPo.Type) == schema.AffordanceTypeProperty {
		if err := l.svcCtx.SchemaManaRepo.CreateProperty(
			l.ctx, relationDB.ToPropertyDo(po.Identifier, &newPo.DmSchemaCore), in.Info.ProductID); err != nil {
			l.Errorf("%s.CreateProperty failure,err:%v", utils.FuncName(), err)
			return nil, errors.Database.AddDetail(err)
		}
	}
	err = l.PsDB.Update(l.ctx, newPo)
	if err != nil {
		return nil, err
	}
	//清除缓存
	err = l.svcCtx.ProductSchemaRepo.SetData(l.ctx, in.Info.ProductID, nil)
	if err != nil {
		return nil, err
	}
	return &dm.Empty{}, nil
}
