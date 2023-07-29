package logic

import (
	"context"
	"github.com/i-Things/things/src/stocksvr/internal/models/stock_location"
	"github.com/i-Things/things/src/stocksvr/internal/svc"
	"github.com/i-Things/things/src/stocksvr/types/pb/stock"
	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"
)

type StockLocationUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStockLocationUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StockLocationUpdateLogic {
	return &StockLocationUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 规格修改
func (l *StockLocationUpdateLogic) StockLocationUpdate(in *stock.ReqStockLocationUpdate) (out *stock.ResStockLocationUpdate, err error) {
	if err = l.checkReq(in); err != nil {
		return
	}
	//批量更新主体
	model := l.svcCtx.Model.StockLocationModel
	dto := stock_location.Dto{}
	_ = copier.Copy(&dto, in)
	err = model.UpdateDto(dto)
	if err != nil {
		return
	}
	return &stock.ResStockLocationUpdate{Success: success}, nil
}
func (l *StockLocationUpdateLogic) checkReq(in *stock.ReqStockLocationUpdate) (err error) {

	//if in.UserId == 0 {
	//	err = status.Error(10002, "用户不存在")
	//	return
	//}
	return
}
