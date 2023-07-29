package stock_move

import (
	"github.com/gogf/gf/os/gtime"
	"github.com/i-Things/things/shared/xzero/xmodel"
)

const Table = "`stock_move`"

var (
	EntityFieldNames      = xmodel.RawFieldNames(&Entity{})
	CacheEntityFieldNames = xmodel.RawFieldNames(&CacheEntity{})
	IdNameField           = []string{"id", "name"}
	cacheEntityIdPrefix   = "cache:dbName:stock_move:id:"
)

type (
	Entities []*Entity
	Entity   struct {
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

	CacheEntities []*CacheEntity
	CacheEntity   struct {
	}
	Status uint32
)

var (
	Columns = struct {
		Id                    string //
		Sequence              string // Sequence
		CompanyId             string // Company
		ProductId             string // Product
		ProductUom            string // UoM
		LocationId            string // Source Location
		LocationDestId        string // Destination Location
		PartnerId             string // Destination Address
		PickingId             string // Transfer
		GroupId               string // Procurement Group
		RuleId                string // Stock Rule
		PickingTypeId         string // Operation Type
		OriginReturnedMoveId  string // Origin return move
		RestrictPartnerId     string // Owner
		WarehouseId           string // Warehouse
		PackageLevelId        string // Package Level
		NextSerialCount       string // Number of SN
		OrderpointId          string // Original Reordering Rule
		ProductPackagingId    string // Packaging
		CreateUid             string // Created by
		WriteUid              string // Last Updated by
		Name                  string // Description
		Priority              string // Priority
		State                 string // Status
		Origin                string // Source Document
		ProcureMethod         string // Supply Method
		Reference             string // Reference
		NextSerial            string // First SN
		ReservationDate       string // Date to Reserve
		DescriptionPicking    string // Description of Picking
		ProductQty            string // Real Quantity
		ProductUomQty         string // Demand
		QuantityDone          string // Quantity Done
		Scrapped              string // Scrapped
		PropagateCancel       string // Propagate cancel and split
		IsInventory           string // Inventory
		Additional            string // Whether the move was added after the picking's confirmation
		Date                  string // Date Scheduled
		DateDeadline          string // Deadline
		DelayAlertDate        string // Delay Alert Date
		CreateDate            string // Created on
		WriteDate             string // Last Updated on
		PriceUnit             string // Unit Price
		AnalyticAccountLineId string // Analytic Account Line
		ToRefund              string // Update quantities on SO/PO
		SaleLineId            string // Sale Line
		PurchaseLineId        string // Purchase Order Line
		CreatedPurchaseLineId string // Created Purchase Order Line

	}{
		Id:                    "id",
		Sequence:              "sequence",
		CompanyId:             "company_id",
		ProductId:             "product_id",
		ProductUom:            "product_uom",
		LocationId:            "location_id",
		LocationDestId:        "location_dest_id",
		PartnerId:             "partner_id",
		PickingId:             "picking_id",
		GroupId:               "group_id",
		RuleId:                "rule_id",
		PickingTypeId:         "picking_type_id",
		OriginReturnedMoveId:  "origin_returned_move_id",
		RestrictPartnerId:     "restrict_partner_id",
		WarehouseId:           "warehouse_id",
		PackageLevelId:        "package_level_id",
		NextSerialCount:       "next_serial_count",
		OrderpointId:          "orderpoint_id",
		ProductPackagingId:    "product_packaging_id",
		CreateUid:             "create_uid",
		WriteUid:              "write_uid",
		Name:                  "name",
		Priority:              "priority",
		State:                 "state",
		Origin:                "origin",
		ProcureMethod:         "procure_method",
		Reference:             "reference",
		NextSerial:            "next_serial",
		ReservationDate:       "reservation_date",
		DescriptionPicking:    "description_picking",
		ProductQty:            "product_qty",
		ProductUomQty:         "product_uom_qty",
		QuantityDone:          "quantity_done",
		Scrapped:              "scrapped",
		PropagateCancel:       "propagate_cancel",
		IsInventory:           "is_inventory",
		Additional:            "additional",
		Date:                  "date",
		DateDeadline:          "date_deadline",
		DelayAlertDate:        "delay_alert_date",
		CreateDate:            "create_date",
		WriteDate:             "write_date",
		PriceUnit:             "price_unit",
		AnalyticAccountLineId: "analytic_account_line_id",
		ToRefund:              "to_refund",
		SaleLineId:            "sale_line_id",
		PurchaseLineId:        "purchase_line_id",
		CreatedPurchaseLineId: "created_purchase_line_id",
	}
)
