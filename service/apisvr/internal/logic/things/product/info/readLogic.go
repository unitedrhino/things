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
	dmResp, err := l.svcCtx.ProductM.ProductInfoRead(l.ctx, &dm.ProductInfoReadReq{
		ProductID:    req.ProductID,
		WithProtocol: req.WithProtocol,
		WithCategory: req.WithCategory,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s rpc.GetDeviceInfo req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return productInfoToApi(l.ctx, dmResp), nil
}
