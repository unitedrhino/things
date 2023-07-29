package logic

import (
	"context"
	"github.com/i-Things/things/src/stocksvr/internal/svc"
	"github.com/i-Things/things/src/stocksvr/types/pb/stock"

	"github.com/zeromicro/go-zero/core/logx"
)

type StockMoveUpdateAutomaticPartialLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStockMoveUpdateAutomaticPartialLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StockMoveUpdateAutomaticPartialLogic {
	return &StockMoveUpdateAutomaticPartialLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// StockMoveUpdateAutomaticPartial 根据租户ID
func (l *StockMoveUpdateAutomaticPartialLogic) StockMoveUpdateAutomaticPartial(in *stock.ReqStockMoveUpdateAutomaticPartial) (out *stock.ResStockMoveUpdateAutomaticPartial, err error) {
	out = &stock.ResStockMoveUpdateAutomaticPartial{}
	//var dto stock_move.DtoUpdateAutomatic
	//_ = copier.Copy(&dto, in.Item)
	//
	//err = l.svcCtx.Model.StockMoveModel.UpdateAutomaticById(in.Id, dto, in.Fields)
	//if err != nil {
	//	err = status.Error(1997, err.Error())
	//}
	return &stock.ResStockMoveUpdateAutomaticPartial{Success: success}, nil
}
