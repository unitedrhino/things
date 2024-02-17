package productmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/events"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductSchemaCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
	PsDB *relationDB.ProductSchemaRepo
}

func NewProductSchemaCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductSchemaCreateLogic {
	return &ProductSchemaCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		PsDB:   relationDB.NewProductSchemaRepo(ctx),
	}
}

func (l *ProductSchemaCreateLogic) ruleCheck(in *dm.ProductSchemaCreateReq) (*relationDB.DmProductSchema, error) {
	_, err := l.PiDB.FindOneByFilter(l.ctx, relationDB.ProductFilter{ProductIDs: []string{in.Info.ProductID}})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddMsgf("找不到该产品:%v", in.Info.ProductID)
		}
		return nil, err
	}
	_, err = l.PsDB.FindOneByFilter(l.ctx, relationDB.ProductSchemaFilter{
		ProductID: in.Info.ProductID, Identifiers: []string{in.Info.Identifier},
	})
	if err == nil {
		return nil, errors.Parameter.AddMsgf("标识符在该产品中已经被使用:%s", in.Info.Identifier)
	}
	if err != nil {
		if !errors.Cmp(err, errors.NotFind) {
			return nil, err
		}
	}
	po := ToProductSchemaPo(in.Info)

	var cs *relationDB.DmCommonSchema
	if in.Info.Tag != int64(schema.TagCustom) {
		cs, err = relationDB.NewCommonSchemaRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.CommonSchemaFilter{Identifiers: []string{in.Info.Identifier}})
		if err != nil {
			return nil, err
		}
		po.IsCanSceneLinkage = cs.IsCanSceneLinkage
		po.IsShareAuthPerm = cs.IsShareAuthPerm
		po.IsHistory = cs.IsHistory
	}
	if po.Name == "" {
		if cs == nil {
			return nil, errors.Parameter.AddMsg("功能名称不能为空")
		}
		po.Name = cs.Name
	}
	if po.Required == 0 && cs != nil {
		po.Required = cs.Required
	}

	if po.IsCanSceneLinkage == 0 && cs != nil {
		po.IsCanSceneLinkage = cs.IsCanSceneLinkage
	}
	if po.IsShareAuthPerm == 0 && cs != nil {
		po.IsShareAuthPerm = cs.IsShareAuthPerm
	}
	if po.IsHistory == 0 && cs != nil {
		po.IsHistory = cs.IsHistory
	}
	if po.Order == 0 && cs != nil {
		po.Order = cs.Order
	}
	if po.ExtendConfig == "" && cs != nil {
		po.ExtendConfig = cs.ExtendConfig
	}

	if err = CheckAffordance(&po.DmSchemaCore, cs); err != nil {
		return nil, err
	}
	return po, nil
}

// 新增产品
func (l *ProductSchemaCreateLogic) ProductSchemaCreate(in *dm.ProductSchemaCreateReq) (*dm.Empty, error) {
	l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(in))
	po, err := l.ruleCheck(in)
	if err != nil {
		l.Errorf("%s.ruleCheck err:%v", err)
		return nil, err
	}

	if schema.AffordanceType(po.Type) == schema.AffordanceTypeProperty && po.Tag == int64(schema.TagCustom) {
		if err := l.svcCtx.SchemaManaRepo.CreateProperty(l.ctx, relationDB.ToPropertyDo(&po.DmSchemaCore), in.Info.ProductID); err != nil {
			l.Errorf("%s.CreateProperty failure,err:%v", utils.FuncName(), err)
			return nil, errors.Database.AddDetail(err)
		}
	}
	err = l.PsDB.Insert(l.ctx, po)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.ServerMsg.Publish(l.ctx, eventBus.DmProductSchemaUpdate, &events.DeviceUpdateInfo{ProductID: in.Info.ProductID})
	if err != nil {
		return nil, err
	}
	return &dm.Empty{}, nil
}
