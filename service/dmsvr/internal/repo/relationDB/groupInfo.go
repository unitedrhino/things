package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/share/devices"
	"gorm.io/gorm"
)

type GroupInfoRepo struct {
	db *gorm.DB
}

type GroupInfoFilter struct {
	AreaID      int64
	ID          int64
	IDs         []int64
	Names       []string
	Purpose     string
	Purposes    []string
	ParentID    int64
	IDPath      string
	Name        string
	Tags        map[string]string
	HasDevice   *devices.Core
	HasDevices  []*devices.Core
	WithProduct bool
}

func NewGroupInfoRepo(in any) *GroupInfoRepo {
	return &GroupInfoRepo{db: stores.GetCommonConn(in)}
}

func (p GroupInfoRepo) fmtFilter(ctx context.Context, f GroupInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.WithProduct {
		db = db.Preload("ProductInfo")
	}
	if f.Purpose == "" && len(f.Purposes) == 0 {
		f.Purpose = "default"
	}
	if len(f.Purposes) != 0 {
		db = db.Where("purpose in ?", f.Purposes)
	} else {
		db = db.Where("purpose = ?", f.Purpose)
	}
	if f.HasDevice != nil {
		subQuery := p.db.Model(&DmGroupDevice{}).Select("group_id").Where("product_id=? and device_name=?",
			f.HasDevice.ProductID, f.HasDevice.DeviceName)
		db = db.Where("id in (?)", subQuery)
	}
	if len(f.HasDevices) != 0 {
		scope := func(db *gorm.DB) *gorm.DB {
			for i, d := range f.HasDevices {
				if i == 0 {
					db = db.Where("product_id = ? and device_name = ?", d.ProductID, d.DeviceName)
					continue
				}
				db = db.Or("product_id = ? and device_name = ?", d.ProductID, d.DeviceName)
			}
			return db
		}
		subQuery := p.db.Model(&DmGroupDevice{}).Select("group_id")
		subQuery = subQuery.Where(scope(subQuery))
		db = db.Where("id in (?)", subQuery)
	}
	if f.ID != 0 {
		db = db.Where("id = ?", f.ID)
	}
	if len(f.IDs) != 0 {
		db = db.Where("id in ?", f.IDs)
	}
	if f.AreaID != 0 {
		db = db.Where("area_id=?", f.AreaID)
	}
	if len(f.Names) != 0 {
		db = db.Where("name in ?", f.Names)
	}
	if f.ParentID != 0 {
		db = db.Where("parent_id=?", f.ParentID)
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	if f.IDPath != "" {
		db = db.Where("id_path like ?", f.IDPath+"%")
	}
	if f.Tags != nil {
		for k, v := range f.Tags {
			db = stores.CmpJsonObjEq(k, v).Where(db, "tags")
		}
	}
	return db
}

func (g GroupInfoRepo) Insert(ctx context.Context, data *DmGroupInfo) error {
	result := g.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (g GroupInfoRepo) UpdateGroupDeviceCount(ctx context.Context, id int64) error {
	subQuery := g.db.WithContext(ctx).Model(&DmGroupDevice{}).Select("count(1)").Where("group_id=?", id)
	result := g.db.WithContext(ctx).Model(&DmGroupInfo{}).Where("id=?", id).Update("device_count", subQuery)
	return stores.ErrFmt(result.Error)
}

func (g GroupInfoRepo) FindOneByFilter(ctx context.Context, f GroupInfoFilter) (*DmGroupInfo, error) {
	var result DmGroupInfo
	db := g.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p GroupInfoRepo) FindOne(ctx context.Context, id int64) (*DmGroupInfo, error) {
	var result DmGroupInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p GroupInfoRepo) FindByFilter(ctx context.Context, f GroupInfoFilter, page *stores.PageInfo) ([]*DmGroupInfo, error) {
	var results []*DmGroupInfo
	db := p.fmtFilter(ctx, f).Model(&DmGroupInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p GroupInfoRepo) CountByFilter(ctx context.Context, f GroupInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmGroupInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (g GroupInfoRepo) Update(ctx context.Context, data *DmGroupInfo) error {
	err := g.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (d GroupInfoRepo) UpdateWithField(ctx context.Context, f GroupInfoFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&DmGroupInfo{}).Updates(updates).Error
	return stores.ErrFmt(err)
}

func (g GroupInfoRepo) DeleteByFilter(ctx context.Context, f GroupInfoFilter) error {
	db := g.fmtFilter(ctx, f)
	err := db.Delete(&DmGroupInfo{}).Error
	return stores.ErrFmt(err)
}
func (p GroupInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmGroupInfo{}).Error
	return stores.ErrFmt(err)
}
