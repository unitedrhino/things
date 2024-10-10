package info

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.ProductInfoDeleteReq) error {
	_, err := l.svcCtx.ProductM.ProductInfoDelete(l.ctx, &dm.ProductInfoDeleteReq{ProductID: req.ProductID})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ManageProduct req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
