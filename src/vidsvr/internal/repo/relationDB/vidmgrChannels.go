package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
)

// 国标设备-通道
type VidmgrChannelsRepo struct {
	db *gorm.DB
}

func NewVidmgrChannelsRepo(in any) *VidmgrChannelsRepo {
	return &VidmgrChannelsRepo{db: stores.GetCommonConn(in)}
}

type VidmgrChannelsFilter struct {
	ChannelIDs []string
	DeviceIDs  []string
	ChannelID  string
	MeMeo      string
	Model      string
}

func (p VidmgrChannelsRepo) fmtFilter(ctx context.Context, f VidmgrChannelsFilter) *gorm.DB {
	db := p.db.WithContext(ctx)

	if len(f.ChannelIDs) != 0 {
		db = db.Where("channel_id in?", f.ChannelIDs)
	} else if len(f.DeviceIDs) != 0 {
		db = db.Where("device_id in?", f.DeviceIDs)
	}

	if f.ChannelID != "" {
		db = db.Where("channel_id =?", f.ChannelID)
	}
	if f.MeMeo != "" {
		db = db.Where("memo =?", f.MeMeo)
	}
	if f.Model != "" {
		db = db.Where("model =?", f.Model)
	}
	return db
}

func (p VidmgrChannelsRepo) Insert(ctx context.Context, data *VidmgrChannels) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p VidmgrChannelsRepo) FindOneByFilter(ctx context.Context, f VidmgrChannelsFilter) (*VidmgrChannels, error) {
	var result VidmgrChannels
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p VidmgrChannelsRepo) Update(ctx context.Context, data *VidmgrChannels) error {
	err := p.db.WithContext(ctx).Where("channel_id = ?", data.ChannelID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p VidmgrChannelsRepo) DeleteByFilter(ctx context.Context, f VidmgrChannelsFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&VidmgrChannels{}).Error
	return stores.ErrFmt(err)
}

func (p VidmgrChannelsRepo) FindByFilter(ctx context.Context, f VidmgrChannelsFilter, page *def.PageInfo) ([]*VidmgrChannels, error) {
	var results []*VidmgrChannels
	db := p.fmtFilter(ctx, f).Model(&VidmgrChannels{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p VidmgrChannelsRepo) CountByFilter(ctx context.Context, f VidmgrChannelsFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&VidmgrChannels{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}
