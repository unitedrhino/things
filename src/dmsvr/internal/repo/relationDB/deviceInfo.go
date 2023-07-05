package relationDB

import (
	"context"
	"database/sql"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/store"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"gorm.io/gorm"
)

type DeviceInfoRepo struct {
	db *gorm.DB
}

func NewDeviceInfoRepo(in any) *DeviceInfoRepo {
	return &DeviceInfoRepo{db: store.GetCommonConn(in)}
}

func (d DeviceInfoRepo) Insert(ctx context.Context, data *mysql.DmDeviceInfo) (sql.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (d DeviceInfoRepo) FindOne(ctx context.Context, id int64) (*mysql.DmDeviceInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (d DeviceInfoRepo) FindOneByIccid(ctx context.Context, iccid sql.NullString) (*mysql.DmDeviceInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (d DeviceInfoRepo) FindOneByPhone(ctx context.Context, phone sql.NullString) (*mysql.DmDeviceInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (d DeviceInfoRepo) FindOneByProductIDDeviceName(ctx context.Context, productID string, deviceName string) (*mysql.DmDeviceInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (d DeviceInfoRepo) Update(ctx context.Context, data *mysql.DmDeviceInfo) error {
	//TODO implement me
	panic("implement me")
}

func (d DeviceInfoRepo) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (d DeviceInfoRepo) InsertDeviceInfo(ctx context.Context, data *mysql.DmDeviceInfo) error {
	//TODO implement me
	panic("implement me")
}

func (d DeviceInfoRepo) FindByFilter(ctx context.Context, filter mysql.DeviceFilter, page *def.PageInfo) ([]*mysql.DmDeviceInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (d DeviceInfoRepo) CountByFilter(ctx context.Context, filter mysql.DeviceFilter) (size int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (d DeviceInfoRepo) CountGroupByField(ctx context.Context, filter mysql.DeviceFilter, fieldName string) (map[string]int64, error) {
	//TODO implement me
	panic("implement me")
}

func (d DeviceInfoRepo) UpdateDeviceInfo(ctx context.Context, data *mysql.DmDeviceInfo) error {
	//TODO implement me
	panic("implement me")
}
