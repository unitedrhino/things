package relationDB

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"gorm.io/gorm"
)

type ProductInfoRepo struct {
	db *gorm.DB
}

type ProductFilter struct {
	DeviceType   int64
	DeviceTypes  []int64
	ProductName  string
	ProductIDs   []string
	ProductNames []string
	Tags         map[string]string
	ProtocolConf map[string]string
	WithProtocol bool
	WithCategory bool
	ProtocolCode string
	CategoryIDs  []int64
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
	if len(f.CategoryIDs) != 0 {
		db = db.Where("category_id in ?", f.CategoryIDs)
	}
	if f.ProtocolCode != "" {
		db = db.Where("protocol_code=?", f.ProtocolCode)
	}
	if f.WithProtocol {
		db = db.Preload("Protocol")
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
			db = db.Where("JSON_CONTAINS(protocol_conf, JSON_OBJECT(?,?))",
				k, v)
		}
	}
	if f.Tags != nil {
		for k, v := range f.Tags {
			db = db.Where("JSON_CONTAINS(tags, JSON_OBJECT(?,?))",
				k, v)
		}
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
			var schemas []*DmProductSchema
			for _, v := range ids {
				v.ID = 0
				v.Tag = schema.TagRequired
				schemas = append(schemas, &DmProductSchema{
					ProductID:    data.ProductID,
					DmSchemaCore: v.DmSchemaCore,
				})
			}
			err = NewProductSchemaRepo(tx).MultiInsert(ctx, schemas)
		}
		return err
	})
	return stores.ErrFmt(err)
}

func (p ProductInfoRepo) FindOneByFilter(ctx context.Context, f ProductFilter) (*DmProductInfo, error) {
	var result DmProductInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p ProductInfoRepo) Update(ctx context.Context, data *DmProductInfo) error {
	err := p.db.WithContext(ctx).Where("product_id = ?", data.ProductID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProductInfoRepo) DeleteByFilter(ctx context.Context, f ProductFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmProductInfo{}).Error
	return stores.ErrFmt(err)
}

func (p ProductInfoRepo) FindByFilter(ctx context.Context, f ProductFilter, page *def.PageInfo) ([]*DmProductInfo, error) {
	var results []*DmProductInfo
	db := p.fmtFilter(ctx, f).Model(&DmProductInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
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
