package relationDB

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"gorm.io/gorm"
)

type ProductInfoRepo struct {
	db *gorm.DB
}

type ProductFilter struct {
	DeviceType    int64
	DeviceTypes   []int64
	ProductName   string
	ProductIDs    []string
	ProductNames  []string
	Tags          map[string]string
	ProtocolConf  map[string]string
	WithProtocol  bool
	WithCategory  bool
	ProtocolCode  string
	ProtocolType  string
	ProtocolTrans string
	CategoryIDs   []int64
	SceneMode     string
	SceneModes    []string
	Status        devices.ProductStatus
	Statuses      []devices.ProductStatus
	AreaID        int64
	AreaIDPath    string
	NetType       int64
}

func NewProductInfoRepo(in any) *ProductInfoRepo {
	return &ProductInfoRepo{db: stores.GetCommonConn(in)}
}

func (p ProductInfoRepo) fmtFilter(ctx context.Context, f ProductFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.DeviceType != 0 {
		db = db.Where("device_type=?", f.DeviceType)
	}
	if len(f.DeviceTypes) != 0 {
		db = db.Where("device_type in ?", f.DeviceTypes)
	}
	if f.AreaID != 0 {
		subQuery := p.db.Model(&DmDeviceInfo{}).Distinct("product_id").Select("product_id").Where("area_id=?", f.AreaID)
		db = db.Where("product_id in (?)", subQuery)
	}
	if f.AreaIDPath != "" {
		subQuery := p.db.Model(&DmDeviceInfo{}).Distinct("product_id").Select("product_id").Where("area_id_path like ?", f.AreaIDPath+"%")
		db = db.Where("product_id in (?)", subQuery)
	}
	if len(f.CategoryIDs) != 0 {
		db = db.Where("category_id in ?", f.CategoryIDs)
	}

	if f.ProtocolCode != "" {
		db = db.Where("protocol_code=? or sub_protocol_code=?", f.ProtocolCode, f.ProtocolCode)
	}
	if f.ProtocolType != "" || f.ProtocolTrans != "" {
		subQuery := p.db.Model(&DmProtocolInfo{}).Select("code")
		if f.ProtocolType != "" {
			subQuery = subQuery.Where(
				fmt.Sprintf("%s = ?", stores.Col("type")), f.ProtocolType)
		}
		if f.ProtocolTrans != "" {
			subQuery = subQuery.Where(
				fmt.Sprintf("%s = ?", stores.Col("trans_protocol")), f.ProtocolTrans)
		}
		db = db.Where("protocol_code in (?) or sub_protocol_code in (?)", subQuery, subQuery)
	}
	if f.SceneMode != "" {
		db = db.Where("scene_mode=?", f.SceneMode)
	}
	if len(f.SceneModes) != 0 {
		db = db.Where("scene_mode in ?", f.SceneModes)
	}
	if f.Status != 0 {
		db = db.Where("status=?", f.Status)
	}
	if len(f.Statuses) != 0 {
		db = db.Where("status in ?", f.Statuses)
	}
	if f.WithProtocol {
		db = db.Preload("Protocol").Preload("SubProtocol")
	}
	if f.WithCategory {
		db = db.Preload("Category")
	}
	if f.ProductName != "" {
		db = db.Where("product_name like ?", "%"+f.ProductName+"%")
	}
	if len(f.ProductIDs) != 0 {
		db = db.Where("product_id in ?", f.ProductIDs)
	}
	if len(f.ProductNames) != 0 {
		db = db.Where("product_name in ?", f.ProductNames)
	}
	if f.ProtocolConf != nil {
		for k, v := range f.ProtocolConf {
			db = stores.CmpJsonObjEq(k, v).Where(db, "protocol_conf")
		}
	}
	if f.Tags != nil {
		for k, v := range f.Tags {
			db = stores.CmpJsonObjEq(k, v).Where(db, "tags")
		}
	}
	if f.NetType != 0 {
		db = db.Where("net_type = ?", f.NetType)
	}
	return db
}

func (p ProductInfoRepo) Insert(ctx context.Context, data *DmProductInfo) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		err := p.db.WithContext(ctx).Create(data).Error
		if err != nil {
			return err
		}

		if data.CategoryID != def.NotClassified {
			var ProductCategoryIDs = []int64{data.CategoryID}
			if data.CategoryID > def.NotClassified {
				pcDB := NewProductCategoryRepo(tx)
				pcs, err := pcDB.FindOne(ctx, data.CategoryID)
				if err != nil {
					return err
				}
				ProductCategoryIDs = append(ProductCategoryIDs, utils.GetIDPath(pcs.IDPath)...)
			}
			pcs, err := NewProductCategorySchemaRepo(tx).FindByFilter(ctx, ProductCategorySchemaFilter{ProductCategoryIDs: ProductCategoryIDs}, nil)
			if err != nil {
				return err
			}
			if len(pcs) == 0 {
				return nil
			}
			identifiers := utils.ToSliceWithFunc(pcs, func(in *DmProductCategorySchema) string {
				return in.Identifier
			})
			ids, err := NewCommonSchemaRepo(tx).FindByFilter(ctx, CommonSchemaFilter{Identifiers: identifiers}, nil)
			if err != nil {
				return err
			}
			if len(ids) == 0 {
				return nil
			}
			var schemas []*DmSchemaInfo
			var idents []string
			for _, v := range ids {
				v.ID = 0
				v.Tag = schema.TagRequired
				idents = append(idents, v.Identifier)
				schemas = append(schemas, &DmSchemaInfo{
					ProductID:    data.ProductID,
					DmSchemaCore: v.DmSchemaCore,
				})
			}
			//如果定义了产品级的,需要删除设备级的
			err = NewSchemaInfoRepo(tx).DeleteByFilter(ctx, SchemaInfoFilter{ProductID: data.ProductID, Tag: schema.TagDevice, Identifiers: idents})
			if err != nil {
				return err
			}
			err = NewProductSchemaRepo(tx).MultiInsert(ctx, schemas)
			if err != nil {
				return err
			}
			err = NewProductConfigRepo(tx).Insert(ctx, &DmProductConfig{ProductID: data.ProductID})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return stores.ErrFmt(err)
}

func (p ProductInfoRepo) FindOneByFilter(ctx context.Context, f ProductFilter) (*DmProductInfo, error) {
	var result DmProductInfo
	db := p.fmtFilter(ctx, f)
	err := db.Preload("Config").First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p ProductInfoRepo) Update(ctx context.Context, data *DmProductInfo) error {
	data2 := *data
	data2.Config = nil
	err := p.db.WithContext(ctx).Where("product_id = ?", data.ProductID).Save(&data2).Error
	return stores.ErrFmt(err)
}

func (d ProductInfoRepo) UpdateWithField(ctx context.Context, f ProductFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&DmProductInfo{}).Updates(updates).Error
	return stores.ErrFmt(err)
}

func (p ProductInfoRepo) Delete(ctx context.Context, productID string) error {
	db := p.db.WithContext(ctx)
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("product_id=?", productID).Delete(&DmProductInfo{}).Error
		if err != nil {
			return err
		}
		return tx.Where("product_id=?", productID).Delete(&DmProductConfig{}).Error
	})
	return stores.ErrFmt(err)
}

func (p ProductInfoRepo) FindByFilter(ctx context.Context, f ProductFilter, page *stores.PageInfo) ([]*DmProductInfo, error) {
	var results []*DmProductInfo
	db := p.fmtFilter(ctx, f).Model(&DmProductInfo{})
	db = page.ToGorm(db)
	err := db.Preload("Config").Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ProductInfoRepo) CountByFilter(ctx context.Context, f ProductFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmProductInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}
