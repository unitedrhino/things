package logic

import (
	"context"
	"github.com/i-Things/things/src/stocksvr/internal/models/stock_move"
	"google.golang.org/grpc/status"

	"github.com/i-Things/things/src/stocksvr/internal/svc"
	"github.com/i-Things/things/src/stocksvr/types/pb/stock"

	"github.com/zeromicro/go-zero/core/logx"
)

type StockMoveBatchCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStockMoveBatchCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StockMoveBatchCreateLogic {
	return &StockMoveBatchCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StockMoveBatchCreateLogic) StockMoveBatchCreate(in *stock.ReqStockMoveBatchCreate) (out *stock.ResStockMoveBatchCreate, err error) {
	out = &stock.ResStockMoveBatchCreate{}
	if err = l.checkReq(in); err != nil {
		return
	}
	var entities []map[string]interface{}
	model := l.svcCtx.Model.StockMoveModel
	for _, v := range in.Items {
		entities = append(entities, map[string]interface{}{
			stock_move.Columns.Sequence:             v.Sequence,             // Sequence
			stock_move.Columns.CompanyId:            v.CompanyId,            // Company
			stock_move.Columns.ProductId:            v.ProductId,            // Product
			stock_move.Columns.ProductUom:           v.ProductUom,           // UoM
			stock_move.Columns.LocationId:           v.LocationId,           // Source Location
			stock_move.Columns.LocationDestId:       v.LocationDestId,       // Destination Location
			stock_move.Columns.PartnerId:            v.PartnerId,            // Destination Address
			stock_move.Columns.PickingId:            v.PickingId,            // Transfer
			stock_move.Columns.GroupId:              v.GroupId,              // Procurement Group
			stock_move.Columns.RuleId:               v.RuleId,               // Stock Rule
			stock_move.Columns.PickingTypeId:        v.PickingTypeId,        // Operation Type
			stock_move.Columns.OriginReturnedMoveId: v.OriginReturnedMoveId, // Origin return move
			stock_move.Columns.RestrictPartnerId:    v.RestrictPartnerId,    // Owner
			stock_move.Columns.WarehouseId:          v.WarehouseId,          // Warehouse
			stock_move.Columns.PackageLevelId:       v.PackageLevelId,       // Package Level
			stock_move.Columns.NextSerialCount:      v.NextSerialCount,      // Number of SN
			stock_move.Columns.OrderpointId:         v.OrderpointId,         // Original Reordering Rule
			stock_move.Columns.ProductPackagingId:   v.ProductPackagingId,   // Packaging
			stock_move.Columns.CreateUid:            v.CreateUid,            // Created by
			stock_move.Columns.WriteUid:             v.WriteUid,             // Last Updated by
			stock_move.Columns.Name:                 v.Name,                 // Description
			stock_move.Columns.Priority:             v.Priority,             // Priority
			stock_move.Columns.State:                v.State,                // Status
			stock_move.Columns.Origin:               v.Origin,               // Source Document
			stock_move.Columns.ProcureMethod:        v.ProcureMethod,        // Supply Method
			stock_move.Columns.Reference:            v.Reference,            // Reference
			stock_move.Columns.NextSerial:           v.NextSerial,           // First SN
			//stock_move.Columns.ReservationDate:       gtime.NewFromStr(v.ReservationDate), // Date to Reserve
			stock_move.Columns.DescriptionPicking: v.DescriptionPicking, // Description of Picking
			stock_move.Columns.ProductQty:         v.ProductQty,         // Real Quantity
			stock_move.Columns.ProductUomQty:      v.ProductUomQty,      // Demand
			stock_move.Columns.QuantityDone:       v.QuantityDone,       // Quantity Done
			stock_move.Columns.Scrapped:           v.Scrapped,           // Scrapped
			stock_move.Columns.PropagateCancel:    v.PropagateCancel,    // Propagate cancel and split
			stock_move.Columns.IsInventory:        v.IsInventory,        // Inventory
			stock_move.Columns.Additional:         v.Additional,         // Whether the move was added after the picking's confirmation
			//stock_move.Columns.Date:                  gtime.NewFromStr(v.Date),            // Date Scheduled
			//stock_move.Columns.DateDeadline:          gtime.NewFromStr(v.DateDeadline),    // Deadline
			//stock_move.Columns.DelayAlertDate:        gtime.NewFromStr(v.DelayAlertDate),  // Delay Alert Date
			//stock_move.Columns.CreateDate:            gtime.NewFromStr(v.CreateDate),      // Created on
			//stock_move.Columns.WriteDate:             gtime.NewFromStr(v.WriteDate),       // Last Updated on
			stock_move.Columns.PriceUnit:             v.PriceUnit,             // Unit Price
			stock_move.Columns.AnalyticAccountLineId: v.AnalyticAccountLineId, // Analytic Account Line
			stock_move.Columns.ToRefund:              v.ToRefund,              // Update quantities on SO/PO
			stock_move.Columns.SaleLineId:            v.SaleLineId,            // Sale Line
			stock_move.Columns.PurchaseLineId:        v.PurchaseLineId,        // Purchase Order Line
			stock_move.Columns.CreatedPurchaseLineId: v.CreatedPurchaseLineId, // Created Purchase Order Line

		})
	}
	result, err := model.BatchInsert(entities)
	if err != nil {
		return
	}
	id, _ := result.LastInsertId()
	out = &stock.ResStockMoveBatchCreate{Id: uint64(id)}
	return
}

func (l *StockMoveBatchCreateLogic) checkReq(in *stock.ReqStockMoveBatchCreate) (err error) {
	if in.UserId == 0 {
		err = status.Error(10002, "用户不存在")
		return
	}
	for _, v := range in.Items {
		if v.Sequence == 0 {
			return status.Error(56439, "Sequence有误")
		}
		if v.CompanyId == 0 {
			return status.Error(46752, "Company有误")
		}
		if v.ProductId == 0 {
			return status.Error(82222, "Product有误")
		}
		if v.ProductUom == 0 {
			return status.Error(82748, "UoM有误")
		}
		if v.LocationId == 0 {
			return status.Error(69407, "Source Location有误")
		}
		if v.LocationDestId == 0 {
			return status.Error(38875, "Destination Location有误")
		}
		if v.PartnerId == 0 {
			return status.Error(28353, "Destination Address有误")
		}
		if v.PickingId == 0 {
			return status.Error(88167, "Transfer有误")
		}
		if v.GroupId == 0 {
			return status.Error(72966, "Procurement Group有误")
		}
		if v.RuleId == 0 {
			return status.Error(60054, "Stock Rule有误")
		}
		if v.PickingTypeId == 0 {
			return status.Error(12449, "Operation Type有误")
		}
		if v.OriginReturnedMoveId == 0 {
			return status.Error(20405, "Origin return move有误")
		}
		if v.RestrictPartnerId == 0 {
			return status.Error(54339, "Owner有误")
		}
		if v.WarehouseId == 0 {
			return status.Error(18076, "Warehouse有误")
		}
		if v.PackageLevelId == 0 {
			return status.Error(34402, "Package Level有误")
		}
		if v.NextSerialCount == 0 {
			return status.Error(39238, "Number of SN有误")
		}
		if v.OrderpointId == 0 {
			return status.Error(41574, "Original Reordering Rule有误")
		}
		if v.ProductPackagingId == 0 {
			return status.Error(14365, "Packaging有误")
		}
		if v.CreateUid == 0 {
			return status.Error(75155, "Created by有误")
		}
		if v.WriteUid == 0 {
			return status.Error(99959, "Last Updated by有误")
		}
		if v.Name == "" {
			return status.Error(37216, "Description有误")
		}
		if v.Priority == "" {
			return status.Error(56738, "Priority有误")
		}
		if v.State == "" {
			return status.Error(15654, "Status有误")
		}
		if v.Origin == "" {
			return status.Error(75115, "Source Document有误")
		}
		if v.ProcureMethod == "" {
			return status.Error(23912, "Supply Method有误")
		}
		if v.Reference == "" {
			return status.Error(37234, "Reference有误")
		}
		if v.NextSerial == "" {
			return status.Error(40095, "First SN有误")
		}
		if v.DescriptionPicking == "" {
			return status.Error(78540, "Description of Picking有误")
		}
		if v.ProductQty == 0 {
			return status.Error(78062, "Real Quantity有误")
		}
		if v.ProductUomQty == 0 {
			return status.Error(97064, "Demand有误")
		}
		if v.QuantityDone == 0 {
			return status.Error(70506, "Quantity Done有误")
		}
		if v.Scrapped {
			return status.Error(63094, "Scrapped有误")
		}
		if v.PropagateCancel {
			return status.Error(10726, "Propagate cancel and split有误")
		}
		if v.IsInventory {
			return status.Error(63310, "Inventory有误")
		}
		if v.Additional {
			return status.Error(27718, "Whether the move was added after the picking's confirmation有误")
		}
		if v.PriceUnit == 0 {
			return status.Error(64171, "Unit Price有误")
		}
		if v.AnalyticAccountLineId == 0 {
			return status.Error(76750, "Analytic Account Line有误")
		}
		if v.ToRefund {
			return status.Error(58158, "Update quantities on SO/PO有误")
		}
		if v.SaleLineId == 0 {
			return status.Error(67992, "Sale Line有误")
		}
		if v.PurchaseLineId == 0 {
			return status.Error(45269, "Purchase Order Line有误")
		}
		if v.CreatedPurchaseLineId == 0 {
			return status.Error(47913, "Created Purchase Order Line有误")
		}

	}
	return
}
