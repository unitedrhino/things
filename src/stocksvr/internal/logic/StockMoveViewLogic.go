package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/i-Things/things/src/stocksvr/internal/models/stock_move"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/stocksvr/internal/svc"
	"github.com/i-Things/things/src/stocksvr/types/pb/stock"

	"github.com/zeromicro/go-zero/core/logx"
)

type StockMoveViewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStockMoveViewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StockMoveViewLogic {
	return &StockMoveViewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StockMoveViewLogic) StockMoveView(in *stock.ReqStockMoveView) (out *stock.ResStockMoveView, err error) {
	out = &stock.ResStockMoveView{}
	if err = l.checkReq(in); err != nil {
		return
	}
	model := l.svcCtx.Model.StockMoveModel
	selectBuilder := model.NewSelectBuilder().Where(squirrel.Eq{stock_move.Columns.Id: in.Id})
	err, entity := model.QueryEntity(selectBuilder, true)
	if err != nil {
		return
	}
	_ = copier.Copy(&out, entity)
	//out.CreatedAt = entity.CreatedAt.UnixNano() / 1e6
	//out.UpdatedAt = entity.UpdatedAt.UnixNano() / 1e6
	return
}
func (l *StockMoveViewLogic) checkReq(in *stock.ReqStockMoveView) (err error) {
	return
}
