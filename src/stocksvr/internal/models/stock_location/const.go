package stock_location

import (
	"database/sql"
	"github.com/gogf/gf/os/gtime"
	"github.com/i-Things/things/shared/xzero/xmodel"
)

const Table = `public.stock_location`

var (
	EntityFieldNames      = xmodel.RawFieldNames(&Entity{}, true)
	CacheEntityFieldNames = xmodel.RawFieldNames(&CacheEntity{}, true)
	IdNameField           = []string{"id", "name"}
	cacheEntityIdPrefix   = "cache:dbName:stock_location:id:"
)

type (
	Entities []*Entity
	Entity   struct {
		Id                       uint64         `json:"id"                         db:"id"`                         //
		LocationId               sql.NullInt64  `json:"location_id"                db:"location_id"`                // Parent Location
		Posx                     sql.NullInt64  `json:"posx"                       db:"posx"`                       // Corridor (X)
		Posy                     sql.NullInt64  `json:"posy"                       db:"posy"`                       // Shelves (Y)
		Posz                     sql.NullInt64  `json:"posz"                       db:"posz"`                       // Height (Z)
		CompanyId                sql.NullInt64  `json:"company_id"                 db:"company_id"`                 // Company
		RemovalStrategyId        sql.NullInt64  `json:"removal_strategy_id"        db:"removal_strategy_id"`        // Removal Strategy
		CyclicInventoryFrequency sql.NullInt64  `json:"cyclic_inventory_frequency" db:"cyclic_inventory_frequency"` // Inventory Frequency (Days)
		WarehouseId              sql.NullInt64  `json:"warehouse_id"               db:"warehouse_id"`               // Warehouse
		StorageCategoryId        sql.NullInt64  `json:"storage_category_id"        db:"storage_category_id"`        // Storage Category
		CreateUid                sql.NullInt64  `json:"create_uid"                 db:"create_uid"`                 // Created by
		WriteUid                 sql.NullInt64  `json:"write_uid"                  db:"write_uid"`                  // Last Updated by
		Name                     string         `json:"name"                       db:"name"`                       // Location Name
		CompleteName             sql.NullString `json:"complete_name"              db:"complete_name"`              // Full Location Name
		Usage                    string         `json:"usage"                      db:"usage"`                      // Location Type
		ParentPath               sql.NullString `json:"parent_path"                db:"parent_path"`                // Parent Path
		Barcode                  sql.NullString `json:"barcode"                    db:"barcode"`                    // Barcode
		LastInventoryDate        *gtime.Time    `json:"last_inventory_date"        db:"last_inventory_date"`        // Last Effective Inventory
		NextInventoryDate        *gtime.Time    `json:"next_inventory_date"        db:"next_inventory_date"`        // Next Expected Inventory
		Comment                  sql.NullString `json:"comment"                    db:"comment"`                    // Additional Information
		Active                   sql.NullBool   `json:"active"                     db:"active"`                     // Active
		ScrapLocation            sql.NullBool   `json:"scrap_location"             db:"scrap_location"`             // Is a Scrap Location?
		ReturnLocation           sql.NullBool   `json:"return_location"            db:"return_location"`            // Is a Return Location?
		ReplenishLocation        sql.NullBool   `json:"replenish_location"         db:"replenish_location"`         // Replenish Location
		CreateDate               *gtime.Time    `json:"create_date"                db:"create_date"`                // Created on
		WriteDate                *gtime.Time    `json:"write_date"                 db:"write_date"`                 // Last Updated on
		ValuationInAccountId     sql.NullInt64  `json:"valuation_in_account_id"    db:"valuation_in_account_id"`    // Stock Valuation Account (Incoming)
		ValuationOutAccountId    sql.NullInt64  `json:"valuation_out_account_id"   db:"valuation_out_account_id"`   // Stock Valuation Account (Outgoing)

	}

	CacheEntities []*CacheEntity
	CacheEntity   struct {
	}
	Status uint32
)

var (
	Columns = struct {
		Id                       string //
		LocationId               string // Parent Location
		Posx                     string // Corridor (X)
		Posy                     string // Shelves (Y)
		Posz                     string // Height (Z)
		CompanyId                string // Company
		RemovalStrategyId        string // Removal Strategy
		CyclicInventoryFrequency string // Inventory Frequency (Days)
		WarehouseId              string // Warehouse
		StorageCategoryId        string // Storage Category
		CreateUid                string // Created by
		WriteUid                 string // Last Updated by
		Name                     string // Location Name
		CompleteName             string // Full Location Name
		Usage                    string // Location Type
		ParentPath               string // Parent Path
		Barcode                  string // Barcode
		LastInventoryDate        string // Last Effective Inventory
		NextInventoryDate        string // Next Expected Inventory
		Comment                  string // Additional Information
		Active                   string // Active
		ScrapLocation            string // Is a Scrap Location?
		ReturnLocation           string // Is a Return Location?
		ReplenishLocation        string // Replenish Location
		CreateDate               string // Created on
		WriteDate                string // Last Updated on
		ValuationInAccountId     string // Stock Valuation Account (Incoming)
		ValuationOutAccountId    string // Stock Valuation Account (Outgoing)

	}{
		Id:                       "id",
		LocationId:               "location_id",
		Posx:                     "posx",
		Posy:                     "posy",
		Posz:                     "posz",
		CompanyId:                "company_id",
		RemovalStrategyId:        "removal_strategy_id",
		CyclicInventoryFrequency: "cyclic_inventory_frequency",
		WarehouseId:              "warehouse_id",
		StorageCategoryId:        "storage_category_id",
		CreateUid:                "create_uid",
		WriteUid:                 "write_uid",
		Name:                     "name",
		CompleteName:             "complete_name",
		Usage:                    "usage",
		ParentPath:               "parent_path",
		Barcode:                  "barcode",
		LastInventoryDate:        "last_inventory_date",
		NextInventoryDate:        "next_inventory_date",
		Comment:                  "comment",
		Active:                   "active",
		ScrapLocation:            "scrap_location",
		ReturnLocation:           "return_location",
		ReplenishLocation:        "replenish_location",
		CreateDate:               "create_date",
		WriteDate:                "write_date",
		ValuationInAccountId:     "valuation_in_account_id",
		ValuationOutAccountId:    "valuation_out_account_id",
	}
)
