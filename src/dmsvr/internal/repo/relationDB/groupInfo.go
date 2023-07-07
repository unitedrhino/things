package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/store"
	"gorm.io/gorm"
)

type GroupInfoRepo struct {
	db *gorm.DB
}

type GroupInfoFilter struct {
	GroupID     int64
	GroupNames  []string
	ParentID    int64
	GroupName   string
	Tags        map[string]string
	WithProduct bool
}

func NewGroupInfoRepo(in any) *GroupInfoRepo {
	return &GroupInfoRepo{db: store.GetCommonConn(in)}
}

func (p GroupInfoRepo) fmtFilter(ctx context.Context, f GroupInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.WithProduct {
		db = db.Preload("ProductInfo")
	}
	if f.GroupID != 0 {
		db = db.Where("`groupID`=?", f.GroupID)
	}
	if len(f.GroupNames) != 0 {
		db = db.Where("`groupName` in ?", f.GroupNames)
	}
	if f.ParentID != 0 {
		db = db.Where("`parentID`=?", f.ParentID)
	}
	if f.GroupName != "" {
		db = db.Where("`groupName` like ?", "%"+f.GroupName+"%")
	}
	if f.Tags != nil {
		for k, v := range f.Tags {
			db = db.Where("JSON_CONTAINS(`tags`, JSON_OBJECT(?,?))",
				k, v)
		}
	}
	return db
}

func (g GroupInfoRepo) Insert(ctx context.Context, data *DmGroupInfo) error {
	result := g.db.WithContext(ctx).Create(data)
	return store.ErrFmt(result.Error)
}

func (g GroupInfoRepo) FindOneByFilter(ctx context.Context, f GroupInfoFilter) (*DmGroupInfo, error) {
	var result DmGroupInfo
	db := g.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, store.ErrFmt(err)
	}
	return &result, nil
}
func (p GroupInfoRepo) FindByFilter(ctx context.Context, f GroupInfoFilter, page *def.PageInfo) ([]*DmGroupInfo, error) {
	var results []*DmGroupInfo
	db := p.fmtFilter(ctx, f).Model(&DmGroupInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, store.ErrFmt(err)
	}
	return results, nil
}

func (p GroupInfoRepo) CountByFilter(ctx context.Context, f GroupInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmGroupInfo{})
	err = db.Count(&size).Error
	return size, store.ErrFmt(err)
}

func (g GroupInfoRepo) Update(ctx context.Context, data *DmGroupInfo) error {
	err := g.db.WithContext(ctx).Where("`groupID` = ?", data.GroupID).Save(data).Error
	return store.ErrFmt(err)
}

func (g GroupInfoRepo) DeleteByFilter(ctx context.Context, f GroupInfoFilter) error {
	db := g.fmtFilter(ctx, f)
	err := db.Delete(&DmGroupInfo{}).Error
	return store.ErrFmt(err)
}
