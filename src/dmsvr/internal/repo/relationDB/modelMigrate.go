package relationDB

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/stores"
)

func Migrate(c conf.Database) error {
	if c.IsInitTable == false {
		return nil
	}
	db := stores.GetCommonConn(nil)
	return db.AutoMigrate(
		&DmProductInfo{},
		&DmDeviceInfo{},
		&DmProductCustom{},
		&DmProductSchema{},
		&DmGroupInfo{},
		&DmGroupDevice{},
		&DmGatewayDevice{},
		&DmProductRemoteConfig{},
		&DmOtaFirmware{},
		&DmOtaTask{},
		&DmOtaFirmwareFile{},
		&DmOtaTaskDevices{},
		&DmDeviceShadow{},
	)
}
