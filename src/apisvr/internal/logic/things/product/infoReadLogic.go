package product

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoReadLogic {
	return &InfoReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoReadLogic) InfoRead(req *types.ProductInfoReadReq) (resp *types.ProductInfo, err error) {
	dmResp, err := l.svcCtx.DmRpc.ProductInfoRead(l.ctx, &dm.ProductInfoReadReq{ProductID: req.ProductID})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return productInfoToApi(dmResp), nil
}
