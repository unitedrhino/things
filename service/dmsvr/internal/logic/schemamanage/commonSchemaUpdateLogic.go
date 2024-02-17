package schemamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

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

func (l *CommonSchemaUpdateLogic) ruleCheck(in *dm.CommonSchemaUpdateReq) (*relationDB.DmCommonSchema, *relationDB.DmCommonSchema, error) {
	po, err := l.PsDB.FindOneByFilter(l.ctx, relationDB.CommonSchemaFilter{
		Identifiers: []string{in.Info.Identifier},
	})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, nil, errors.Parameter.AddMsgf("标识符不存在:" + in.Info.Identifier)
		}
		return nil, nil, err
	}
	newPo := ToCommonSchemaPo(in.Info)
	newPo.ID = po.ID
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
	if in.Info.Order == 0 {
		newPo.Order = po.Order
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
	if in.Info.ExtendConfig == "" {
		newPo.ExtendConfig = po.ExtendConfig
		if newPo.ExtendConfig == "" {
			newPo.ExtendConfig = "{}"
		}
	}
	if err := CheckAffordance(&newPo.DmSchemaCore); err != nil {
		return nil, nil, err
	}
	return po, newPo, nil
}

// 更新产品物模型
func (l *CommonSchemaUpdateLogic) CommonSchemaUpdate(in *dm.CommonSchemaUpdateReq) (*dm.Empty, error) {
	po, newPo, err := l.ruleCheck(in)
	if err != nil {
		return nil, err
	}
	if schema.AffordanceType(po.Type) == schema.AffordanceTypeProperty {
		if err := l.svcCtx.SchemaManaRepo.DeleteProperty(
			l.ctx, "", in.Info.Identifier); err != nil {
			l.Errorf("%s.DeleteProperty failure,err:%v", utils.FuncName(), err)
			return nil, errors.Database.AddDetail(err)
		}
	}
	if schema.AffordanceType(newPo.Type) == schema.AffordanceTypeProperty {
		if err := l.svcCtx.SchemaManaRepo.CreateProperty(
			l.ctx, relationDB.ToPropertyDo(&newPo.DmSchemaCore), ""); err != nil {
			l.Errorf("%s.CreateProperty failure,err:%v", utils.FuncName(), err)
			return nil, errors.Database.AddDetail(err)
		}
	}
	err = l.PsDB.Update(l.ctx, newPo)
	if err != nil {
		return nil, err
	}
	return &dm.Empty{}, nil
}
