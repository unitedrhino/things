package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/i-Things/things/src/stocksvr/internal/models/stock_move"
	"google.golang.org/grpc/status"

	"github.com/i-Things/things/src/stocksvr/internal/svc"
	"github.com/i-Things/things/src/stocksvr/types/pb/stock"

	"github.com/i-Things/things/shared/tools/helpers"
	"github.com/zeromicro/go-zero/core/logx"
)

type StockMoveUpdatePartialLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStockMoveUpdatePartialLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StockMoveUpdatePartialLogic {
	return &StockMoveUpdatePartialLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StockMoveUpdatePartialLogic) StockMoveUpdatePartial(in *stock.ReqStockMoveUpdatePartial) (out *stock.ResStockMoveUpdatePartial, err error) {
	out = &stock.ResStockMoveUpdatePartial{}
	if err = l.checkReq(in); err != nil {
		return
	}
	model := l.svcCtx.Model.StockMoveModel
	updateBuilder := model.NewUpdateBuilder()
	updateBuilder = updateBuilder.Where(squirrel.Eq{stock_move.Columns.Id: in.Ids})
	var updateMap = make(map[string]interface{})

	//更新条件拼接
	if helpers.StrInArray(stock_move.Columns.Sequence, in.Fields) {
		updateMap[stock_move.Columns.Sequence] = in.Sequence
	}
	if helpers.StrInArray(stock_move.Columns.CompanyId, in.Fields) {
		updateMap[stock_move.Columns.CompanyId] = in.CompanyId
	}
	if helpers.StrInArray(stock_move.Columns.ProductId, in.Fields) {
		updateMap[stock_move.Columns.ProductId] = in.ProductId
	}
	if helpers.StrInArray(stock_move.Columns.ProductUom, in.Fields) {
		updateMap[stock_move.Columns.ProductUom] = in.ProductUom
	}
	if helpers.StrInArray(stock_move.Columns.LocationId, in.Fields) {
		updateMap[stock_move.Columns.LocationId] = in.LocationId
	}
	if helpers.StrInArray(stock_move.Columns.LocationDestId, in.Fields) {
		updateMap[stock_move.Columns.LocationDestId] = in.LocationDestId
	}
	if helpers.StrInArray(stock_move.Columns.PartnerId, in.Fields) {
		updateMap[stock_move.Columns.PartnerId] = in.PartnerId
	}
	if helpers.StrInArray(stock_move.Columns.PickingId, in.Fields) {
		updateMap[stock_move.Columns.PickingId] = in.PickingId
	}
	if helpers.StrInArray(stock_move.Columns.GroupId, in.Fields) {
		updateMap[stock_move.Columns.GroupId] = in.GroupId
	}
	if helpers.StrInArray(stock_move.Columns.RuleId, in.Fields) {
		updateMap[stock_move.Columns.RuleId] = in.RuleId
	}
	if helpers.StrInArray(stock_move.Columns.PickingTypeId, in.Fields) {
		updateMap[stock_move.Columns.PickingTypeId] = in.PickingTypeId
	}
	if helpers.StrInArray(stock_move.Columns.OriginReturnedMoveId, in.Fields) {
		updateMap[stock_move.Columns.OriginReturnedMoveId] = in.OriginReturnedMoveId
	}
	if helpers.StrInArray(stock_move.Columns.RestrictPartnerId, in.Fields) {
		updateMap[stock_move.Columns.RestrictPartnerId] = in.RestrictPartnerId
	}
	if helpers.StrInArray(stock_move.Columns.WarehouseId, in.Fields) {
		updateMap[stock_move.Columns.WarehouseId] = in.WarehouseId
	}
	if helpers.StrInArray(stock_move.Columns.PackageLevelId, in.Fields) {
		updateMap[stock_move.Columns.PackageLevelId] = in.PackageLevelId
	}
	if helpers.StrInArray(stock_move.Columns.NextSerialCount, in.Fields) {
		updateMap[stock_move.Columns.NextSerialCount] = in.NextSerialCount
	}
	if helpers.StrInArray(stock_move.Columns.OrderpointId, in.Fields) {
		updateMap[stock_move.Columns.OrderpointId] = in.OrderpointId
	}
	if helpers.StrInArray(stock_move.Columns.ProductPackagingId, in.Fields) {
		updateMap[stock_move.Columns.ProductPackagingId] = in.ProductPackagingId
	}
	if helpers.StrInArray(stock_move.Columns.CreateUid, in.Fields) {
		updateMap[stock_move.Columns.CreateUid] = in.CreateUid
	}
	if helpers.StrInArray(stock_move.Columns.WriteUid, in.Fields) {
		updateMap[stock_move.Columns.WriteUid] = in.WriteUid
	}
	if helpers.StrInArray(stock_move.Columns.Name, in.Fields) {
		updateMap[stock_move.Columns.Name] = in.Name
	}
	if helpers.StrInArray(stock_move.Columns.Priority, in.Fields) {
		updateMap[stock_move.Columns.Priority] = in.Priority
	}
	if helpers.StrInArray(stock_move.Columns.State, in.Fields) {
		updateMap[stock_move.Columns.State] = in.State
	}
	if helpers.StrInArray(stock_move.Columns.Origin, in.Fields) {
		updateMap[stock_move.Columns.Origin] = in.Origin
	}
	if helpers.StrInArray(stock_move.Columns.ProcureMethod, in.Fields) {
		updateMap[stock_move.Columns.ProcureMethod] = in.ProcureMethod
	}
	if helpers.StrInArray(stock_move.Columns.Reference, in.Fields) {
		updateMap[stock_move.Columns.Reference] = in.Reference
	}
	if helpers.StrInArray(stock_move.Columns.NextSerial, in.Fields) {
		updateMap[stock_move.Columns.NextSerial] = in.NextSerial
	}
	if helpers.StrInArray(stock_move.Columns.ReservationDate, in.Fields) {
		updateMap[stock_move.Columns.ReservationDate] = in.ReservationDate
	}
	if helpers.StrInArray(stock_move.Columns.DescriptionPicking, in.Fields) {
		updateMap[stock_move.Columns.DescriptionPicking] = in.DescriptionPicking
	}
	if helpers.StrInArray(stock_move.Columns.ProductQty, in.Fields) {
		updateMap[stock_move.Columns.ProductQty] = in.ProductQty
	}
	if helpers.StrInArray(stock_move.Columns.ProductUomQty, in.Fields) {
		updateMap[stock_move.Columns.ProductUomQty] = in.ProductUomQty
	}
	if helpers.StrInArray(stock_move.Columns.QuantityDone, in.Fields) {
		updateMap[stock_move.Columns.QuantityDone] = in.QuantityDone
	}
	if helpers.StrInArray(stock_move.Columns.Scrapped, in.Fields) {
		updateMap[stock_move.Columns.Scrapped] = in.Scrapped
	}
	if helpers.StrInArray(stock_move.Columns.PropagateCancel, in.Fields) {
		updateMap[stock_move.Columns.PropagateCancel] = in.PropagateCancel
	}
	if helpers.StrInArray(stock_move.Columns.IsInventory, in.Fields) {
		updateMap[stock_move.Columns.IsInventory] = in.IsInventory
	}
	if helpers.StrInArray(stock_move.Columns.Additional, in.Fields) {
		updateMap[stock_move.Columns.Additional] = in.Additional
	}
	if helpers.StrInArray(stock_move.Columns.Date, in.Fields) {
		updateMap[stock_move.Columns.Date] = in.Date
	}
	if helpers.StrInArray(stock_move.Columns.DateDeadline, in.Fields) {
		updateMap[stock_move.Columns.DateDeadline] = in.DateDeadline
	}
	if helpers.StrInArray(stock_move.Columns.DelayAlertDate, in.Fields) {
		updateMap[stock_move.Columns.DelayAlertDate] = in.DelayAlertDate
	}
	if helpers.StrInArray(stock_move.Columns.CreateDate, in.Fields) {
		updateMap[stock_move.Columns.CreateDate] = in.CreateDate
	}
	if helpers.StrInArray(stock_move.Columns.WriteDate, in.Fields) {
		updateMap[stock_move.Columns.WriteDate] = in.WriteDate
	}
	if helpers.StrInArray(stock_move.Columns.PriceUnit, in.Fields) {
		updateMap[stock_move.Columns.PriceUnit] = in.PriceUnit
	}
	if helpers.StrInArray(stock_move.Columns.AnalyticAccountLineId, in.Fields) {
		updateMap[stock_move.Columns.AnalyticAccountLineId] = in.AnalyticAccountLineId
	}
	if helpers.StrInArray(stock_move.Columns.ToRefund, in.Fields) {
		updateMap[stock_move.Columns.ToRefund] = in.ToRefund
	}
	if helpers.StrInArray(stock_move.Columns.SaleLineId, in.Fields) {
		updateMap[stock_move.Columns.SaleLineId] = in.SaleLineId
	}
	if helpers.StrInArray(stock_move.Columns.PurchaseLineId, in.Fields) {
		updateMap[stock_move.Columns.PurchaseLineId] = in.PurchaseLineId
	}
	if helpers.StrInArray(stock_move.Columns.CreatedPurchaseLineId, in.Fields) {
		updateMap[stock_move.Columns.CreatedPurchaseLineId] = in.CreatedPurchaseLineId
	}

	err = model.Update(updateBuilder, updateMap)
	if err != nil {
		return
	}
	return &stock.ResStockMoveUpdatePartial{Success: success}, nil
}

func (l *StockMoveUpdatePartialLogic) checkReq(in *stock.ReqStockMoveUpdatePartial) (err error) {

	if in.UserId == 0 {
		err = status.Error(10002, "用户不存在")
		return
	}
	return
}
