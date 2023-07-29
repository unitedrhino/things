package logic

import (
	"context"
	"github.com/i-Things/things/src/stocksvr/internal/models/stock_move"
	"github.com/jinzhu/copier"
	"google.golang.org/grpc/status"

	"github.com/i-Things/things/src/stocksvr/internal/svc"
	"github.com/i-Things/things/src/stocksvr/types/pb/stock"

	"github.com/zeromicro/go-zero/core/logx"
)

type StockMoveCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStockMoveCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StockMoveCreateLogic {
	return &StockMoveCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// StockMoveCreate 新增
func (l *StockMoveCreateLogic) StockMoveCreate(in *stock.ReqStockMoveCreate) (out *stock.ResStockMoveCreate, err error) {
	out = &stock.ResStockMoveCreate{}
	if err = l.checkReq(in); err != nil {
		return
	}
	model := l.svcCtx.Model.StockMoveModel
	dto := stock_move.DtoCreate{}
	_ = copier.Copy(&dto, in)

	lastInsertId, err := model.Create(dto)
	if err != nil {
		return
	}
	out = &stock.ResStockMoveCreate{Id: uint64(lastInsertId)}
	out.Id = uint64(lastInsertId)
	return
}

func (l *StockMoveCreateLogic) checkReq(in *stock.ReqStockMoveCreate) (err error) {

	if in.UserId == 0 {
		err = status.Error(10002, "用户不存在")
		return
	}
	if in.Sequence == 0 {
		return status.Error(99857, "Sequence有误")
	}
	if in.CompanyId == 0 {
		return status.Error(47750, "Company有误")
	}
	if in.ProductId == 0 {
		return status.Error(22593, "Product有误")
	}
	if in.ProductUom == 0 {
		return status.Error(61653, "UoM有误")
	}
	if in.LocationId == 0 {
		return status.Error(26672, "Source Location有误")
	}
	if in.LocationDestId == 0 {
		return status.Error(48414, "Destination Location有误")
	}
	if in.PartnerId == 0 {
		return status.Error(52320, "Destination Address有误")
	}
	if in.PickingId == 0 {
		return status.Error(11834, "Transfer有误")
	}
	if in.GroupId == 0 {
		return status.Error(39603, "Procurement Group有误")
	}
	if in.RuleId == 0 {
		return status.Error(74958, "Stock Rule有误")
	}
	if in.PickingTypeId == 0 {
		return status.Error(14960, "Operation Type有误")
	}
	if in.OriginReturnedMoveId == 0 {
		return status.Error(11253, "Origin return move有误")
	}
	if in.RestrictPartnerId == 0 {
		return status.Error(18817, "Owner有误")
	}
	if in.WarehouseId == 0 {
		return status.Error(28982, "Warehouse有误")
	}
	if in.PackageLevelId == 0 {
		return status.Error(68680, "Package Level有误")
	}
	if in.NextSerialCount == 0 {
		return status.Error(84478, "Number of SN有误")
	}
	if in.OrderpointId == 0 {
		return status.Error(31717, "Original Reordering Rule有误")
	}
	if in.ProductPackagingId == 0 {
		return status.Error(43408, "Packaging有误")
	}
	if in.CreateUid == 0 {
		return status.Error(70019, "Created by有误")
	}
	if in.WriteUid == 0 {
		return status.Error(80250, "Last Updated by有误")
	}
	if in.Name == "" {
		return status.Error(35135, "Description有误")
	}
	if in.Priority == "" {
		return status.Error(29306, "Priority有误")
	}
	if in.State == "" {
		return status.Error(67424, "Status有误")
	}
	if in.Origin == "" {
		return status.Error(65242, "Source Document有误")
	}
	if in.ProcureMethod == "" {
		return status.Error(84091, "Supply Method有误")
	}
	if in.Reference == "" {
		return status.Error(60518, "Reference有误")
	}
	if in.NextSerial == "" {
		return status.Error(12678, "First SN有误")
	}
	if in.DescriptionPicking == "" {
		return status.Error(76828, "Description of Picking有误")
	}
	if in.ProductQty == 0 {
		return status.Error(83984, "Real Quantity有误")
	}
	if in.ProductUomQty == 0 {
		return status.Error(13755, "Demand有误")
	}
	if in.QuantityDone == 0 {
		return status.Error(77190, "Quantity Done有误")
	}
	if in.Scrapped {
		return status.Error(41318, "Scrapped有误")
	}
	if in.PropagateCancel {
		return status.Error(57407, "Propagate cancel and split有误")
	}
	if in.IsInventory {
		return status.Error(13138, "Inventory有误")
	}
	if in.Additional {
		return status.Error(95603, "Whether the move was added after the picking's confirmation有误")
	}
	if in.PriceUnit == 0 {
		return status.Error(85042, "Unit Price有误")
	}
	if in.AnalyticAccountLineId == 0 {
		return status.Error(68662, "Analytic Account Line有误")
	}
	if in.ToRefund {
		return status.Error(74295, "Update quantities on SO/PO有误")
	}
	if in.SaleLineId == 0 {
		return status.Error(94923, "Sale Line有误")
	}
	if in.PurchaseLineId == 0 {
		return status.Error(12635, "Purchase Order Line有误")
	}
	if in.CreatedPurchaseLineId == 0 {
		return status.Error(82377, "Created Purchase Order Line有误")
	}

	return
}
