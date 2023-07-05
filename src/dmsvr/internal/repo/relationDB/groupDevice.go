package relationDB

import (
	"context"
	"database/sql"
	"github.com/i-Things/things/shared/store"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"gorm.io/gorm"
)

type GroupDeviceRepo struct {
	db *gorm.DB
}

func NewGroupDeviceRepo(in any) *GroupDeviceRepo {
	return &GroupDeviceRepo{db: store.GetCommonConn(in)}
}

func (g GroupDeviceRepo) Insert(ctx context.Context, data *mysql.DmGroupDevice) (sql.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (g GroupDeviceRepo) FindOne(ctx context.Context, id int64) (*mysql.DmGroupDevice, error) {
	//TODO implement me
	panic("implement me")
}

func (g GroupDeviceRepo) FindOneByGroupIDProductIDDeviceName(ctx context.Context, groupID int64, productID string, deviceName string) (*mysql.DmGroupDevice, error) {
	//TODO implement me
	panic("implement me")
}

func (g GroupDeviceRepo) Update(ctx context.Context, data *mysql.DmGroupDevice) error {
	//TODO implement me
	panic("implement me")
}

func (g GroupDeviceRepo) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}
