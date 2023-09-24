package relationDB

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/stores"
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
		SubDevice *devices.Core
	}
)

func NewGatewayDeviceRepo(in any) *GatewayDeviceRepo {
	return &GatewayDeviceRepo{db: stores.GetCommonConn(in)}
}
func (p GatewayDeviceRepo) fmtFilter(ctx context.Context, f GatewayDeviceFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	di := DmDeviceInfo{}
	gd := DmGatewayDevice{}
	if f.Gateway != nil { //通过网关获取旗下子设备列表
		db = db.Table(gd.TableName()+" as gd").Joins(fmt.Sprintf(
			"left join %s as di on di.product_id=gd.product_id and di.device_name=gd.device_name", di.TableName())).
			Where("gateway_product_id=? and gateway_device_name=? and di.id IS NOT NULL", f.Gateway.ProductID, f.Gateway.DeviceName)
	} else {
		db = db.Table(gd.TableName()+" as gd").Joins(fmt.Sprintf(
			"left join %s as di on di.product_id=gd.gateway_product_id and di.device_name=gd.gateway_device_name", di.TableName())).
			Where("gd.product_id=? and gd.device_name=? and di.id IS NOT NULL", f.SubDevice.ProductID, f.SubDevice.DeviceName)
	}
	return db
}

func (g GatewayDeviceRepo) FindByFilter(ctx context.Context, f GatewayDeviceFilter, page *def.PageInfo) ([]*DmDeviceInfo, error) {
	var results []*DmDeviceInfo
	db := g.fmtFilter(ctx, f)
	db = page.ToGorm(db)
	err := db.Select("di.*").Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (g GatewayDeviceRepo) CountByFilter(ctx context.Context, f GatewayDeviceFilter) (size int64, err error) {
	db := g.fmtFilter(ctx, f)
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}
func (p GatewayDeviceRepo) FindOneByFilter(ctx context.Context, f GatewayDeviceFilter) (*DmDeviceInfo, error) {
	var result DmDeviceInfo
	db := p.fmtFilter(ctx, f)
	err := db.Select("di.*").First(&result).Error
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
	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmGatewayDevice{}).Create(data).Error
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
	db = db.Where("gateway_product_id=? and gateway_device_name`=?", gateway.ProductID, gateway.DeviceName).Where(scope(db))
	err := db.Delete(&DmGatewayDevice{}).Error
	return stores.ErrFmt(err)
}
