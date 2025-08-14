package productmanagelogic

import (
	"context"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"golang.org/x/sync/errgroup"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductSchemaMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductSchemaMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductSchemaMultiCreateLogic {
	return &ProductSchemaMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量新增物模型,只新增没有的,已有的不处理
func (l *ProductSchemaMultiCreateLogic) ProductSchemaMultiCreate(in *dm.ProductSchemaMultiCreateReq) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	createLogic := NewProductSchemaCreateLogic(l.ctx, l.svcCtx)
	var errGroup errgroup.Group
	var pos []*relationDB.DmSchemaInfo
	for _, v := range in.List {
		info := v
		info.ProductID = in.ProductID
		info.Tag = schema.TagOptional
		po, err := createLogic.RuleCheck(&dm.ProductSchemaCreateReq{Info: info})
		if err != nil {
			if !errors.Cmp(errors.Duplicate, err) {
				l.Errorf("RuleCheck err:%v", err)
				return nil, err
			}
			continue
		}
		pos = append(pos, po)
		errGroup.Go(func() error {
			if schema.AffordanceType(po.Type) == schema.AffordanceTypeProperty && po.Tag == int64(schema.TagCustom) {
				if err := l.svcCtx.SchemaManaRepo.CreateProperty(l.ctx, relationDB.ToPropertyDo(po.Identifier, &po.DmSchemaCore), po.ProductID); err != nil {
					l.Errorf("%s.CreateProperty failure,err:%v", utils.FuncName(), err)
					return errors.Database.AddDetail(err)
				}
			}
			return nil
		})
	}
	err := errGroup.Wait()
	if err != nil {
		return nil, err
	}
	if len(pos) == 0 {
		return &dm.Empty{}, err
	}
	err = relationDB.NewProductSchemaRepo(l.ctx).MultiInsert(l.ctx, pos)
	if err != nil {
		return nil, err
	}
	for _, po := range pos {
		//清除缓存
		err = l.svcCtx.ProductSchemaRepo.SetData(l.ctx, po.ProductID, nil)
		if err != nil {
			l.Errorf("%s.SetData failure,err:%v", utils.FuncName(), err)
		}
	}

	return &dm.Empty{}, nil
}
