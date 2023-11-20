package relationDB

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/shared/utils"
	"gorm.io/gorm"
)

type VidmgrStreamRepo struct {
	db *gorm.DB
}

type VidmgrStreamFilter struct {
	VidmgrID string //流服务ID

	App    string
	Schema string
	Stream string
	Vhost  string

	Identifier string
	LocalIP    int64
	LocalPort  int64
	PeerIP     int64
	PeerPort   int64

	IsOnline      bool
	LastLoginTime struct {
		Start int64
		End   int64
	}

	StreamIDs []int64 //流ID
}

func NewVidmgrStreamRepo(in any) *VidmgrStreamRepo {
	return &VidmgrStreamRepo{db: stores.GetCommonConn(in)}
}

func (p VidmgrStreamRepo) fmtFilter(ctx context.Context, f VidmgrStreamFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.VidmgrID != "" {
		db = db.Where("vidmgr_id=?", f.VidmgrID)
	}
	if f.App != "" {
		db = db.Where("app=?", f.App)
	}
	if f.Schema != "" {
		db = db.Where("schema=?", f.Schema)
	}
	if f.Stream != "" {
		db = db.Where("stream=?", f.Stream)
	}
	if f.Vhost != "" {
		db = db.Where("vhost=?", f.Vhost)
	}
	if f.Identifier != "" {
		db = db.Where("identifier=?", f.Identifier)
	}
	if f.LocalIP != 0 {
		db = db.Where("local_ip=?", f.LocalIP)
	}
	if f.LocalPort != 0 {
		db = db.Where("local_port=?", f.LocalPort)
	}
	if f.PeerIP != 0 {
		db = db.Where("peer_ip=?", f.PeerIP)
	}
	if f.PeerPort != 0 {
		db = db.Where("peer_port=?", f.PeerPort)
	}
	if len(f.StreamIDs) != 0 {
		db = db.Where("stream_id=?", f.StreamIDs)
	}
	if f.IsOnline {
		db = db.Where("is_online=?", f.IsOnline)
	}
	if f.LastLoginTime.Start != 0 {
		db = db.Where("last_login >= ?", utils.ToYYMMddHHSS(f.LastLoginTime.Start*1000))
	}
	if f.LastLoginTime.End != 0 {
		db = db.Where("last_login <= ?", utils.ToYYMMddHHSS(f.LastLoginTime.End*1000))
	}
	return db
}

func (p VidmgrStreamRepo) Insert(ctx context.Context, data *VidmgrStream) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p VidmgrStreamRepo) FindOneByFilter(ctx context.Context, f VidmgrStreamFilter) (*VidmgrStream, error) {
	var result VidmgrStream
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p VidmgrStreamRepo) Update(ctx context.Context, data *VidmgrStream) error {
	err := p.db.WithContext(ctx).Where("mediaServerId = ?", data.StreamID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p VidmgrStreamRepo) DeleteByFilter(ctx context.Context, f VidmgrStreamFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&VidmgrStream{}).Error
	return stores.ErrFmt(err)
}

// 用于定时check是否有在线状态
func (p VidmgrStreamRepo) FindAllFilter(ctx context.Context, f VidmgrStreamFilter) ([]*VidmgrStream, error) {
	var results []*VidmgrStream
	db := p.fmtFilter(ctx, f).Model(&VidmgrStream{})
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p VidmgrStreamRepo) FindByFilter(ctx context.Context, f VidmgrStreamFilter, page *def.PageInfo) ([]*VidmgrStream, error) {
	var results []*VidmgrStream
	db := p.fmtFilter(ctx, f).Model(&VidmgrStream{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p VidmgrStreamRepo) CountByFilter(ctx context.Context, f VidmgrStreamFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&VidmgrStream{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (d VidmgrStreamRepo) CountStreamByField(ctx context.Context, f VidmgrStreamFilter, columnName string) (map[string]int64, error) {
	db := d.fmtFilter(ctx, f).Model(&VidmgrInfo{})
	countModelList := make([]*countModel, 0)
	err := db.Select(fmt.Sprintf("%s as CountKey", columnName), "count(1) as count").Group(columnName).Find(&countModelList).Error
	result := make(map[string]int64, 0)
	for _, v := range countModelList {
		result[v.CountKey] = v.Count
	}
	return result, stores.ErrFmt(err)
}
