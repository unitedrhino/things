package logic

import (
	"context"
	"github.com/i-Things/things/src/stocksvr/internal/models/stock_move"
	"google.golang.org/grpc/status"

	"github.com/i-Things/things/src/stocksvr/internal/svc"
	"github.com/i-Things/things/src/stocksvr/types/pb/stock"
	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"
)

type StockMoveUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStockMoveUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StockMoveUpdateLogic {
	return &StockMoveUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 规格修改
func (l *StockMoveUpdateLogic) StockMoveUpdate(in *stock.ReqStockMoveUpdate) (out *stock.ResStockMoveUpdate, err error) {
	out = &stock.ResStockMoveUpdate{}
	if err = l.checkReq(in); err != nil {
		return
	}
	//批量更新主体
	model := l.svcCtx.Model.StockMoveModel

	dto := stock_move.Dto{}
	_ = copier.Copy(&dto, in)
	err = model.UpdateDto(dto)
	if err != nil {
		return
	}
	return &stock.ResStockMoveUpdate{Success: success}, nil
}
func (l *StockMoveUpdateLogic) checkReq(in *stock.ReqStockMoveUpdate) (err error) {

	if in.UserId == 0 {
		err = status.Error(10002, "用户不存在")
		return
	}
	return
}
