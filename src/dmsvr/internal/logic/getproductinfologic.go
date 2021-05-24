package logic

import (
	"context"

	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/dmsvr/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetProductInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductInfoLogic {
	return &GetProductInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProductInfoLogic) GetProductInfo(in *dm.GetProductInfoReq) (*dm.ProductInfo, error) {
	l.Infof("GetProductInfo|req=%+v",in)
	di,err := l.svcCtx.ProductInfo.FindOne(in.ProductID)
	if err != nil {
		return nil, err
	}
	return &dm.ProductInfo{
		ProductID:di.ProductID,
		ProductName: di.ProductName,
		CreatedTime: di.CreatedTime.Unix(),
	}, nil
}
