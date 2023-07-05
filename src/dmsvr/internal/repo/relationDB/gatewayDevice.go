package relationDB

import (
	"context"
	"database/sql"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/store"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"gorm.io/gorm"
)

type GatewayDeviceRepo struct {
	db *gorm.DB
}

func NewGatewayDeviceRepo(in any) *GatewayDeviceRepo {
	return &GatewayDeviceRepo{db: store.GetCommonConn(in)}
}

func (g GatewayDeviceRepo) Insert(ctx context.Context, data *mysql.DmGatewayDevice) (sql.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (g GatewayDeviceRepo) FindOne(ctx context.Context, id int64) (*mysql.DmGatewayDevice, error) {
	//TODO implement me
	panic("implement me")
}

func (g GatewayDeviceRepo) FindOneByGatewayProductIDGatewayDeviceNameProductIDDeviceName(ctx context.Context, gatewayProductID string, gatewayDeviceName string, productID string, deviceName string) (*mysql.DmGatewayDevice, error) {
	//TODO implement me
	panic("implement me")
}

func (g GatewayDeviceRepo) Update(ctx context.Context, data *mysql.DmGatewayDevice) error {
	//TODO implement me
	panic("implement me")
}

func (g GatewayDeviceRepo) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (g GatewayDeviceRepo) CreateList(ctx context.Context, gateway *devices.Core, subDevices []*devices.Core) error {
	//TODO implement me
	panic("implement me")
}

func (g GatewayDeviceRepo) DeleteList(ctx context.Context, gateway *devices.Core, subDevices []*devices.Core) error {
	//TODO implement me
	panic("implement me")
}

func (g GatewayDeviceRepo) FindByFilter(ctx context.Context, filter mysql.GatewayDeviceFilter, page *def.PageInfo) ([]*mysql.DmDeviceInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (g GatewayDeviceRepo) CountByFilter(ctx context.Context, filter mysql.GatewayDeviceFilter) (size int64, err error) {
	//TODO implement me
	panic("implement me")
}
