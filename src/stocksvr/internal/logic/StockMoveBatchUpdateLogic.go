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

type StockMoveBatchUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStockMoveBatchUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StockMoveBatchUpdateLogic {
	return &StockMoveBatchUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StockMoveBatchUpdateLogic) StockMoveBatchUpdate(in *stock.ReqStockMoveBatchUpdate) (out *stock.ResStockMoveBatchUpdate, err error) {
	out = &stock.ResStockMoveBatchUpdate{}
	if err = l.checkReq(in); err != nil {
		return
	}
	model := l.svcCtx.Model.StockMoveModel
	updateBuilder := model.NewUpdateBuilder()
	updateBuilder = updateBuilder.Where(squirrel.Eq{stock_move.Columns.Id: in.Ids})
	err = model.
		Update(updateBuilder, map[string]interface{}{
			stock_move.Columns.Sequence:             in.Sequence,             // Sequence
			stock_move.Columns.CompanyId:            in.CompanyId,            // Company
			stock_move.Columns.ProductId:            in.ProductId,            // Product
			stock_move.Columns.ProductUom:           in.ProductUom,           // UoM
			stock_move.Columns.LocationId:           in.LocationId,           // Source Location
			stock_move.Columns.LocationDestId:       in.LocationDestId,       // Destination Location
			stock_move.Columns.PartnerId:            in.PartnerId,            // Destination Address
			stock_move.Columns.PickingId:            in.PickingId,            // Transfer
			stock_move.Columns.GroupId:              in.GroupId,              // Procurement Group
			stock_move.Columns.RuleId:               in.RuleId,               // Stock Rule
			stock_move.Columns.PickingTypeId:        in.PickingTypeId,        // Operation Type
			stock_move.Columns.OriginReturnedMoveId: in.OriginReturnedMoveId, // Origin return move
			stock_move.Columns.RestrictPartnerId:    in.RestrictPartnerId,    // Owner
			stock_move.Columns.WarehouseId:          in.WarehouseId,          // Warehouse
			stock_move.Columns.PackageLevelId:       in.PackageLevelId,       // Package Level
			stock_move.Columns.NextSerialCount:      in.NextSerialCount,      // Number of SN
			stock_move.Columns.OrderpointId:         in.OrderpointId,         // Original Reordering Rule
			stock_move.Columns.ProductPackagingId:   in.ProductPackagingId,   // Packaging
			stock_move.Columns.CreateUid:            in.CreateUid,            // Created by
			stock_move.Columns.WriteUid:             in.WriteUid,             // Last Updated by
			stock_move.Columns.Name:                 in.Name,                 // Description
			stock_move.Columns.Priority:             in.Priority,             // Priority
			stock_move.Columns.State:                in.State,                // Status
			stock_move.Columns.Origin:               in.Origin,               // Source Document
			stock_move.Columns.ProcureMethod:        in.ProcureMethod,        // Supply Method
			stock_move.Columns.Reference:            in.Reference,            // Reference
			stock_move.Columns.NextSerial:           in.NextSerial,           // First SN
			//stock_move.Columns.ReservationDate:       gtime.NewFromStr(in.ReservationDate), // Date to Reserve
			stock_move.Columns.DescriptionPicking: in.DescriptionPicking, // Description of Picking
			stock_move.Columns.ProductQty:         in.ProductQty,         // Real Quantity
			stock_move.Columns.ProductUomQty:      in.ProductUomQty,      // Demand
			stock_move.Columns.QuantityDone:       in.QuantityDone,       // Quantity Done
			stock_move.Columns.Scrapped:           in.Scrapped,           // Scrapped
			stock_move.Columns.PropagateCancel:    in.PropagateCancel,    // Propagate cancel and split
			stock_move.Columns.IsInventory:        in.IsInventory,        // Inventory
			stock_move.Columns.Additional:         in.Additional,         // Whether the move was added after the picking's confirmation
			//stock_move.Columns.Date:                  gtime.NewFromStr(in.Date),            // Date Scheduled
			//stock_move.Columns.DateDeadline:          gtime.NewFromStr(in.DateDeadline),    // Deadline
			//stock_move.Columns.DelayAlertDate:        gtime.NewFromStr(in.DelayAlertDate),  // Delay Alert Date
			//stock_move.Columns.CreateDate:            gtime.NewFromStr(in.CreateDate),      // Created on
			//stock_move.Columns.WriteDate:             gtime.NewFromStr(in.WriteDate),       // Last Updated on
			stock_move.Columns.PriceUnit:             in.PriceUnit,             // Unit Price
			stock_move.Columns.AnalyticAccountLineId: in.AnalyticAccountLineId, // Analytic Account Line
			stock_move.Columns.ToRefund:              in.ToRefund,              // Update quantities on SO/PO
			stock_move.Columns.SaleLineId:            in.SaleLineId,            // Sale Line
			stock_move.Columns.PurchaseLineId:        in.PurchaseLineId,        // Purchase Order Line
			stock_move.Columns.CreatedPurchaseLineId: in.CreatedPurchaseLineId, // Created Purchase Order Line

		})
	if err != nil {
		return
	}
	return &stock.ResStockMoveBatchUpdate{Success: success}, nil
}

func (l *StockMoveBatchUpdateLogic) checkReq(in *stock.ReqStockMoveBatchUpdate) (err error) {

	if in.UserId == 0 {
		err = status.Error(10002, "用户不存在")
		return
	}
	return
}
