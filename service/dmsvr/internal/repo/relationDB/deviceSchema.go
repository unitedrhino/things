package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/domain/schema"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeviceSchemaRepo struct {
	db *gorm.DB
}

type (
	DeviceSchemaFilter struct {
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
		WithProductSchema bool
	}
)

func NewDeviceSchemaRepo(in any) *DeviceSchemaRepo {
	return &DeviceSchemaRepo{db: stores.GetCommonConn(in)}
}

func (p DeviceSchemaRepo) filter(db *gorm.DB, f DeviceSchemaFilter) *gorm.DB {
	if !f.WithProductSchema {
		db = db.Where("product_id=?", f.ProductID)
		db = db.Where("device_name=?", f.DeviceName)
		db = db.Where("tag=?", schema.TagDevice)
	} else {
		db = db.Where("product_id=? and device_name=? and tag=?", f.ProductID, f.DeviceName, schema.TagDevice).
			Or("product_id=? and tag !=?", f.ProductID, schema.TagDevice)

	}

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

func (p DeviceSchemaRepo) fmtFilter(ctx context.Context, f DeviceSchemaFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	return p.filter(db, f)
}
func (p DeviceSchemaRepo) Insert(ctx context.Context, data *DmDeviceSchema) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p DeviceSchemaRepo) FindOneByFilter(ctx context.Context, f DeviceSchemaFilter) (*DmDeviceSchema, error) {
	var result DmDeviceSchema
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p DeviceSchemaRepo) Update(ctx context.Context, data *DmDeviceSchema) error {
	err := p.db.WithContext(ctx).Omit("product_id", "device_name", "identifier").Where("product_id = ? and device_name=?  and identifier = ?", data.ProductID, data.DeviceName, data.Identifier).Save(data).Error
	return stores.ErrFmt(err)
}

func (p DeviceSchemaRepo) UpdateTag(ctx context.Context, productIDs []string, identifiers []string, oldTag, newTag int64) error {
	err := p.db.WithContext(ctx).Model(&DmDeviceSchema{}).Where(
		"product_id in ? and identifier in ? and tag =?", productIDs, identifiers, oldTag).Update("tag", newTag).Error
	return stores.ErrFmt(err)
}

func (p DeviceSchemaRepo) DeleteByFilter(ctx context.Context, f DeviceSchemaFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmDeviceSchema{}).Error
	return stores.ErrFmt(err)
}

func (p DeviceSchemaRepo) FindByFilter(ctx context.Context, f DeviceSchemaFilter, page *stores.PageInfo) ([]*DmDeviceSchema, error) {
	var results []*DmDeviceSchema
	var db *gorm.DB
	db = p.fmtFilter(ctx, f)
	db = db.Model(&DmDeviceSchema{})
	db = page.ToGorm(db)
	//var rst = []map[string]any{}
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p DeviceSchemaRepo) CountByFilter(ctx context.Context, f DeviceSchemaFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmDeviceSchema{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

// 批量插入 LightStrategyDevice 记录
func (p DeviceSchemaRepo) MultiInsert(ctx context.Context, data []*DmDeviceSchema) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmDeviceSchema{}).Create(data).Error
	return stores.ErrFmt(err)
}

// 批量插入 LightStrategyDevice 记录
func (p DeviceSchemaRepo) MultiInsert2(ctx context.Context, productID string, deviceName string, schemaInfo *schema.Model) error {
	var datas []*DmDeviceSchema
	for _, item := range schemaInfo.Property {
		datas = append(datas, &DmDeviceSchema{
			ProductID:    productID,
			DeviceName:   deviceName,
			DmSchemaCore: ToPropertyPo(item),
		})
	}
	for _, item := range schemaInfo.Event {
		datas = append(datas, &DmDeviceSchema{
			ProductID:    productID,
			DeviceName:   deviceName,
			DmSchemaCore: ToEventPo(item),
		})
	}
	for _, item := range schemaInfo.Action {
		datas = append(datas, &DmDeviceSchema{
			ProductID:    productID,
			DeviceName:   deviceName,
			DmSchemaCore: ToActionPo(item),
		})
	}
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Model(&DmDeviceSchema{}).CreateInBatches(datas, 100).Error
	return stores.ErrFmt(err)
}

func (p DeviceSchemaRepo) MultiUpdate(ctx context.Context, productID string, deviceName string, schemaInfo *schema.Model) error {
	var datas []*DmDeviceSchema
	for _, item := range schemaInfo.Property {
		item.Tag = schema.TagDevice
		datas = append(datas, &DmDeviceSchema{
			ProductID:    productID,
			DeviceName:   deviceName,
			DmSchemaCore: ToPropertyPo(item),
		})
	}
	for _, item := range schemaInfo.Event {
		item.Tag = schema.TagDevice
		datas = append(datas, &DmDeviceSchema{
			ProductID:    productID,
			DeviceName:   deviceName,
			DmSchemaCore: ToEventPo(item),
		})
	}
	for _, item := range schemaInfo.Action {
		item.Tag = schema.TagDevice
		datas = append(datas, &DmDeviceSchema{
			ProductID:    productID,
			DeviceName:   deviceName,
			DmSchemaCore: ToActionPo(item),
		})
	}
	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewDeviceSchemaRepo(tx)
		err := rm.DeleteByFilter(ctx, DeviceSchemaFilter{ProductID: productID})
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
