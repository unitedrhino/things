package relationDB

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GatewayDeviceRepo struct {
	db *gorm.DB
}

type (
	GatewayDeviceFilter struct {
		//网关和子设备 至少要有一个填写
		Gateway *devices.Core
		//网关和子设备 至少要有一个填写
		SubDevice  *devices.Core
		SubDevices []*devices.Core
	}
)

func NewGatewayDeviceRepo(in any) *GatewayDeviceRepo {
	return &GatewayDeviceRepo{db: stores.GetCommonConn(in)}
}
func (p GatewayDeviceRepo) fmtFilter(ctx context.Context, f GatewayDeviceFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Gateway != nil { //通过网关获取旗下子设备列表
		db = db.Where("gateway_product_id=? and gateway_device_name=?", f.Gateway.ProductID, f.Gateway.DeviceName)
	}
	if f.SubDevice != nil { //根据子设备获取网关
		db = db.Where("product_id=? and device_name=?", f.SubDevice.ProductID, f.SubDevice.DeviceName)
	}
	if len(f.SubDevices) != 0 {
		db = db.Where(fmt.Sprintf("(product_id, device_name) in (%s)", utils.JoinWithFunc(f.SubDevices, ",", func(in *devices.Core) string {
			return fmt.Sprintf("('%s','%s')", in.ProductID, in.DeviceName)
		})))
	}
	return db
}

func (g GatewayDeviceRepo) FindByFilter(ctx context.Context, f GatewayDeviceFilter, page *stores.PageInfo) ([]*DmGatewayDevice, error) {
	var results []*DmGatewayDevice
	db := g.fmtFilter(ctx, f)
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (g GatewayDeviceRepo) CountByFilter(ctx context.Context, f GatewayDeviceFilter) (size int64, err error) {
	db := g.fmtFilter(ctx, f)
	err = db.Model(&DmGatewayDevice{}).Count(&size).Error
	return size, stores.ErrFmt(err)
}
func (p GatewayDeviceRepo) FindOneByFilter(ctx context.Context, f GatewayDeviceFilter) (*DmGatewayDevice, error) {
	var result DmGatewayDevice
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (m GatewayDeviceRepo) MultiInsert(ctx context.Context, gateway *devices.Core, subDevice []*devices.Core) error {
	var data []*DmGatewayDevice
	for _, v := range subDevice {
		data = append(data, &DmGatewayDevice{
			GatewayProductID:  gateway.ProductID,
			GatewayDeviceName: gateway.DeviceName,
			ProductID:         v.ProductID,
			DeviceName:        v.DeviceName,
		})
	}
	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Model(&DmGatewayDevice{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (m GatewayDeviceRepo) DeleteDevAll(ctx context.Context, dev devices.Core) error {
	db := m.db.WithContext(ctx).Model(&DmGatewayDevice{})
	db = db.Where("(gateway_product_id=? and gateway_device_name=?) or (product_id = ? and device_name = ?)",
		dev.ProductID, dev.DeviceName, dev.ProductID, dev.DeviceName)
	err := db.Delete(&DmGatewayDevice{}).Error
	return stores.ErrFmt(err)
}

// 批量插入 LightStrategyDevice 记录
func (m GatewayDeviceRepo) MultiDelete(ctx context.Context, gateway *devices.Core, subDevice []*devices.Core) error {
	if len(subDevice) < 1 {
		return nil
	}
	scope := func(db *gorm.DB) *gorm.DB {
		for i, d := range subDevice {
			if i == 0 {
				db = db.Where("product_id = ? and device_name = ?", d.ProductID, d.DeviceName)
				continue
			}
			db = db.Or("product_id = ? and device_name = ?", d.ProductID, d.DeviceName)
		}
		return db
	}
	db := m.db.WithContext(ctx).Model(&DmGatewayDevice{})
	db = db.Where("gateway_product_id=? and gateway_device_name=?", gateway.ProductID, gateway.DeviceName).Where(scope(db))
	err := db.Delete(&DmGatewayDevice{}).Error
	return stores.ErrFmt(err)
}
