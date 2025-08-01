package schemamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/domain/schema"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommonSchemaUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PsDB *relationDB.CommonSchemaRepo
}

func NewCommonSchemaUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommonSchemaUpdateLogic {
	return &CommonSchemaUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PsDB:   relationDB.NewCommonSchemaRepo(ctx),
	}
}

func (l *CommonSchemaUpdateLogic) ruleCheck(in *dm.CommonSchemaUpdateReq) (*relationDB.DmCommonSchema, *relationDB.DmCommonSchema, bool, error) {
	var schemaIsUpdate bool
	po, err := l.PsDB.FindOneByFilter(l.ctx, relationDB.CommonSchemaFilter{
		Identifiers: []string{in.Info.Identifier},
	})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, nil, schemaIsUpdate, errors.Parameter.AddMsgf("标识符不存在:" + in.Info.Identifier)
		}
		return nil, nil, schemaIsUpdate, err
	}
	newPo := ToCommonSchemaPo(in.Info)
	newPo.ID = po.ID
	if in.Info.Affordance == nil && in.Info.Affordance.Value != newPo.Affordance {
		newPo.Affordance = po.Affordance
		schemaIsUpdate = true
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
	if in.Info.Order == 0 {
		newPo.Order = po.Order
	}
	if in.Info.IsCanSceneLinkage == 0 {
		newPo.IsCanSceneLinkage = po.IsCanSceneLinkage
	}
	if in.Info.FuncGroup == 0 {
		newPo.FuncGroup = po.FuncGroup
	}
	if in.Info.ControlMode == 0 {
		newPo.ControlMode = po.ControlMode
	}
	if in.Info.UserPerm == 0 {
		newPo.UserPerm = po.UserPerm
	}
	if in.Info.RecordMode == 0 {
		newPo.RecordMode = po.RecordMode
	}
	if in.Info.IsPassword == 0 {
		newPo.IsPassword = po.IsPassword
	}
	if in.Info.ExtendConfig == "" {
		newPo.ExtendConfig = po.ExtendConfig
		if newPo.ExtendConfig == "" {
			newPo.ExtendConfig = "{}"
		}
	}
	if err := CheckAffordance(newPo.Identifier, &newPo.DmSchemaCore); err != nil {
		return nil, nil, schemaIsUpdate, err
	}
	return po, newPo, schemaIsUpdate, nil
}

// 更新产品物模型
func (l *CommonSchemaUpdateLogic) CommonSchemaUpdate(in *dm.CommonSchemaUpdateReq) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	po, newPo, schemaIsUpdate, err := l.ruleCheck(in)
	if err != nil {
		return nil, err
	}
	if schemaIsUpdate {
		if schema.AffordanceType(po.Type) == schema.AffordanceTypeProperty {
			if err := l.svcCtx.SchemaManaRepo.DeleteProperty(
				l.ctx, nil, "", in.Info.Identifier); err != nil {
				l.Errorf("%s.DeleteProperty failure,err:%v", utils.FuncName(), err)
				return nil, errors.Database.AddDetail(err)
			}
		}
		if schema.AffordanceType(newPo.Type) == schema.AffordanceTypeProperty {
			if err := l.svcCtx.SchemaManaRepo.CreateProperty(
				l.ctx, relationDB.ToPropertyDo(newPo.Identifier, &newPo.DmSchemaCore), ""); err != nil {
				l.Errorf("%s.CreateProperty failure,err:%v", utils.FuncName(), err)
				return nil, errors.Database.AddDetail(err)
			}
		}
	}
	err = l.PsDB.Update(l.ctx, newPo)
	if err != nil {
		return nil, err
	}
	err = relationDB.NewProductSchemaRepo(l.ctx).UpdateWithCommon(l.ctx, newPo)
	if err != nil {
		return nil, err
	}
	return &dm.Empty{}, nil
}
