package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/i-Things/things/src/stocksvr/internal/models/stock_location"
	"github.com/i-Things/things/src/stocksvr/internal/svc"
	"github.com/i-Things/things/src/stocksvr/types/pb/stock"

	"github.com/zeromicro/go-zero/core/logx"
)

const success = 1

type StockLocationDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStockLocationDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StockLocationDeleteLogic {
	return &StockLocationDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除
func (l *StockLocationDeleteLogic) StockLocationDelete(in *stock.ReqStockLocationDelete) (out *stock.ResStockLocationDelete, err error) {
	if err = l.checkReq(in); err != nil {
		return
	}
	model := l.svcCtx.Model.StockLocationModel
	updateBuilder := model.NewDeleteBuilder()
	updateBuilder = updateBuilder.Where(squirrel.Eq{stock_location.Columns.Id: in.Id})
	err = model.Delete(updateBuilder)
	if err != nil {
		return
	}
	return &stock.ResStockLocationDelete{Success: success}, nil
}

func (l *StockLocationDeleteLogic) checkReq(in *stock.ReqStockLocationDelete) (err error) {

	//if in.UserId == 0 {
	//	err = status.Error(10002, "用户不存在")
	//	return
	//}
	return
}
