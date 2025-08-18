package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/share/domain/schema"
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
		ControlMode       int64
		ProductSceneMode  string
	}
)

func NewProductSchemaRepo(in any) *ProductSchemaRepo {
	return &ProductSchemaRepo{db: stores.GetCommonConn(in)}
}

func (p ProductSchemaRepo) filter(db *gorm.DB, f ProductSchemaFilter) *gorm.DB {
	db = db.Where("tag !=?", schema.TagDeviceCustom)

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

func (p ProductSchemaRepo) fmtFilter(ctx context.Context, f ProductSchemaFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	return p.filter(db, f)
}
func (p ProductSchemaRepo) Insert(ctx context.Context, data *DmSchemaInfo) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		err := NewSchemaInfoRepo(tx).DeleteByFilter(ctx, SchemaInfoFilter{ProductID: data.ProductID, Tag: schema.TagDeviceCustom, Identifiers: []string{data.Identifier}})
		if err != nil {
			return err
		}
		return tx.Create(data).Error
	})
	return stores.ErrFmt(err)
}

func (p ProductSchemaRepo) FindOneByFilter(ctx context.Context, f ProductSchemaFilter) (*DmSchemaInfo, error) {
	var result DmSchemaInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p ProductSchemaRepo) Update(ctx context.Context, data *DmSchemaInfo) error {
	err := p.db.WithContext(ctx).Omit("product_id", "identifier").Where("product_id = ? and identifier = ?", data.ProductID, data.Identifier).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProductSchemaRepo) UpdateTag(ctx context.Context, productIDs []string, identifiers []string, oldTag, newTag int64) error {
	err := p.db.WithContext(ctx).Model(&DmSchemaInfo{}).Where(
		"product_id in ? and identifier in ? and tag =?", productIDs, identifiers, oldTag).Update("tag", newTag).Error
	return stores.ErrFmt(err)
}

func (p ProductSchemaRepo) UpdateWithCommon(ctx context.Context, common *DmCommonSchema) error {
	data := DmSchemaInfo{
		DmSchemaCore: DmSchemaCore{
			//ExtendConfig:      common.ExtendConfig,
			Name:              common.Name,
			Required:          common.Required,
			IsCanSceneLinkage: common.IsCanSceneLinkage,
			FuncGroup:         common.FuncGroup,
			ControlMode:       common.ControlMode,
			UserPerm:          common.UserPerm,
			RecordMode:        common.RecordMode,
			Affordance:        common.Affordance,
			Order:             common.Order,
		},
	}
	err := p.db.WithContext(ctx).Select("Order", "Name", "ControlMode", "Required", "IsCanSceneLinkage", "UserPerm", "FuncGroup", "IsHistory", "Affordance").
		Where("identifier = ? and (tag = ? or tag=?)",
			common.Identifier, schema.TagOptional, schema.TagRequired).Updates(&data).Error
	return stores.ErrFmt(err)
}

func (p ProductSchemaRepo) DeleteByFilter(ctx context.Context, f ProductSchemaFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmSchemaInfo{}).Error
	return stores.ErrFmt(err)
}

func (p ProductSchemaRepo) FindByFilter(ctx context.Context, f ProductSchemaFilter, page *stores.PageInfo) ([]*DmSchemaInfo, error) {
	var results []*DmSchemaInfo
	db := p.fmtFilter(ctx, f).Model(&DmSchemaInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}
func (p ProductSchemaRepo) FindProductIDByFilter(ctx context.Context, f ProductSchemaFilter) ([]string, error) {
	var results []*DmSchemaInfo
	db := p.fmtFilter(ctx, f).Model(&DmSchemaInfo{})
	err := db.Select("ProductID").Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return utils.ToSliceWithFunc(results, func(in *DmSchemaInfo) string {
		return in.ProductID
	}), nil

}

func (p ProductSchemaRepo) CountByFilter(ctx context.Context, f ProductSchemaFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmSchemaInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

// 批量插入 LightStrategyDevice 记录
func (p ProductSchemaRepo) MultiInsert(ctx context.Context, data []*DmSchemaInfo) error {
	var pmap = make(map[string][]string)
	for _, item := range data {
		p := pmap[item.ProductID]
		if p == nil {
			pmap[item.ProductID] = []string{item.Identifier}
			continue
		}
		pmap[item.ProductID] = append(p, item.Identifier)
	}
	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{UpdateAll: true,
			Columns: stores.SetColumnsWithPg(p.db, &DmSchemaInfo{}, "idx_dm_schema_info_identifier")}).Model(&DmSchemaInfo{}).Create(data).Error
		if err != nil {
			return err
		}
		for productID, ids := range pmap { //如果定义了设备物模型,需要清除
			err = NewDeviceSchemaRepo(tx).DeleteByFilter(ctx, DeviceSchemaFilter{ProductID: productID, Identifiers: ids})
			if err != nil {
				return err
			}
		}
		return err
	})
	return stores.ErrFmt(err)
}

func (p ProductSchemaRepo) MultiUpdate(ctx context.Context, productID string, schemaInfo *schema.Model) error {
	var datas []*DmSchemaInfo
	var idents []string

	for _, item := range schemaInfo.Property {
		idents = append(idents, item.Identifier)
		datas = append(datas, &DmSchemaInfo{
			ProductID:    productID,
			Identifier:   item.Identifier,
			DmSchemaCore: ToPropertyPo(item),
		})
	}
	for _, item := range schemaInfo.Event {
		idents = append(idents, item.Identifier)
		datas = append(datas, &DmSchemaInfo{
			ProductID:    productID,
			Identifier:   item.Identifier,
			DmSchemaCore: ToEventPo(item),
		})
	}
	for _, item := range schemaInfo.Action {
		idents = append(idents, item.Identifier)
		datas = append(datas, &DmSchemaInfo{
			ProductID:    productID,
			Identifier:   item.Identifier,
			DmSchemaCore: ToActionPo(item),
		})
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
		//如果定义了产品级的,需要删除设备级的
		err = NewSchemaInfoRepo(tx).DeleteByFilter(ctx, SchemaInfoFilter{ProductID: productID, Tag: schema.TagDeviceCustom,
			Identifiers: idents})
		if err != nil {
			return err
		}
		return nil
	})
	return stores.ErrFmt(err)
}
