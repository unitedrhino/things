package logic

import (
	"context"
	"github.com/i-Things/things/shared/def"

	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/i-Things/things/src/dmsvr/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *GetProductInfoLogic) GetProductInfo(in *dm.GetProductInfoReq) (resp *dm.GetProductInfoResp, err error) {
	l.Infof("GetProductInfo|req=%+v", in)
	var info []*dm.ProductInfo
	var size int64
	if in.Page == nil || in.Page.Page == 0 {
		di, err := l.svcCtx.ProductInfo.FindOne(l.ctx, in.ProductID)
		if err != nil {
			return nil, err
		}
		info = append(info, ToProductInfo(di))
	} else {
		size, err = l.svcCtx.DmDB.GetCountByProductInfo(l.ctx)
		if err != nil {
			return nil, err
		}
		di, err := l.svcCtx.DmDB.FindByProductInfo(l.ctx, def.PageInfo{PageSize: in.Page.PageSize, Page: in.Page.Page})
		if err != nil {
			return nil, err
		}
		info = make([]*dm.ProductInfo, 0, len(di))
		for _, v := range di {
			info = append(info, ToProductInfo(v))
		}
	}
	return &dm.GetProductInfoResp{List: info, Total: size}, nil
}
