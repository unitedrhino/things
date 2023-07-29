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

type StockMoveBatchDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStockMoveBatchDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StockMoveBatchDeleteLogic {
	return &StockMoveBatchDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StockMoveBatchDeleteLogic) StockMoveBatchDelete(in *stock.ReqStockMoveBatchDelete) (out *stock.ResStockMoveBatchDelete, err error) {
	out = &stock.ResStockMoveBatchDelete{}
	if err = l.checkReq(in); err != nil {
		return
	}
	model := l.svcCtx.Model.StockMoveModel
	deleteBuilder := model.NewDeleteBuilder()
	deleteBuilder = deleteBuilder.Where(squirrel.Eq{stock_move.Columns.Id: in.Ids})
	err = model.Delete(deleteBuilder)
	if err != nil {
		return
	}
	return &stock.ResStockMoveBatchDelete{Success: success}, nil
}

func (l *StockMoveBatchDeleteLogic) checkReq(in *stock.ReqStockMoveBatchDelete) (err error) {

	if in.UserId == 0 {
		err = status.Error(10002, "用户不存在")
		return
	}
	return
}
