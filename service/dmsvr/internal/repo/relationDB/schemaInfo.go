package relationDB

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/domain/schema"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SchemaInfoRepo struct {
	db *gorm.DB
}

type (
	SchemaInfoFilter struct {
		ID                int64
		ProductID         string //产品id  必填
		DeviceName        string //设备ID
		Type              int64  //物模型类型 1:property属性 2:event事件 3:action行为
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

func NewSchemaInfoRepo(in any) *SchemaInfoRepo {
	return &SchemaInfoRepo{db: stores.GetCommonConn(in)}
}

func (p SchemaInfoRepo) fmtFilter(ctx context.Context, f SchemaInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	db = db.Where("product_id=?", f.ProductID)
	db = db.Where("device_name=?", f.DeviceName)

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
		db = db.Where("JSON_CONTAINS(affordance, JSON_OBJECT('mode',?))",
			f.PropertyMode)
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
func (p SchemaInfoRepo) Insert(ctx context.Context, data *DmSchemaInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p SchemaInfoRepo) FindOneByFilter(ctx context.Context, f SchemaInfoFilter) (*DmSchemaInfo, error) {
	var result DmSchemaInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p SchemaInfoRepo) Update(ctx context.Context, data *DmSchemaInfo) error {
	err := p.db.WithContext(ctx).Omit("product_id", "identifier").Where("product_id = ? and identifier = ?", data.ProductID, data.Identifier).Save(data).Error
	return stores.ErrFmt(err)
}

func (p SchemaInfoRepo) UpdateTag(ctx context.Context, productIDs []string, identifiers []string, oldTag, newTag int64) error {
	err := p.db.WithContext(ctx).Model(&DmSchemaInfo{}).Where(
		"product_id in ? and identifier in ? and tag =?", productIDs, identifiers, oldTag).Update("tag", newTag).Error
	return stores.ErrFmt(err)
}

func (p SchemaInfoRepo) DeleteByFilter(ctx context.Context, f SchemaInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmSchemaInfo{}).Error
	return stores.ErrFmt(err)
}

func (p SchemaInfoRepo) FindByFilter(ctx context.Context, f SchemaInfoFilter, page *stores.PageInfo) ([]*DmSchemaInfo, error) {
	var results []*DmSchemaInfo
	db := p.fmtFilter(ctx, f).Model(&DmSchemaInfo{})
	db = page.ToGorm(db)
	newDB := NewProductSchemaRepo(ctx).fmtFilter(ctx, utils.Copy2[ProductSchemaFilter](f))
	db2 := p.db.Raw("(?) union all (?)", db, newDB.Select(fmt.Sprintf("'%v' as device_name,dm_product_schema.*", f.DeviceName)).Model(&DmProductSchema{}))
	db2 = page.ToGorm(db2)
	//var rst = []map[string]any{}
	err := db2.Scan(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p SchemaInfoRepo) CountByFilter(ctx context.Context, f SchemaInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmSchemaInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

// 批量插入 LightStrategyDevice 记录
func (p SchemaInfoRepo) MultiInsert(ctx context.Context, data []*DmSchemaInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmSchemaInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p SchemaInfoRepo) MultiUpdate(ctx context.Context, productID string, deviceName string, schemaInfo *schema.Model) error {
	var datas []*DmSchemaInfo
	//toD := func(in *DmProductSchema) *DmSchemaInfo {
	//	return &DmSchemaInfo{
	//		DeviceName:      deviceName,
	//		DmProductSchema: *in,
	//	}
	//}
	//for _, item := range schemaInfo.Property {
	//	datas = append(datas, toD(ToPropertyPo(productID, item)))
	//}
	//for _, item := range schemaInfo.Event {
	//	datas = append(datas, toD(ToEventPo(productID, item)))
	//}
	//for _, item := range schemaInfo.Action {
	//	datas = append(datas, toD(ToActionPo(productID, item)))
	//}
	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewSchemaInfoRepo(tx)
		err := rm.DeleteByFilter(ctx, SchemaInfoFilter{ProductID: productID})
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
