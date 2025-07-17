package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductSchemaOldRepo struct {
	db *gorm.DB
}

type (
	ProductSchemaOldFilter struct {
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
		ControlMode       int64
		ProductSceneMode  string
	}
)

func NewProductSchemaOldRepo(in any) *ProductSchemaOldRepo {
	return &ProductSchemaOldRepo{db: stores.GetCommonConn(in)}
}

func (p ProductSchemaOldRepo) fmtFilter(ctx context.Context, f ProductSchemaFilter) *gorm.DB {
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
	if f.ControlMode != 0 {
		db = db.Where("control_mode=?", f.ControlMode)
	}
	if f.PropertyMode != "" {
		db = stores.CmpJsonObjEq("mode", f.PropertyMode).Where(db, "affordance")

	}
	if len(f.ProductIDs) != 0 {
		db = db.Where("product_id in ?", f.ProductIDs)
	}
	if f.ProductID != "" {
		db = db.Where("product_id=?", f.ProductID)
	}
	if f.ProductSceneMode != "" {
		db = db.Where("product_id in (?)", p.db.Select("product_id").Model(DmProductInfo{}).Where("scene_mode = ?", f.ProductSceneMode))
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
func (p ProductSchemaOldRepo) Insert(ctx context.Context, data *DmProductSchema) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ProductSchemaOldRepo) FindOneByFilter(ctx context.Context, f ProductSchemaFilter) (*DmProductSchema, error) {
	var result DmProductSchema
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p ProductSchemaOldRepo) Update(ctx context.Context, data *DmProductSchema) error {
	err := p.db.WithContext(ctx).Omit("product_id", "identifier").Where("product_id = ? and identifier = ?", data.ProductID, data.Identifier).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProductSchemaOldRepo) UpdateTag(ctx context.Context, productIDs []string, identifiers []string, oldTag, newTag int64) error {
	err := p.db.WithContext(ctx).Model(&DmProductSchema{}).Where(
		"product_id in ? and identifier in ? and tag =?", productIDs, identifiers, oldTag).Update("tag", newTag).Error
	return stores.ErrFmt(err)
}

func (p ProductSchemaOldRepo) UpdateWithCommon(ctx context.Context, common *DmCommonSchema) error {
	data := DmProductSchema{
		DmSchemaCore: DmSchemaCore{
			//ExtendConfig:      common.ExtendConfig,
			Name:              common.Name,
			Required:          common.Required,
			IsCanSceneLinkage: common.IsCanSceneLinkage,
			FuncGroup:         common.FuncGroup,
			ControlMode:       common.ControlMode,
			UserPerm:          common.UserPerm,
			IsHistory:         common.IsHistory,
			Affordance:        common.Affordance,
		},
	}
	err := p.db.WithContext(ctx).Select("Name", "ControlMode", "ExtendConfig", "Required", "IsCanSceneLinkage", "UserPerm", "FuncGroup", "IsHistory", "Affordance").
		Where("identifier = ? and (tag = ? or tag=?)",
			common.Identifier, schema.TagOptional, schema.TagRequired).Updates(&data).Error
	return stores.ErrFmt(err)
}

func (p ProductSchemaOldRepo) DeleteByFilter(ctx context.Context, f ProductSchemaFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmProductSchema{}).Error
	return stores.ErrFmt(err)
}

func (p ProductSchemaOldRepo) FindByFilter(ctx context.Context, f ProductSchemaFilter, page *stores.PageInfo) ([]*DmProductSchema, error) {
	var results []*DmProductSchema
	db := p.fmtFilter(ctx, f).Model(&DmProductSchema{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}
func (p ProductSchemaOldRepo) FindProductIDByFilter(ctx context.Context, f ProductSchemaFilter) ([]string, error) {
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

func (p ProductSchemaOldRepo) CountByFilter(ctx context.Context, f ProductSchemaFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmProductSchema{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

// 批量插入 LightStrategyDevice 记录
func (p ProductSchemaOldRepo) MultiInsert(ctx context.Context, data []*DmProductSchema) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmProductSchema{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p ProductSchemaOldRepo) MultiUpdate(ctx context.Context, productID string, schemaInfo *schema.Model) error {
	var datas []*DmProductSchema
	for _, item := range schemaInfo.Property {
		datas = append(datas, &DmProductSchema{
			ProductID:    productID,
			Identifier:   item.Identifier,
			DmSchemaCore: ToPropertyPo(item),
		})
	}
	for _, item := range schemaInfo.Event {
		datas = append(datas, &DmProductSchema{
			ProductID:    productID,
			Identifier:   item.Identifier,
			DmSchemaCore: ToEventPo(item),
		})
	}
	for _, item := range schemaInfo.Action {
		datas = append(datas, &DmProductSchema{
			ProductID:    productID,
			Identifier:   item.Identifier,
			DmSchemaCore: ToActionPo(item),
		})
	}
	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewProductSchemaOldRepo(tx)
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
