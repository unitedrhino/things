package info

import (
	"context"

	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.ProductInfoReadReq) (resp *types.ProductInfo, err error) {
	dmResp, err := l.svcCtx.ProductM.ProductInfoRead(l.ctx,
		&dm.ProductInfoReadReq{ProductID: req.ProductID, WithProtocol: req.WithProtocol, WithCategory: req.WithCategory})
	return productInfoToApi(l.ctx, dmResp), err
}
