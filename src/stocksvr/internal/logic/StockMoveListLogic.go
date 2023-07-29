package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/tools/reqUtil"
	"github.com/i-Things/things/src/stocksvr/internal/models/stock_move"

	"github.com/i-Things/things/src/stocksvr/internal/svc"
	"github.com/i-Things/things/src/stocksvr/types/pb/stock"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type StockMoveListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStockMoveListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StockMoveListLogic {
	return &StockMoveListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 规格列表
func (l *StockMoveListLogic) StockMoveList(in *stock.ReqStockMoveList) (out *stock.ResStockMoveList, err error) {
	out = &stock.ResStockMoveList{}
	model := l.svcCtx.Model.StockMoveModel
	selectBuilder := model.NewSelectBuilder()
	eqBuilder := squirrel.Eq{}
	//搜索条件拼装
	if in.Id > 0 {
		eqBuilder[stock_move.Columns.Id] = in.Id
	}
	if in.Sequence > 0 {
		eqBuilder[stock_move.Columns.Sequence] = in.Sequence
	}
	if in.CompanyId > 0 {
		eqBuilder[stock_move.Columns.CompanyId] = in.CompanyId
	}
	if in.ProductId > 0 {
		eqBuilder[stock_move.Columns.ProductId] = in.ProductId
	}
	if in.ProductUom > 0 {
		eqBuilder[stock_move.Columns.ProductUom] = in.ProductUom
	}
	if in.LocationId > 0 {
		eqBuilder[stock_move.Columns.LocationId] = in.LocationId
	}
	if in.LocationDestId > 0 {
		eqBuilder[stock_move.Columns.LocationDestId] = in.LocationDestId
	}
	if in.PartnerId > 0 {
		eqBuilder[stock_move.Columns.PartnerId] = in.PartnerId
	}
	if in.PickingId > 0 {
		eqBuilder[stock_move.Columns.PickingId] = in.PickingId
	}
	if in.GroupId > 0 {
		eqBuilder[stock_move.Columns.GroupId] = in.GroupId
	}
	if in.RuleId > 0 {
		eqBuilder[stock_move.Columns.RuleId] = in.RuleId
	}
	if in.PickingTypeId > 0 {
		eqBuilder[stock_move.Columns.PickingTypeId] = in.PickingTypeId
	}
	if in.OriginReturnedMoveId > 0 {
		eqBuilder[stock_move.Columns.OriginReturnedMoveId] = in.OriginReturnedMoveId
	}
	if in.RestrictPartnerId > 0 {
		eqBuilder[stock_move.Columns.RestrictPartnerId] = in.RestrictPartnerId
	}
	if in.WarehouseId > 0 {
		eqBuilder[stock_move.Columns.WarehouseId] = in.WarehouseId
	}
	if in.PackageLevelId > 0 {
		eqBuilder[stock_move.Columns.PackageLevelId] = in.PackageLevelId
	}
	if in.NextSerialCount > 0 {
		eqBuilder[stock_move.Columns.NextSerialCount] = in.NextSerialCount
	}
	if in.OrderpointId > 0 {
		eqBuilder[stock_move.Columns.OrderpointId] = in.OrderpointId
	}
	if in.ProductPackagingId > 0 {
		eqBuilder[stock_move.Columns.ProductPackagingId] = in.ProductPackagingId
	}
	if in.CreateUid > 0 {
		eqBuilder[stock_move.Columns.CreateUid] = in.CreateUid
	}
	if in.WriteUid > 0 {
		eqBuilder[stock_move.Columns.WriteUid] = in.WriteUid
	}
	if in.Name != "" {
		eqBuilder[stock_move.Columns.Name] = in.Name
	}
	if in.Priority != "" {
		eqBuilder[stock_move.Columns.Priority] = in.Priority
	}
	if in.State != "" {
		eqBuilder[stock_move.Columns.State] = in.State
	}
	if in.Origin != "" {
		eqBuilder[stock_move.Columns.Origin] = in.Origin
	}
	if in.ProcureMethod != "" {
		eqBuilder[stock_move.Columns.ProcureMethod] = in.ProcureMethod
	}
	if in.Reference != "" {
		eqBuilder[stock_move.Columns.Reference] = in.Reference
	}
	if in.NextSerial != "" {
		eqBuilder[stock_move.Columns.NextSerial] = in.NextSerial
	}
	if in.DescriptionPicking != "" {
		eqBuilder[stock_move.Columns.DescriptionPicking] = in.DescriptionPicking
	}
	if in.ProductQty > 0 {
		eqBuilder[stock_move.Columns.ProductQty] = in.ProductQty
	}
	if in.ProductUomQty > 0 {
		eqBuilder[stock_move.Columns.ProductUomQty] = in.ProductUomQty
	}
	if in.QuantityDone > 0 {
		eqBuilder[stock_move.Columns.QuantityDone] = in.QuantityDone
	}
	if in.Scrapped {
		eqBuilder[stock_move.Columns.Scrapped] = in.Scrapped
	}
	if in.PropagateCancel {
		eqBuilder[stock_move.Columns.PropagateCancel] = in.PropagateCancel
	}
	if in.IsInventory {
		eqBuilder[stock_move.Columns.IsInventory] = in.IsInventory
	}
	if in.Additional {
		eqBuilder[stock_move.Columns.Additional] = in.Additional
	}
	if in.PriceUnit > 0 {
		eqBuilder[stock_move.Columns.PriceUnit] = in.PriceUnit
	}
	if in.AnalyticAccountLineId > 0 {
		eqBuilder[stock_move.Columns.AnalyticAccountLineId] = in.AnalyticAccountLineId
	}
	if in.ToRefund {
		eqBuilder[stock_move.Columns.ToRefund] = in.ToRefund
	}
	if in.SaleLineId > 0 {
		eqBuilder[stock_move.Columns.SaleLineId] = in.SaleLineId
	}
	if in.PurchaseLineId > 0 {
		eqBuilder[stock_move.Columns.PurchaseLineId] = in.PurchaseLineId
	}
	if in.CreatedPurchaseLineId > 0 {
		eqBuilder[stock_move.Columns.CreatedPurchaseLineId] = in.CreatedPurchaseLineId
	}

	selectBuilder = selectBuilder.Where(eqBuilder)
	order := reqUtil.FormOrderBy{OrderBy: in.OrderBy}
	selectBuilder = selectBuilder.OrderBy(order.GetOrder("id desc",
		//允许的排序字段
		stock_move.Columns.Id, //id= asc -id=desc
		//cas_menus.Columns.UpdatedAt,
	))
	formPage := reqUtil.FormPage{
		WithTotal:  in.WithTotal,
		WithNoPage: in.WithNoPage,
		Page:       int(in.Page),
		PerPage:    int(in.PerPage),
	}
	if formPage.ShowPerPage(false, 10, 100) {
		selectBuilder = selectBuilder.Offset(formPage.XOffset()).Limit(formPage.XLimit())
	}
	var items []*stock.StockMoveItem
	//查询数据
	if err, entities := model.QueryEntities(selectBuilder); err != nil {
		return nil, err
	} else {
		if len(entities) <= 0 {
			return out, err
		}
		for _, entity := range entities {
			var item stock.StockMoveItem
			_ = copier.Copy(&item, entity)
			//item.CreatedAt = entity.CreatedAt.UnixNano() / 1e6
			//item.UpdatedAt = entity.UpdatedAt.UnixNano() / 1e6
			items = append(items, &item)
		}

	}

	//确定统计数据
	total, err := formPage.ShowTotal(false, func() (i int, err error) {
		return model.Count(selectBuilder)
	}, len(items))

	if err != nil {
		return nil, err
	}
	return &stock.ResStockMoveList{Total: uint64(total), Items: items}, nil
}
