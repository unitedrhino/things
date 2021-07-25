package logic

import (
	"context"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ManageProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewManageProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) ManageProductLogic {
	return ManageProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ManageProductLogic) ManageProduct(req types.ManageProductReq) (*types.ProductInfo, error) {
	// todo: add your logic here and delete this line

	return &types.ProductInfo{}, nil
}
