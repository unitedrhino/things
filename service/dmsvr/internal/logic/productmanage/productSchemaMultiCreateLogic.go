package productmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"golang.org/x/sync/errgroup"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

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
	for _, v := range in.List {
		info := v
		info.ProductID = in.ProductID
		info.Tag = schema.TagOptional
		errGroup.Go(func() error {
			defer utils.Recover(l.ctx)
			_, err := createLogic.ProductSchemaCreate(&dm.ProductSchemaCreateReq{Info: info})
			if err != nil && !errors.Cmp(errors.Duplicate, err) {
				return err
			}
			return nil
		})
	}
	err := errGroup.Wait()
	return &dm.Empty{}, err
}
