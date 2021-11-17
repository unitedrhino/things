package dm

import (
	"context"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/dmsvr/dm"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetProductInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetProductInfoLogic {
	return GetProductInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductInfoLogic) GetProductInfo(req types.GetProductInfoReq) (*types.GetProductInfoResp, error) {
	l.Infof("GetProductInfo|req=%+v", req)
	dmReq := &dm.GetProductInfoReq{
		ProductID: req.ProductID, //产品id
	}
	if req.Page != nil {
		if req.Page.PageSize == 0 || req.Page.Page == 0 {
			return nil, errors.Parameter.AddDetail("pageSize and page can't equal 0")
		}
		dmReq.Page = &dm.PageInfo{
			Page:     req.Page.Page,
			PageSize: req.Page.PageSize,
		}
	} else if req.ProductID == "" {
		return nil, errors.Parameter.AddDetail("need page or productID")
	}
	resp, err := l.svcCtx.DmRpc.GetProductInfo(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	pis := make([]*types.ProductInfo, 0, len(resp.Info))
	for _, v := range resp.Info {
		pi := types.ProductInfoToApi(v)
		pis = append(pis, pi)
	}
	return &types.GetProductInfoResp{
		Total: resp.Total,
		Info:  pis,
		Num:   int64(len(pis)),
	}, nil
	return &types.GetProductInfoResp{}, nil
}
