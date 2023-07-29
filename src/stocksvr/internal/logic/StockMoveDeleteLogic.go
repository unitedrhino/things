package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/i-Things/things/src/stocksvr/internal/models/stock_move"
	"google.golang.org/grpc/status"

	"github.com/i-Things/things/src/stocksvr/internal/svc"
	"github.com/i-Things/things/src/stocksvr/types/pb/stock"

	"github.com/zeromicro/go-zero/core/logx"
)

type StockMoveDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStockMoveDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StockMoveDeleteLogic {
	return &StockMoveDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除
func (l *StockMoveDeleteLogic) StockMoveDelete(in *stock.ReqStockMoveDelete) (out *stock.ResStockMoveDelete, err error) {
	out = &stock.ResStockMoveDelete{}
	if err = l.checkReq(in); err != nil {
		return
	}
	model := l.svcCtx.Model.StockMoveModel
	deleteBuilder := model.NewDeleteBuilder()
	deleteBuilder = deleteBuilder.Where(squirrel.Eq{stock_move.Columns.Id: in.Id})
	err = model.Delete(deleteBuilder)
	if err != nil {
		return
	}
	return &stock.ResStockMoveDelete{Success: success}, nil
}

func (l *StockMoveDeleteLogic) checkReq(in *stock.ReqStockMoveDelete) (err error) {

	if in.UserId == 0 {
		err = status.Error(10002, "用户不存在")
		return
	}
	return
}
