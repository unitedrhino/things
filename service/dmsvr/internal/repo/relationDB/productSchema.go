package relationDB

import (
	"context"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductSchemaRepo struct {
	db *gorm.DB
}

type (
	ProductSchemaFilter struct {
		ID                int64
		ProductID         string   //产品id  必填
		ProductIDs        []string //产品id列表
		Type              int64    //物模型类型 1:property属性 2:event事件 3:action行为
		Types             []int64
		Tag               schema.Tag //过滤条件: 物模型标签 1:自定义 2:可选 3:必选
		Tags              []schema.Tag
		Identifiers       []string //过滤标识符列表
		Name              string
		IsCanSceneLinkage int64
		FuncGroup         int64
		UserPerm          int64
		PropertyMode      string
	}
)

func NewProductSchemaRepo(in any) *ProductSchemaRepo {
	return &ProductSchemaRepo{db: stores.GetCommonConn(in)}
}

func (p ProductSchemaRepo) fmtFilter(ctx context.Context, f ProductSchemaFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.IsCanSceneLinkage != 0 {
		db = db.Where("is_can_scene_linkage = ?", f.IsCanSceneLinkage)
	}
	if f.FuncGroup != 0 {
		db = db.Where("func_group = ?", f.FuncGroup)
	}
	if f.UserPerm != 0 {
		db = db.Where("user_auth = ?", f.UserPerm)
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	if f.ID != 0 {
		db = db.Where("id=?", f.ID)
	}
	if f.PropertyMode != "" {
		db = db.Where("JSON_CONTAINS(affordance, JSON_OBJECT('mode',?))",
			f.PropertyMode)
	}
	if len(f.ProductIDs) != 0 {
		db = db.Where("product_id in ?", f.ProductIDs)
	}
	if f.ProductID != "" {
		db = db.Where("product_id=?", f.ProductID)
	}
	if f.Type != 0 {
		db = db.Where("type=?", f.Type)
	}
	if len(f.Types) != 0 {
		db = db.Where("type in ?", f.Types)
	}
	if f.Tag != 0 {
		db = db.Where("tag=?", f.Tag)
	}
	if len(f.Tags) != 0 {
		db = db.Where("tag in ?", f.Tags)
	}
	if len(f.Identifiers) != 0 {
		db = db.Where("identifier in ?", f.Identifiers)
	}
	return db
}
func (p ProductSchemaRepo) Insert(ctx context.Context, data *DmProductSchema) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ProductSchemaRepo) FindOneByFilter(ctx context.Context, f ProductSchemaFilter) (*DmProductSchema, error) {
	var result DmProductSchema
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p ProductSchemaRepo) Update(ctx context.Context, data *DmProductSchema) error {
	err := p.db.WithContext(ctx).Omit("product_id", "identifier").Where("product_id = ? and identifier = ?", data.ProductID, data.Identifier).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProductSchemaRepo) UpdateTag(ctx context.Context, productIDs []string, identifiers []string, oldTag, newTag int64) error {
	err := p.db.WithContext(ctx).Model(&DmProductSchema{}).Where(
		"product_id in ? and identifier in ? and tag =?", productIDs, identifiers, oldTag).Update("tag", newTag).Error
	return stores.ErrFmt(err)
}

func (p ProductSchemaRepo) UpdateWithCommon(ctx context.Context, common *DmCommonSchema) error {
	data := DmProductSchema{
		DmSchemaCore: DmSchemaCore{
			ExtendConfig:      common.ExtendConfig,
			Name:              common.Name,
			Required:          common.Required,
			IsCanSceneLinkage: common.IsCanSceneLinkage,
			FuncGroup:         common.FuncGroup,
			UserPerm:          common.UserPerm,
			IsHistory:         common.IsHistory,
			Affordance:        common.Affordance,
		},
	}
	err := p.db.WithContext(ctx).Select("Name", "ExtendConfig", "Required", "IsCanSceneLinkage", "UserPerm", "FuncGroup", "IsHistory", "Affordance").
		Where("identifier = ? and (tag = ? or tag=?)",
			common.Identifier, schema.TagOptional, schema.TagRequired).Updates(&data).Error
	return stores.ErrFmt(err)
}

func (p ProductSchemaRepo) DeleteByFilter(ctx context.Context, f ProductSchemaFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmProductSchema{}).Error
	return stores.ErrFmt(err)
}

func (p ProductSchemaRepo) FindByFilter(ctx context.Context, f ProductSchemaFilter, page *stores.PageInfo) ([]*DmProductSchema, error) {
	var results []*DmProductSchema
	db := p.fmtFilter(ctx, f).Model(&DmProductSchema{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}
func (p ProductSchemaRepo) FindProductIDByFilter(ctx context.Context, f ProductSchemaFilter) ([]string, error) {
	var results []*DmProductSchema
	db := p.fmtFilter(ctx, f).Model(&DmProductSchema{})
	err := db.Select("ProductID").Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return utils.ToSliceWithFunc(results, func(in *DmProductSchema) string {
		return in.ProductID
	}), nil

}

func (p ProductSchemaRepo) CountByFilter(ctx context.Context, f ProductSchemaFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmProductSchema{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

// 批量插入 LightStrategyDevice 记录
func (p ProductSchemaRepo) MultiInsert(ctx context.Context, data []*DmProductSchema) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmProductSchema{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p ProductSchemaRepo) MultiUpdate(ctx context.Context, productID string, schemaInfo *schema.Model) error {
	var datas []*DmProductSchema
	for _, item := range schemaInfo.Property {
		datas = append(datas, ToPropertyPo(productID, item))
	}
	for _, item := range schemaInfo.Event {
		datas = append(datas, ToEventPo(productID, item))
	}
	for _, item := range schemaInfo.Action {
		datas = append(datas, ToActionPo(productID, item))
	}
	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewProductSchemaRepo(tx)
		err := rm.DeleteByFilter(ctx, ProductSchemaFilter{ProductID: productID})
		if err != nil {
			return err
		}
		if len(datas) != 0 {
			err = rm.MultiInsert(ctx, datas)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return stores.ErrFmt(err)
}
