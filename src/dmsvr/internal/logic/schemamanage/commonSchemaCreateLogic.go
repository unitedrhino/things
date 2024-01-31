package schemamanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/domain/schema"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

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
func (l *CommonSchemaCreateLogic) CommonSchemaCreate(in *dm.CommonSchemaCreateReq) (*dm.Response, error) {
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
	return &dm.Response{}, nil
}
