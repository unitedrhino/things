package logic

import (
	"context"
	"github.com/i-Things/things/src/stocksvr/internal/models/stock_location"
	"github.com/jinzhu/copier"
	"google.golang.org/grpc/status"

	"github.com/i-Things/things/src/stocksvr/internal/svc"
	"github.com/i-Things/things/src/stocksvr/types/pb/stock"

	"github.com/zeromicro/go-zero/core/logx"
)

type StockLocationCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStockLocationCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StockLocationCreateLogic {
	return &StockLocationCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// StockLocationCreate 新增
func (l *StockLocationCreateLogic) StockLocationCreate(in *stock.ReqStockLocationCreate) (out *stock.ResStockLocationCreate, err error) {
	if err = l.checkReq(in); err != nil {
		return
	}
	model := l.svcCtx.Model.StockLocationModel

	dto := stock_location.DtoCreate{}
	_ = copier.Copy(&dto, in)
	lastInsertId, err := model.Create(dto)

	out = &stock.ResStockLocationCreate{Id: uint64(lastInsertId)}
	out.Id = uint64(lastInsertId)
	return
}

func (l *StockLocationCreateLogic) checkReq(in *stock.ReqStockLocationCreate) (err error) {

	//if in.UserId == 0 {
	//	err = status.Error(10002, "用户不存在")
	//	return
	//}
	//if in.LocationId == 0 {
	//	return status.Error(87036, "Parent Location有误")
	//}
	//if in.Posx == 0 {
	//	return status.Error(42432, "Corridor (X)有误")
	//}
	//if in.Posy == 0 {
	//	return status.Error(89531, "Shelves (Y)有误")
	//}
	//if in.Posz == 0 {
	//	return status.Error(52649, "Height (Z)有误")
	//}
	//if in.CompanyId == 0 {
	//	return status.Error(28891, "Company有误")
	//}
	//if in.RemovalStrategyId == 0 {
	//	return status.Error(41358, "Removal Strategy有误")
	//}
	//if in.CyclicInventoryFrequency == 0 {
	//	return status.Error(14008, "Inventory Frequency (Days)有误")
	//}
	//if in.WarehouseId == 0 {
	//	return status.Error(60764, "Warehouse有误")
	//}
	//if in.StorageCategoryId == 0 {
	//	return status.Error(95635, "Storage Category有误")
	//}
	//if in.CreateUid == 0 {
	//	return status.Error(57359, "Created by有误")
	//}
	//if in.WriteUid == 0 {
	//	return status.Error(82957, "Last Updated by有误")
	//}
	if in.Name == "" {
		return status.Error(43057, "Location Name有误")
	}
	//if in.CompleteName == "" {
	//	return status.Error(72118, "Full Location Name有误")
	//}
	if in.Usage == "" {
		return status.Error(89965, "Location Type有误")
	}
	//if in.ParentPath == "" {
	//	return status.Error(45310, "Parent Path有误")
	//}
	//if in.Barcode == "" {
	//	return status.Error(77417, "Barcode有误")
	//}
	//if in.Comment == "" {
	//	return status.Error(50854, "Additional Information有误")
	//}
	//if in.Active {
	//	return status.Error(64635, "Active有误")
	//}
	//if in.ScrapLocation {
	//	return status.Error(26567, "Is a Scrap Location?有误")
	//}
	//if in.ReturnLocation {
	//	return status.Error(95129, "Is a Return Location?有误")
	//}
	//if in.ReplenishLocation {
	//	return status.Error(61373, "Replenish Location有误")
	//}
	//if in.ValuationInAccountId == 0 {
	//	return status.Error(73719, "Stock Valuation Account (Incoming)有误")
	//}
	//if in.ValuationOutAccountId == 0 {
	//	return status.Error(41185, "Stock Valuation Account (Outgoing)有误")
	//}

	return
}
