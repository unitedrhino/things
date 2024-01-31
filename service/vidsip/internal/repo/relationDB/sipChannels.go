package relationDB

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"
)

// 国标设备-通道
type SipChannelsRepo struct {
	db *gorm.DB
}

func NewSipChannelsRepo(in any) *SipChannelsRepo {
	return &SipChannelsRepo{db: stores.GetCommonConn(in)}
}

type SipChannelsFilter struct {
	ChannelIDs []string
	DeviceIDs  []string
	ChannelID  string
	Stream     string
	MeMeo      string
	Model      string
}

func (p SipChannelsRepo) fmtFilter(ctx context.Context, f SipChannelsFilter) *gorm.DB {
	db := p.db.WithContext(ctx)

	if len(f.ChannelIDs) != 0 {
		db = db.Where("channel_id in?", f.ChannelIDs)
	} else if len(f.DeviceIDs) != 0 {
		db = db.Where("device_id in?", f.DeviceIDs)
	}

	if f.Stream != "" {
		db = db.Where("stream =?", f.Stream)
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

func (p SipChannelsRepo) Insert(ctx context.Context, data *SipChannels) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p SipChannelsRepo) FindOneByFilter(ctx context.Context, f SipChannelsFilter) (*SipChannels, error) {
	var result SipChannels
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p SipChannelsRepo) Update(ctx context.Context, data *SipChannels) error {
	err := p.db.WithContext(ctx).Where("channel_id = ?", data.ChannelID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p SipChannelsRepo) DeleteByFilter(ctx context.Context, f SipChannelsFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SipChannels{}).Error
	return stores.ErrFmt(err)
}

func (p SipChannelsRepo) FindByFilter(ctx context.Context, f SipChannelsFilter, page *def.PageInfo) ([]*SipChannels, error) {
	var results []*SipChannels
	db := p.fmtFilter(ctx, f).Model(&SipChannels{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p SipChannelsRepo) CountByFilter(ctx context.Context, f SipChannelsFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SipChannels{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}
