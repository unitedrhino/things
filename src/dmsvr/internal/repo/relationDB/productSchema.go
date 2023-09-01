package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
)

type ProductSchemaRepo struct {
	db *gorm.DB
}

type (
	ProductSchemaFilter struct {
		ID          int64
		ProductID   string   //产品id  必填
		Type        int64    //物模型类型 1:property属性 2:event事件 3:action行为
		Tag         int64    //过滤条件: 物模型标签 1:自定义 2:可选 3:必选
		Identifiers []string //过滤标识符列表
	}
	PropertyDef struct {
		IsUseShadow bool                `json:"isUseShadow"` //是否使用设备影子
		IsNoRecord  bool                `json:"isNoRecord"`  //不存储历史记录
		Define      schema.Define       `json:"define"`      //数据定义
		Mode        schema.PropertyMode `json:"mode"`        //读写类型: 1:r(只读) 2:rw(可读可写)
	}
	EventDef struct {
		Type   schema.EventType `json:"type"`   //事件类型: 1:信息:info  2:告警alert  3:故障:fault
		Params schema.Params    `json:"params"` //事件参数
	}
	ActionDef struct {
		Dir    schema.ActionDir `json:"dir"`    //调用方向
		Input  schema.Params    `json:"input"`  //调用参数
		Output schema.Params    `json:"output"` //返回参数
	}
)

func NewProductSchemaRepo(in any) *ProductSchemaRepo {
	return &ProductSchemaRepo{db: stores.GetCommonConn(in)}
}

func (p ProductSchemaRepo) fmtFilter(ctx context.Context, f ProductSchemaFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.ID != 0 {
		db = db.Where("id=?", f.ID)
	}
	if f.ProductID != "" {
		db = db.Where("product_id=?", f.ProductID)
	}
	if f.Type != 0 {
		db = db.Where("type=?", f.Type)
	}
	if f.Tag != 0 {
		db = db.Where("tag=?", f.Tag)
	}
	if len(f.Identifiers) != 0 {
		db = db.Where("identifier = ?", f.Identifiers)
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
	err := p.db.WithContext(ctx).Where("product_id = ?", data.ProductID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProductSchemaRepo) DeleteByFilter(ctx context.Context, f ProductSchemaFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmProductSchema{}).Error
	return stores.ErrFmt(err)
}

func (p ProductSchemaRepo) FindByFilter(ctx context.Context, f ProductSchemaFilter, page *def.PageInfo) ([]*DmProductSchema, error) {
	var results []*DmProductSchema
	db := p.fmtFilter(ctx, f).Model(&DmProductSchema{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ProductSchemaRepo) CountByFilter(ctx context.Context, f ProductSchemaFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmProductSchema{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}
