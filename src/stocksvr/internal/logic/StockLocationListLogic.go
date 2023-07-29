package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/tools/reqUtil"
	"github.com/i-Things/things/src/stocksvr/internal/models/stock_location"

	"github.com/i-Things/things/src/stocksvr/internal/svc"
	"github.com/i-Things/things/src/stocksvr/types/pb/stock"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type StockLocationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStockLocationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StockLocationListLogic {
	return &StockLocationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 规格列表
func (l *StockLocationListLogic) StockLocationList(in *stock.ReqStockLocationList) (out *stock.ResStockLocationList, err error) {
	out = &stock.ResStockLocationList{}
	model := l.svcCtx.Model.StockLocationModel
	//squirrel.StatementBuilderType{}
	selectBuilder := model.NewSelectBuilder()
	eqBuilder := squirrel.Eq{}
	//搜索条件拼装
	if in.Id > 0 {
		eqBuilder[stock_location.Columns.Id] = in.Id
	}
	if in.LocationId > 0 {
		eqBuilder[stock_location.Columns.LocationId] = in.LocationId
	}
	if in.Posx > 0 {
		eqBuilder[stock_location.Columns.Posx] = in.Posx
	}
	if in.Posy > 0 {
		eqBuilder[stock_location.Columns.Posy] = in.Posy
	}
	if in.Posz > 0 {
		eqBuilder[stock_location.Columns.Posz] = in.Posz
	}
	if in.CompanyId > 0 {
		eqBuilder[stock_location.Columns.CompanyId] = in.CompanyId
	}
	if in.RemovalStrategyId > 0 {
		eqBuilder[stock_location.Columns.RemovalStrategyId] = in.RemovalStrategyId
	}
	if in.CyclicInventoryFrequency > 0 {
		eqBuilder[stock_location.Columns.CyclicInventoryFrequency] = in.CyclicInventoryFrequency
	}
	if in.WarehouseId > 0 {
		eqBuilder[stock_location.Columns.WarehouseId] = in.WarehouseId
	}
	if in.StorageCategoryId > 0 {
		eqBuilder[stock_location.Columns.StorageCategoryId] = in.StorageCategoryId
	}
	if in.CreateUid > 0 {
		eqBuilder[stock_location.Columns.CreateUid] = in.CreateUid
	}
	if in.WriteUid > 0 {
		eqBuilder[stock_location.Columns.WriteUid] = in.WriteUid
	}
	if in.Name != "" {
		eqBuilder[stock_location.Columns.Name] = in.Name
	}
	if in.CompleteName != "" {
		eqBuilder[stock_location.Columns.CompleteName] = in.CompleteName
	}
	if in.Usage != "" {
		eqBuilder[stock_location.Columns.Usage] = in.Usage
	}
	if in.ParentPath != "" {
		eqBuilder[stock_location.Columns.ParentPath] = in.ParentPath
	}
	if in.Barcode != "" {
		eqBuilder[stock_location.Columns.Barcode] = in.Barcode
	}
	if in.Comment != "" {
		eqBuilder[stock_location.Columns.Comment] = in.Comment
	}
	if in.Active {
		eqBuilder[stock_location.Columns.Active] = in.Active
	}
	if in.ScrapLocation {
		eqBuilder[stock_location.Columns.ScrapLocation] = in.ScrapLocation
	}
	if in.ReturnLocation {
		eqBuilder[stock_location.Columns.ReturnLocation] = in.ReturnLocation
	}
	if in.ReplenishLocation {
		eqBuilder[stock_location.Columns.ReplenishLocation] = in.ReplenishLocation
	}
	if in.ValuationInAccountId > 0 {
		eqBuilder[stock_location.Columns.ValuationInAccountId] = in.ValuationInAccountId
	}
	if in.ValuationOutAccountId > 0 {
		eqBuilder[stock_location.Columns.ValuationOutAccountId] = in.ValuationOutAccountId
	}

	selectBuilder = selectBuilder.Where(eqBuilder)
	//order := reqUtil.FormOrderBy{OrderBy: in.OrderBy}
	//selectBuilder = selectBuilder.OrderBy(order.GetOrder("id desc",
	//	//允许的排序字段
	//	stock_location.Columns.Id, //id= asc -id=desc
	//	//cas_menus.Columns.UpdatedAt,
	//))
	formPage := reqUtil.FormPage{
		WithTotal:  in.WithTotal,
		WithNoPage: in.WithNoPage,
		Page:       int(in.Page),
		PerPage:    int(in.PerPage),
	}
	//if formPage.ShowPerPage(false, 10, 100) {
	//	selectBuilder = selectBuilder.Offset(formPage.XOffset()).Limit(formPage.XLimit())
	//}
	var items []*stock.StockLocationItem
	//查询数据
	if err, entities := model.QueryEntities(selectBuilder); err != nil {
		return nil, err
	} else {
		if len(entities) <= 0 {
			return out, err
		}
		for _, entity := range entities {
			var item stock.StockLocationItem
			_ = copier.Copy(&item, entity)
			item.CreateDate = entity.CreateDate.UnixNano() / 1e6
			//item. = entity.UpdatedAt.UnixNano() / 1e6
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
	return &stock.ResStockLocationList{Total: uint64(total), Items: items}, nil
}
