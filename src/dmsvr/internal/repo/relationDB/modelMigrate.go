package relationDB

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm/clause"
)

func Migrate(c conf.Database) error {
	if c.IsInitTable == false {
		return nil
	}
	db := stores.GetCommonConn(nil)
	var needInitColumn bool
	if !db.Migrator().HasTable(&DmProtocolInfo{}) {
		//需要初始化表
		needInitColumn = true
	}
	err := db.AutoMigrate(
		&DmProductInfo{},
		&DmProtocolInfo{},
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
	if err != nil {
		return err
	}
	if needInitColumn {
		return migrateTableColumn()
	}
	return err
}
func migrateTableColumn() error {
	db := stores.GetCommonConn(nil).Clauses(clause.OnConflict{DoNothing: true})
	if err := db.CreateInBatches(&MigrateProtocolInfo, 100).Error; err != nil {
		return err
	}
	return nil
}

var (
	MigrateProtocolInfo = []DmProtocolInfo{{ID: 2, Name: "thingsboard协议", Protocol: "mqtt", ProtocolType: "thingsboard", EtcdKey: "dg.thingsboard.rpc"}}
)
