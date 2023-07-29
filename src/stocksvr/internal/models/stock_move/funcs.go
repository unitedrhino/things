package stock_move

import (
	"github.com/Masterminds/squirrel"
	"github.com/gogf/gf/os/gtime"
	"github.com/i-Things/things/shared/xzero/xsql"
)

type Dto struct {
	Id                    uint64      `json:"id"                       db:"id"`                       //
	Sequence              int64       `json:"sequence"                 db:"sequence"`                 // Sequence
	CompanyId             uint64      `json:"company_id"               db:"company_id"`               // Company
	ProductId             uint64      `json:"product_id"               db:"product_id"`               // Product
	ProductUom            int64       `json:"product_uom"              db:"product_uom"`              // UoM
	LocationId            uint64      `json:"location_id"              db:"location_id"`              // Source Location
	LocationDestId        uint64      `json:"location_dest_id"         db:"location_dest_id"`         // Destination Location
	PartnerId             uint64      `json:"partner_id"               db:"partner_id"`               // Destination Address
	PickingId             uint64      `json:"picking_id"               db:"picking_id"`               // Transfer
	GroupId               uint64      `json:"group_id"                 db:"group_id"`                 // Procurement Group
	RuleId                uint64      `json:"rule_id"                  db:"rule_id"`                  // Stock Rule
	PickingTypeId         uint64      `json:"picking_type_id"          db:"picking_type_id"`          // Operation Type
	OriginReturnedMoveId  uint64      `json:"origin_returned_move_id"  db:"origin_returned_move_id"`  // Origin return move
	RestrictPartnerId     uint64      `json:"restrict_partner_id"      db:"restrict_partner_id"`      // Owner
	WarehouseId           uint64      `json:"warehouse_id"             db:"warehouse_id"`             // Warehouse
	PackageLevelId        uint64      `json:"package_level_id"         db:"package_level_id"`         // Package Level
	NextSerialCount       int64       `json:"next_serial_count"        db:"next_serial_count"`        // Number of SN
	OrderpointId          uint64      `json:"orderpoint_id"            db:"orderpoint_id"`            // Original Reordering Rule
	ProductPackagingId    uint64      `json:"product_packaging_id"     db:"product_packaging_id"`     // Packaging
	CreateUid             uint64      `json:"create_uid"               db:"create_uid"`               // Created by
	WriteUid              uint64      `json:"write_uid"                db:"write_uid"`                // Last Updated by
	Name                  string      `json:"name"                     db:"name"`                     // Description
	Priority              string      `json:"priority"                 db:"priority"`                 // Priority
	State                 string      `json:"state"                    db:"state"`                    // Status
	Origin                string      `json:"origin"                   db:"origin"`                   // Source Document
	ProcureMethod         string      `json:"procure_method"           db:"procure_method"`           // Supply Method
	Reference             string      `json:"reference"                db:"reference"`                // Reference
	NextSerial            string      `json:"next_serial"              db:"next_serial"`              // First SN
	ReservationDate       *gtime.Time `json:"reservation_date"         db:"reservation_date"`         // Date to Reserve
	DescriptionPicking    string      `json:"description_picking"      db:"description_picking"`      // Description of Picking
	ProductQty            float64     `json:"product_qty"              db:"product_qty"`              // Real Quantity
	ProductUomQty         float64     `json:"product_uom_qty"          db:"product_uom_qty"`          // Demand
	QuantityDone          float64     `json:"quantity_done"            db:"quantity_done"`            // Quantity Done
	Scrapped              bool        `json:"scrapped"                 db:"scrapped"`                 // Scrapped
	PropagateCancel       bool        `json:"propagate_cancel"         db:"propagate_cancel"`         // Propagate cancel and split
	IsInventory           bool        `json:"is_inventory"             db:"is_inventory"`             // Inventory
	Additional            bool        `json:"additional"               db:"additional"`               // Whether the move was added after the picking's confirmation
	Date                  *gtime.Time `json:"date"                     db:"date"`                     // Date Scheduled
	DateDeadline          *gtime.Time `json:"date_deadline"            db:"date_deadline"`            // Deadline
	DelayAlertDate        *gtime.Time `json:"delay_alert_date"         db:"delay_alert_date"`         // Delay Alert Date
	CreateDate            *gtime.Time `json:"create_date"              db:"create_date"`              // Created on
	WriteDate             *gtime.Time `json:"write_date"               db:"write_date"`               // Last Updated on
	PriceUnit             float64     `json:"price_unit"               db:"price_unit"`               // Unit Price
	AnalyticAccountLineId uint64      `json:"analytic_account_line_id" db:"analytic_account_line_id"` // Analytic Account Line
	ToRefund              bool        `json:"to_refund"                db:"to_refund"`                // Update quantities on SO/PO
	SaleLineId            uint64      `json:"sale_line_id"             db:"sale_line_id"`             // Sale Line
	PurchaseLineId        uint64      `json:"purchase_line_id"         db:"purchase_line_id"`         // Purchase Order Line
	CreatedPurchaseLineId uint64      `json:"created_purchase_line_id" db:"created_purchase_line_id"` // Created Purchase Order Line

}
type DtoCreate struct {
	Sequence              int64       `json:"sequence"                 db:"sequence"`                 // Sequence
	CompanyId             uint64      `json:"company_id"               db:"company_id"`               // Company
	ProductId             uint64      `json:"product_id"               db:"product_id"`               // Product
	ProductUom            int64       `json:"product_uom"              db:"product_uom"`              // UoM
	LocationId            uint64      `json:"location_id"              db:"location_id"`              // Source Location
	LocationDestId        uint64      `json:"location_dest_id"         db:"location_dest_id"`         // Destination Location
	PartnerId             uint64      `json:"partner_id"               db:"partner_id"`               // Destination Address
	PickingId             uint64      `json:"picking_id"               db:"picking_id"`               // Transfer
	GroupId               uint64      `json:"group_id"                 db:"group_id"`                 // Procurement Group
	RuleId                uint64      `json:"rule_id"                  db:"rule_id"`                  // Stock Rule
	PickingTypeId         uint64      `json:"picking_type_id"          db:"picking_type_id"`          // Operation Type
	OriginReturnedMoveId  uint64      `json:"origin_returned_move_id"  db:"origin_returned_move_id"`  // Origin return move
	RestrictPartnerId     uint64      `json:"restrict_partner_id"      db:"restrict_partner_id"`      // Owner
	WarehouseId           uint64      `json:"warehouse_id"             db:"warehouse_id"`             // Warehouse
	PackageLevelId        uint64      `json:"package_level_id"         db:"package_level_id"`         // Package Level
	NextSerialCount       int64       `json:"next_serial_count"        db:"next_serial_count"`        // Number of SN
	OrderpointId          uint64      `json:"orderpoint_id"            db:"orderpoint_id"`            // Original Reordering Rule
	ProductPackagingId    uint64      `json:"product_packaging_id"     db:"product_packaging_id"`     // Packaging
	CreateUid             uint64      `json:"create_uid"               db:"create_uid"`               // Created by
	WriteUid              uint64      `json:"write_uid"                db:"write_uid"`                // Last Updated by
	Name                  string      `json:"name"                     db:"name"`                     // Description
	Priority              string      `json:"priority"                 db:"priority"`                 // Priority
	State                 string      `json:"state"                    db:"state"`                    // Status
	Origin                string      `json:"origin"                   db:"origin"`                   // Source Document
	ProcureMethod         string      `json:"procure_method"           db:"procure_method"`           // Supply Method
	Reference             string      `json:"reference"                db:"reference"`                // Reference
	NextSerial            string      `json:"next_serial"              db:"next_serial"`              // First SN
	ReservationDate       *gtime.Time `json:"reservation_date"         db:"reservation_date"`         // Date to Reserve
	DescriptionPicking    string      `json:"description_picking"      db:"description_picking"`      // Description of Picking
	ProductQty            float64     `json:"product_qty"              db:"product_qty"`              // Real Quantity
	ProductUomQty         float64     `json:"product_uom_qty"          db:"product_uom_qty"`          // Demand
	QuantityDone          float64     `json:"quantity_done"            db:"quantity_done"`            // Quantity Done
	Scrapped              bool        `json:"scrapped"                 db:"scrapped"`                 // Scrapped
	PropagateCancel       bool        `json:"propagate_cancel"         db:"propagate_cancel"`         // Propagate cancel and split
	IsInventory           bool        `json:"is_inventory"             db:"is_inventory"`             // Inventory
	Additional            bool        `json:"additional"               db:"additional"`               // Whether the move was added after the picking's confirmation
	Date                  *gtime.Time `json:"date"                     db:"date"`                     // Date Scheduled
	DateDeadline          *gtime.Time `json:"date_deadline"            db:"date_deadline"`            // Deadline
	DelayAlertDate        *gtime.Time `json:"delay_alert_date"         db:"delay_alert_date"`         // Delay Alert Date
	CreateDate            *gtime.Time `json:"create_date"              db:"create_date"`              // Created on
	WriteDate             *gtime.Time `json:"write_date"               db:"write_date"`               // Last Updated on
	PriceUnit             float64     `json:"price_unit"               db:"price_unit"`               // Unit Price
	AnalyticAccountLineId uint64      `json:"analytic_account_line_id" db:"analytic_account_line_id"` // Analytic Account Line
	ToRefund              bool        `json:"to_refund"                db:"to_refund"`                // Update quantities on SO/PO
	SaleLineId            uint64      `json:"sale_line_id"             db:"sale_line_id"`             // Sale Line
	PurchaseLineId        uint64      `json:"purchase_line_id"         db:"purchase_line_id"`         // Purchase Order Line
	CreatedPurchaseLineId uint64      `json:"created_purchase_line_id" db:"created_purchase_line_id"` // Created Purchase Order Line

}

func (m *defaultCacheModel) Create(dto DtoCreate) (res int64, err error) {
	//cacheIdKey := m.FormatPrimary(dto.Id)
	insertMap := map[string]interface{}{}
	if dto.Sequence > 0 {
		insertMap[Columns.Sequence] = dto.Sequence
	}
	if dto.CompanyId > 0 {
		insertMap[Columns.CompanyId] = dto.CompanyId
	}
	if dto.ProductId > 0 {
		insertMap[Columns.ProductId] = dto.ProductId
	}
	if dto.ProductUom > 0 {
		insertMap[Columns.ProductUom] = dto.ProductUom
	}
	if dto.LocationId > 0 {
		insertMap[Columns.LocationId] = dto.LocationId
	}
	if dto.LocationDestId > 0 {
		insertMap[Columns.LocationDestId] = dto.LocationDestId
	}
	if dto.PartnerId > 0 {
		insertMap[Columns.PartnerId] = dto.PartnerId
	}
	if dto.PickingId > 0 {
		insertMap[Columns.PickingId] = dto.PickingId
	}
	if dto.GroupId > 0 {
		insertMap[Columns.GroupId] = dto.GroupId
	}
	if dto.RuleId > 0 {
		insertMap[Columns.RuleId] = dto.RuleId
	}
	if dto.PickingTypeId > 0 {
		insertMap[Columns.PickingTypeId] = dto.PickingTypeId
	}
	if dto.OriginReturnedMoveId > 0 {
		insertMap[Columns.OriginReturnedMoveId] = dto.OriginReturnedMoveId
	}
	if dto.RestrictPartnerId > 0 {
		insertMap[Columns.RestrictPartnerId] = dto.RestrictPartnerId
	}
	if dto.WarehouseId > 0 {
		insertMap[Columns.WarehouseId] = dto.WarehouseId
	}
	if dto.PackageLevelId > 0 {
		insertMap[Columns.PackageLevelId] = dto.PackageLevelId
	}
	if dto.NextSerialCount > 0 {
		insertMap[Columns.NextSerialCount] = dto.NextSerialCount
	}
	if dto.OrderpointId > 0 {
		insertMap[Columns.OrderpointId] = dto.OrderpointId
	}
	if dto.ProductPackagingId > 0 {
		insertMap[Columns.ProductPackagingId] = dto.ProductPackagingId
	}
	if dto.CreateUid > 0 {
		insertMap[Columns.CreateUid] = dto.CreateUid
	}
	if dto.WriteUid > 0 {
		insertMap[Columns.WriteUid] = dto.WriteUid
	}
	if dto.Name != "" {
		insertMap[Columns.Name] = dto.Name
	}
	if dto.Priority != "" {
		insertMap[Columns.Priority] = dto.Priority
	}
	if dto.State != "" {
		insertMap[Columns.State] = dto.State
	}
	if dto.Origin != "" {
		insertMap[Columns.Origin] = dto.Origin
	}
	if dto.ProcureMethod != "" {
		insertMap[Columns.ProcureMethod] = dto.ProcureMethod
	}
	if dto.Reference != "" {
		insertMap[Columns.Reference] = dto.Reference
	}
	if dto.NextSerial != "" {
		insertMap[Columns.NextSerial] = dto.NextSerial
	}
	if dto.DescriptionPicking != "" {
		insertMap[Columns.DescriptionPicking] = dto.DescriptionPicking
	}
	if dto.ProductQty > 0 {
		insertMap[Columns.ProductQty] = dto.ProductQty
	}
	if dto.ProductUomQty > 0 {
		insertMap[Columns.ProductUomQty] = dto.ProductUomQty
	}
	if dto.QuantityDone > 0 {
		insertMap[Columns.QuantityDone] = dto.QuantityDone
	}
	if dto.Scrapped {
		insertMap[Columns.Scrapped] = dto.Scrapped
	}
	if dto.PropagateCancel {
		insertMap[Columns.PropagateCancel] = dto.PropagateCancel
	}
	if dto.IsInventory {
		insertMap[Columns.IsInventory] = dto.IsInventory
	}
	if dto.Additional {
		insertMap[Columns.Additional] = dto.Additional
	}
	if dto.PriceUnit > 0 {
		insertMap[Columns.PriceUnit] = dto.PriceUnit
	}
	if dto.AnalyticAccountLineId > 0 {
		insertMap[Columns.AnalyticAccountLineId] = dto.AnalyticAccountLineId
	}
	if dto.ToRefund {
		insertMap[Columns.ToRefund] = dto.ToRefund
	}
	if dto.SaleLineId > 0 {
		insertMap[Columns.SaleLineId] = dto.SaleLineId
	}
	if dto.PurchaseLineId > 0 {
		insertMap[Columns.PurchaseLineId] = dto.PurchaseLineId
	}
	if dto.CreatedPurchaseLineId > 0 {
		insertMap[Columns.CreatedPurchaseLineId] = dto.CreatedPurchaseLineId
	}

	_, err = m.Insert(insertMap)
	return
}

func (m *defaultCacheModel) UpdateDto(dto Dto) (err error) {
	cacheIdKey := m.FormatPrimary(dto.Id)
	updateMap := map[string]interface{}{}
	if dto.Sequence > 0 {
		updateMap[Columns.Sequence] = dto.Sequence
	}
	if dto.CompanyId > 0 {
		updateMap[Columns.CompanyId] = dto.CompanyId
	}
	if dto.ProductId > 0 {
		updateMap[Columns.ProductId] = dto.ProductId
	}
	if dto.ProductUom > 0 {
		updateMap[Columns.ProductUom] = dto.ProductUom
	}
	if dto.LocationId > 0 {
		updateMap[Columns.LocationId] = dto.LocationId
	}
	if dto.LocationDestId > 0 {
		updateMap[Columns.LocationDestId] = dto.LocationDestId
	}
	if dto.PartnerId > 0 {
		updateMap[Columns.PartnerId] = dto.PartnerId
	}
	if dto.PickingId > 0 {
		updateMap[Columns.PickingId] = dto.PickingId
	}
	if dto.GroupId > 0 {
		updateMap[Columns.GroupId] = dto.GroupId
	}
	if dto.RuleId > 0 {
		updateMap[Columns.RuleId] = dto.RuleId
	}
	if dto.PickingTypeId > 0 {
		updateMap[Columns.PickingTypeId] = dto.PickingTypeId
	}
	if dto.OriginReturnedMoveId > 0 {
		updateMap[Columns.OriginReturnedMoveId] = dto.OriginReturnedMoveId
	}
	if dto.RestrictPartnerId > 0 {
		updateMap[Columns.RestrictPartnerId] = dto.RestrictPartnerId
	}
	if dto.WarehouseId > 0 {
		updateMap[Columns.WarehouseId] = dto.WarehouseId
	}
	if dto.PackageLevelId > 0 {
		updateMap[Columns.PackageLevelId] = dto.PackageLevelId
	}
	if dto.NextSerialCount > 0 {
		updateMap[Columns.NextSerialCount] = dto.NextSerialCount
	}
	if dto.OrderpointId > 0 {
		updateMap[Columns.OrderpointId] = dto.OrderpointId
	}
	if dto.ProductPackagingId > 0 {
		updateMap[Columns.ProductPackagingId] = dto.ProductPackagingId
	}
	if dto.CreateUid > 0 {
		updateMap[Columns.CreateUid] = dto.CreateUid
	}
	if dto.WriteUid > 0 {
		updateMap[Columns.WriteUid] = dto.WriteUid
	}
	if dto.Name != "" {
		updateMap[Columns.Name] = dto.Name
	}
	if dto.Priority != "" {
		updateMap[Columns.Priority] = dto.Priority
	}
	if dto.State != "" {
		updateMap[Columns.State] = dto.State
	}
	if dto.Origin != "" {
		updateMap[Columns.Origin] = dto.Origin
	}
	if dto.ProcureMethod != "" {
		updateMap[Columns.ProcureMethod] = dto.ProcureMethod
	}
	if dto.Reference != "" {
		updateMap[Columns.Reference] = dto.Reference
	}
	if dto.NextSerial != "" {
		updateMap[Columns.NextSerial] = dto.NextSerial
	}
	if dto.DescriptionPicking != "" {
		updateMap[Columns.DescriptionPicking] = dto.DescriptionPicking
	}
	if dto.ProductQty > 0 {
		updateMap[Columns.ProductQty] = dto.ProductQty
	}
	if dto.ProductUomQty > 0 {
		updateMap[Columns.ProductUomQty] = dto.ProductUomQty
	}
	if dto.QuantityDone > 0 {
		updateMap[Columns.QuantityDone] = dto.QuantityDone
	}
	if dto.Scrapped {
		updateMap[Columns.Scrapped] = dto.Scrapped
	}
	if dto.PropagateCancel {
		updateMap[Columns.PropagateCancel] = dto.PropagateCancel
	}
	if dto.IsInventory {
		updateMap[Columns.IsInventory] = dto.IsInventory
	}
	if dto.Additional {
		updateMap[Columns.Additional] = dto.Additional
	}
	if dto.PriceUnit > 0 {
		updateMap[Columns.PriceUnit] = dto.PriceUnit
	}
	if dto.AnalyticAccountLineId > 0 {
		updateMap[Columns.AnalyticAccountLineId] = dto.AnalyticAccountLineId
	}
	if dto.ToRefund {
		updateMap[Columns.ToRefund] = dto.ToRefund
	}
	if dto.SaleLineId > 0 {
		updateMap[Columns.SaleLineId] = dto.SaleLineId
	}
	if dto.PurchaseLineId > 0 {
		updateMap[Columns.PurchaseLineId] = dto.PurchaseLineId
	}
	if dto.CreatedPurchaseLineId > 0 {
		updateMap[Columns.CreatedPurchaseLineId] = dto.CreatedPurchaseLineId
	}

	err = m.Update(m.NewUpdateBuilder().Where(squirrel.Eq{
		Columns.Id: dto.Id,
	}), updateMap, cacheIdKey)
	return
}

type DtoUpdateAutomatic struct {
	Id                    uint64      `json:"id"                       db:"id"`                       //
	Sequence              int64       `json:"sequence"                 db:"sequence"`                 // Sequence
	CompanyId             uint64      `json:"company_id"               db:"company_id"`               // Company
	ProductId             uint64      `json:"product_id"               db:"product_id"`               // Product
	ProductUom            int64       `json:"product_uom"              db:"product_uom"`              // UoM
	LocationId            uint64      `json:"location_id"              db:"location_id"`              // Source Location
	LocationDestId        uint64      `json:"location_dest_id"         db:"location_dest_id"`         // Destination Location
	PartnerId             uint64      `json:"partner_id"               db:"partner_id"`               // Destination Address
	PickingId             uint64      `json:"picking_id"               db:"picking_id"`               // Transfer
	GroupId               uint64      `json:"group_id"                 db:"group_id"`                 // Procurement Group
	RuleId                uint64      `json:"rule_id"                  db:"rule_id"`                  // Stock Rule
	PickingTypeId         uint64      `json:"picking_type_id"          db:"picking_type_id"`          // Operation Type
	OriginReturnedMoveId  uint64      `json:"origin_returned_move_id"  db:"origin_returned_move_id"`  // Origin return move
	RestrictPartnerId     uint64      `json:"restrict_partner_id"      db:"restrict_partner_id"`      // Owner
	WarehouseId           uint64      `json:"warehouse_id"             db:"warehouse_id"`             // Warehouse
	PackageLevelId        uint64      `json:"package_level_id"         db:"package_level_id"`         // Package Level
	NextSerialCount       int64       `json:"next_serial_count"        db:"next_serial_count"`        // Number of SN
	OrderpointId          uint64      `json:"orderpoint_id"            db:"orderpoint_id"`            // Original Reordering Rule
	ProductPackagingId    uint64      `json:"product_packaging_id"     db:"product_packaging_id"`     // Packaging
	CreateUid             uint64      `json:"create_uid"               db:"create_uid"`               // Created by
	WriteUid              uint64      `json:"write_uid"                db:"write_uid"`                // Last Updated by
	Name                  string      `json:"name"                     db:"name"`                     // Description
	Priority              string      `json:"priority"                 db:"priority"`                 // Priority
	State                 string      `json:"state"                    db:"state"`                    // Status
	Origin                string      `json:"origin"                   db:"origin"`                   // Source Document
	ProcureMethod         string      `json:"procure_method"           db:"procure_method"`           // Supply Method
	Reference             string      `json:"reference"                db:"reference"`                // Reference
	NextSerial            string      `json:"next_serial"              db:"next_serial"`              // First SN
	ReservationDate       *gtime.Time `json:"reservation_date"         db:"reservation_date"`         // Date to Reserve
	DescriptionPicking    string      `json:"description_picking"      db:"description_picking"`      // Description of Picking
	ProductQty            float64     `json:"product_qty"              db:"product_qty"`              // Real Quantity
	ProductUomQty         float64     `json:"product_uom_qty"          db:"product_uom_qty"`          // Demand
	QuantityDone          float64     `json:"quantity_done"            db:"quantity_done"`            // Quantity Done
	Scrapped              bool        `json:"scrapped"                 db:"scrapped"`                 // Scrapped
	PropagateCancel       bool        `json:"propagate_cancel"         db:"propagate_cancel"`         // Propagate cancel and split
	IsInventory           bool        `json:"is_inventory"             db:"is_inventory"`             // Inventory
	Additional            bool        `json:"additional"               db:"additional"`               // Whether the move was added after the picking's confirmation
	Date                  *gtime.Time `json:"date"                     db:"date"`                     // Date Scheduled
	DateDeadline          *gtime.Time `json:"date_deadline"            db:"date_deadline"`            // Deadline
	DelayAlertDate        *gtime.Time `json:"delay_alert_date"         db:"delay_alert_date"`         // Delay Alert Date
	CreateDate            *gtime.Time `json:"create_date"              db:"create_date"`              // Created on
	WriteDate             *gtime.Time `json:"write_date"               db:"write_date"`               // Last Updated on
	PriceUnit             float64     `json:"price_unit"               db:"price_unit"`               // Unit Price
	AnalyticAccountLineId uint64      `json:"analytic_account_line_id" db:"analytic_account_line_id"` // Analytic Account Line
	ToRefund              bool        `json:"to_refund"                db:"to_refund"`                // Update quantities on SO/PO
	SaleLineId            uint64      `json:"sale_line_id"             db:"sale_line_id"`             // Sale Line
	PurchaseLineId        uint64      `json:"purchase_line_id"         db:"purchase_line_id"`         // Purchase Order Line
	CreatedPurchaseLineId uint64      `json:"created_purchase_line_id" db:"created_purchase_line_id"` // Created Purchase Order Line

}

func (m *defaultCacheModel) UpdateAutomaticById(id uint64, dto DtoUpdateAutomatic, fields []string) error {
	key := m.FormatPrimary(id)
	return m.updateAutomatic(m.NewUpdateBuilder().Where(squirrel.Eq{
		Columns.Id: id,
	}), dto, fields, key)
}

func (m *defaultCacheModel) updateAutomatic(builder squirrel.UpdateBuilder, dto DtoUpdateAutomatic, fields []string, cacheKey ...string) error {
	automatic := xsql.UpdateByFields(dto, fields)
	return m.Update(builder, automatic, cacheKey...)
}
