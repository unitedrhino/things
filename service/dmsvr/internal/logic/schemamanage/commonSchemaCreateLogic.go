package schemamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/domain/schema"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommonSchemaCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PsDB *relationDB.CommonSchemaRepo
}

func NewCommonSchemaCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommonSchemaCreateLogic {
	return &CommonSchemaCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PsDB:   relationDB.NewCommonSchemaRepo(ctx),
	}
}

func (l *CommonSchemaCreateLogic) ruleCheck(in *dm.CommonSchemaCreateReq) (*relationDB.DmCommonSchema, error) {

	_, err := l.PsDB.FindOneByFilter(l.ctx, relationDB.CommonSchemaFilter{
		Identifiers: []string{in.Info.Identifier},
	})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			po := ToCommonSchemaPo(in.Info)
			if po.Name == "" {
				return nil, errors.Parameter.AddMsg("功能名称不能为空")
			}

			if po.ExtendConfig == "" {
				po.ExtendConfig = "{}"
			}
			if err := CheckAffordance(&po.DmSchemaCore); err != nil {
				return nil, err
			}
			return po, nil
		}
		return nil, errors.Database.AddDetail(err)
	}
	return nil, errors.Parameter.AddMsgf("标识符在该产品中已经被使用:%s", in.Info.Identifier)
}

// 新增产品
func (l *CommonSchemaCreateLogic) CommonSchemaCreate(in *dm.CommonSchemaCreateReq) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	po, err := l.ruleCheck(in)
	if err != nil {
		l.Errorf("%s.ruleCheck err:%v", err)
		return nil, err
	}

	if schema.AffordanceType(po.Type) == schema.AffordanceTypeProperty {
		if err := l.svcCtx.SchemaManaRepo.CreateProperty(l.ctx, relationDB.ToPropertyDo(&po.DmSchemaCore), ""); err != nil {
			l.Errorf("%s.CreateProperty failure,err:%v", utils.FuncName(), err)
			return nil, errors.Database.AddDetail(err)
		}
	}
	err = l.PsDB.Insert(l.ctx, po)
	if err != nil {
		return nil, err
	}
	return &dm.Empty{}, nil
}
