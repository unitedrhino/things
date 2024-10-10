package category

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.ProductCategory) error {
	_, err := l.svcCtx.ProductM.ProductCategoryUpdate(l.ctx, productCategoryToRpc(req))
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ProductCategoryUpdate req=%v err=%+v", utils.FuncName(), utils.Fmt(req), er)
		return er
	}
	return nil
}
