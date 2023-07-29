package stock_location

import (
	"github.com/Masterminds/squirrel"
	"github.com/gogf/gf/os/gtime"
	"github.com/i-Things/things/shared/xzero/xsql"
)

type Dto struct {
	Id                       uint64      `json:"id"                         db:"id"`                         //
	LocationId               uint64      `json:"location_id"                db:"location_id"`                // Parent Location
	Posx                     int64       `json:"posx"                       db:"posx"`                       // Corridor (X)
	Posy                     int64       `json:"posy"                       db:"posy"`                       // Shelves (Y)
	Posz                     int64       `json:"posz"                       db:"posz"`                       // Height (Z)
	CompanyId                uint64      `json:"company_id"                 db:"company_id"`                 // Company
	RemovalStrategyId        uint64      `json:"removal_strategy_id"        db:"removal_strategy_id"`        // Removal Strategy
	CyclicInventoryFrequency int64       `json:"cyclic_inventory_frequency" db:"cyclic_inventory_frequency"` // Inventory Frequency (Days)
	WarehouseId              uint64      `json:"warehouse_id"               db:"warehouse_id"`               // Warehouse
	StorageCategoryId        uint64      `json:"storage_category_id"        db:"storage_category_id"`        // Storage Category
	CreateUid                uint64      `json:"create_uid"                 db:"create_uid"`                 // Created by
	WriteUid                 uint64      `json:"write_uid"                  db:"write_uid"`                  // Last Updated by
	Name                     string      `json:"name"                       db:"name"`                       // Location Name
	CompleteName             string      `json:"complete_name"              db:"complete_name"`              // Full Location Name
	Usage                    string      `json:"usage"                      db:"usage"`                      // Location Type
	ParentPath               string      `json:"parent_path"                db:"parent_path"`                // Parent Path
	Barcode                  string      `json:"barcode"                    db:"barcode"`                    // Barcode
	LastInventoryDate        *gtime.Time `json:"last_inventory_date"        db:"last_inventory_date"`        // Last Effective Inventory
	NextInventoryDate        *gtime.Time `json:"next_inventory_date"        db:"next_inventory_date"`        // Next Expected Inventory
	Comment                  string      `json:"comment"                    db:"comment"`                    // Additional Information
	Active                   bool        `json:"active"                     db:"active"`                     // Active
	ScrapLocation            bool        `json:"scrap_location"             db:"scrap_location"`             // Is a Scrap Location?
	ReturnLocation           bool        `json:"return_location"            db:"return_location"`            // Is a Return Location?
	ReplenishLocation        bool        `json:"replenish_location"         db:"replenish_location"`         // Replenish Location
	CreateDate               *gtime.Time `json:"create_date"                db:"create_date"`                // Created on
	WriteDate                *gtime.Time `json:"write_date"                 db:"write_date"`                 // Last Updated on
	ValuationInAccountId     uint64      `json:"valuation_in_account_id"    db:"valuation_in_account_id"`    // Stock Valuation Account (Incoming)
	ValuationOutAccountId    uint64      `json:"valuation_out_account_id"   db:"valuation_out_account_id"`   // Stock Valuation Account (Outgoing)
}

type DtoCreate struct {
	LocationId               uint64      `json:"location_id"                db:"location_id"`                // Parent Location
	Posx                     int64       `json:"posx"                       db:"posx"`                       // Corridor (X)
	Posy                     int64       `json:"posy"                       db:"posy"`                       // Shelves (Y)
	Posz                     int64       `json:"posz"                       db:"posz"`                       // Height (Z)
	CompanyId                uint64      `json:"company_id"                 db:"company_id"`                 // Company
	RemovalStrategyId        uint64      `json:"removal_strategy_id"        db:"removal_strategy_id"`        // Removal Strategy
	CyclicInventoryFrequency int64       `json:"cyclic_inventory_frequency" db:"cyclic_inventory_frequency"` // Inventory Frequency (Days)
	WarehouseId              uint64      `json:"warehouse_id"               db:"warehouse_id"`               // Warehouse
	StorageCategoryId        uint64      `json:"storage_category_id"        db:"storage_category_id"`        // Storage Category
	CreateUid                uint64      `json:"create_uid"                 db:"create_uid"`                 // Created by
	WriteUid                 uint64      `json:"write_uid"                  db:"write_uid"`                  // Last Updated by
	Name                     string      `json:"name"                       db:"name"`                       // Location Name
	CompleteName             string      `json:"complete_name"              db:"complete_name"`              // Full Location Name
	Usage                    string      `json:"usage"                      db:"usage"`                      // Location Type
	ParentPath               string      `json:"parent_path"                db:"parent_path"`                // Parent Path
	Barcode                  string      `json:"barcode"                    db:"barcode"`                    // Barcode
	LastInventoryDate        *gtime.Time `json:"last_inventory_date"        db:"last_inventory_date"`        // Last Effective Inventory
	NextInventoryDate        *gtime.Time `json:"next_inventory_date"        db:"next_inventory_date"`        // Next Expected Inventory
	Comment                  string      `json:"comment"                    db:"comment"`                    // Additional Information
	Active                   bool        `json:"active"                     db:"active"`                     // Active
	ScrapLocation            bool        `json:"scrap_location"             db:"scrap_location"`             // Is a Scrap Location?
	ReturnLocation           bool        `json:"return_location"            db:"return_location"`            // Is a Return Location?
	ReplenishLocation        bool        `json:"replenish_location"         db:"replenish_location"`         // Replenish Location
	CreateDate               *gtime.Time `json:"create_date"                db:"create_date"`                // Created on
	WriteDate                *gtime.Time `json:"write_date"                 db:"write_date"`                 // Last Updated on
	ValuationInAccountId     uint64      `json:"valuation_in_account_id"    db:"valuation_in_account_id"`    // Stock Valuation Account (Incoming)
	ValuationOutAccountId    uint64      `json:"valuation_out_account_id"   db:"valuation_out_account_id"`   // Stock Valuation Account (Outgoing)
}

func (m *defaultCacheModel) Create(dto DtoCreate) (res int64, err error) {
	insertMap := map[string]interface{}{
		Columns.Name:  dto.Name,  // Location Name
		Columns.Usage: dto.Usage, // Location Type
	}
	if dto.Posx != 0 {
		insertMap[Columns.Posx] = dto.Posx
	}
	if dto.Posy != 0 {
		insertMap[Columns.Posy] = dto.Posy
	}
	if dto.Posz != 0 {
		insertMap[Columns.Posz] = dto.Posz
	}
	if dto.Posx != 0 {
		insertMap[Columns.Posx] = dto.Posx
	}
	if dto.CompanyId != 0 {
		insertMap[Columns.CompanyId] = dto.CompanyId
	}
	if dto.RemovalStrategyId != 0 {
		insertMap[Columns.RemovalStrategyId] = dto.RemovalStrategyId
	}
	if dto.CyclicInventoryFrequency != 0 {
		insertMap[Columns.CyclicInventoryFrequency] = dto.CyclicInventoryFrequency
	}
	if dto.WarehouseId != 0 {
		insertMap[Columns.WarehouseId] = dto.WarehouseId
	}
	if dto.StorageCategoryId != 0 {
		insertMap[Columns.StorageCategoryId] = dto.StorageCategoryId
	}
	if dto.CreateUid != 0 {
		insertMap[Columns.CreateUid] = dto.CreateUid
	}
	if dto.WriteUid != 0 {
		insertMap[Columns.WriteUid] = dto.WriteUid
	}

	if dto.CompleteName != "" {
		insertMap[Columns.CompleteName] = dto.CompleteName
	}
	if dto.ParentPath != "" {
		insertMap[Columns.ParentPath] = dto.ParentPath
	}
	if dto.Barcode != "" {
		insertMap[Columns.Barcode] = dto.Barcode
	}
	if dto.Comment != "" {
		insertMap[Columns.Comment] = dto.Comment
	}
	//if dto.LastInventoryDate != nil {
	//	insertMap[Columns.LastInventoryDate] = dto.LastInventoryDate
	//}
	//if dto.NextInventoryDate != nil {
	//	insertMap[Columns.NextInventoryDate] = dto.NextInventoryDate
	//}
	//
	//if dto.CreateDate != nil {
	//	insertMap[Columns.CreateDate] = dto.CreateDate
	//}
	//if dto.WriteDate != nil {
	//	insertMap[Columns.WriteDate] = dto.WriteDate
	//}
	if dto.ValuationInAccountId != 0 {
		insertMap[Columns.ValuationInAccountId] = dto.ValuationInAccountId
	}
	if dto.ValuationOutAccountId != 0 {
		insertMap[Columns.ValuationOutAccountId] = dto.ValuationOutAccountId
	}

	_, err = m.Insert(insertMap)
	if err != nil {
		return
	}
	//result.LastInsertId()
	return
}

func (m *defaultCacheModel) UpdateDto(dto Dto) (err error) {
	updateMap := map[string]interface{}{
		Columns.Name:  dto.Name,  // Location Name
		Columns.Usage: dto.Usage, // Location Type
	}
	if dto.Posx != 0 {
		updateMap[Columns.Posx] = dto.Posx
	}
	if dto.Posy != 0 {
		updateMap[Columns.Posy] = dto.Posy
	}
	if dto.Posz != 0 {
		updateMap[Columns.Posz] = dto.Posz
	}
	if dto.Posx != 0 {
		updateMap[Columns.Posx] = dto.Posx
	}
	if dto.CompanyId != 0 {
		updateMap[Columns.CompanyId] = dto.CompanyId
	}
	if dto.RemovalStrategyId != 0 {
		updateMap[Columns.RemovalStrategyId] = dto.RemovalStrategyId
	}
	if dto.CyclicInventoryFrequency != 0 {
		updateMap[Columns.CyclicInventoryFrequency] = dto.CyclicInventoryFrequency
	}
	if dto.WarehouseId != 0 {
		updateMap[Columns.WarehouseId] = dto.WarehouseId
	}
	if dto.StorageCategoryId != 0 {
		updateMap[Columns.StorageCategoryId] = dto.StorageCategoryId
	}
	if dto.CreateUid != 0 {
		updateMap[Columns.CreateUid] = dto.CreateUid
	}
	if dto.WriteUid != 0 {
		updateMap[Columns.WriteUid] = dto.WriteUid
	}

	if dto.CompleteName != "" {
		updateMap[Columns.CompleteName] = dto.CompleteName
	}
	if dto.ParentPath != "" {
		updateMap[Columns.ParentPath] = dto.ParentPath
	}
	if dto.Barcode != "" {
		updateMap[Columns.Barcode] = dto.Barcode
	}
	if dto.Comment != "" {
		updateMap[Columns.Comment] = dto.Comment
	}
	//if dto.LastInventoryDate != nil {
	//	updateMap[Columns.LastInventoryDate] = dto.LastInventoryDate
	//}
	//if dto.NextInventoryDate != nil {
	//	updateMap[Columns.NextInventoryDate] = dto.NextInventoryDate
	//}
	//
	//if dto.CreateDate != nil {
	//	updateMap[Columns.CreateDate] = dto.CreateDate
	//}
	//if dto.WriteDate != nil {
	//	updateMap[Columns.WriteDate] = dto.WriteDate
	//}
	if dto.ValuationInAccountId != 0 {
		updateMap[Columns.ValuationInAccountId] = dto.ValuationInAccountId
	}
	if dto.ValuationOutAccountId != 0 {
		updateMap[Columns.ValuationOutAccountId] = dto.ValuationOutAccountId
	}

	cacheIdKey := m.FormatPrimary(dto.Id)
	err = m.Update(m.NewUpdateBuilder().Where(squirrel.Eq{
		Columns.Id: dto.Id,
	}), updateMap, cacheIdKey)
	return
}

type DtoUpdateAutomatic struct {
	Id                       uint64      `json:"id"                         db:"id"`                         //
	LocationId               uint64      `json:"location_id"                db:"location_id"`                // Parent Location
	Posx                     int64       `json:"posx"                       db:"posx"`                       // Corridor (X)
	Posy                     int64       `json:"posy"                       db:"posy"`                       // Shelves (Y)
	Posz                     int64       `json:"posz"                       db:"posz"`                       // Height (Z)
	CompanyId                uint64      `json:"company_id"                 db:"company_id"`                 // Company
	RemovalStrategyId        uint64      `json:"removal_strategy_id"        db:"removal_strategy_id"`        // Removal Strategy
	CyclicInventoryFrequency int64       `json:"cyclic_inventory_frequency" db:"cyclic_inventory_frequency"` // Inventory Frequency (Days)
	WarehouseId              uint64      `json:"warehouse_id"               db:"warehouse_id"`               // Warehouse
	StorageCategoryId        uint64      `json:"storage_category_id"        db:"storage_category_id"`        // Storage Category
	CreateUid                uint64      `json:"create_uid"                 db:"create_uid"`                 // Created by
	WriteUid                 uint64      `json:"write_uid"                  db:"write_uid"`                  // Last Updated by
	Name                     string      `json:"name"                       db:"name"`                       // Location Name
	CompleteName             string      `json:"complete_name"              db:"complete_name"`              // Full Location Name
	Usage                    string      `json:"usage"                      db:"usage"`                      // Location Type
	ParentPath               string      `json:"parent_path"                db:"parent_path"`                // Parent Path
	Barcode                  string      `json:"barcode"                    db:"barcode"`                    // Barcode
	LastInventoryDate        *gtime.Time `json:"last_inventory_date"        db:"last_inventory_date"`        // Last Effective Inventory
	NextInventoryDate        *gtime.Time `json:"next_inventory_date"        db:"next_inventory_date"`        // Next Expected Inventory
	Comment                  string      `json:"comment"                    db:"comment"`                    // Additional Information
	Active                   bool        `json:"active"                     db:"active"`                     // Active
	ScrapLocation            bool        `json:"scrap_location"             db:"scrap_location"`             // Is a Scrap Location?
	ReturnLocation           bool        `json:"return_location"            db:"return_location"`            // Is a Return Location?
	ReplenishLocation        bool        `json:"replenish_location"         db:"replenish_location"`         // Replenish Location
	CreateDate               *gtime.Time `json:"create_date"                db:"create_date"`                // Created on
	WriteDate                *gtime.Time `json:"write_date"                 db:"write_date"`                 // Last Updated on
	ValuationInAccountId     uint64      `json:"valuation_in_account_id"    db:"valuation_in_account_id"`    // Stock Valuation Account (Incoming)
	ValuationOutAccountId    uint64      `json:"valuation_out_account_id"   db:"valuation_out_account_id"`   // Stock Valuation Account (Outgoing)

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
