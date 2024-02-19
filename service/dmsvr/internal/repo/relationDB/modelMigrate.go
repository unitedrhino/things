package relationDB

import (
	"context"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm/clause"
)

func Migrate(c conf.Database) error {
	if c.IsInitTable == false {
		return nil
	}
	db := stores.GetCommonConn(context.TODO())
	var needInitColumn bool
	if !db.Migrator().HasTable(&DmProtocolInfo{}) {
		//需要初始化表
		needInitColumn = true
	}
	err := db.AutoMigrate(
		&DmProductInfo{},
		&DmProductCategory{},
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
		&DmOtaJob{},
		&DmOtaUpgradeTask{},
		&DmOtaModule{},
		&DmUserDeviceCollect{},
		&DmUserDeviceShare{},
	)
	if err != nil {
		return err
	}
	//stores.SetAuthIncrement(db, &DmGroupInfo{ID: 10})
	if needInitColumn {
		return migrateTableColumn()
	}
	return err
}
func migrateTableColumn() error {
	db := stores.GetCommonConn(context.TODO()).Clauses(clause.OnConflict{DoNothing: true})
	if err := db.CreateInBatches(&MigrateProtocolInfo, 100).Error; err != nil {
		return err
	}
	return nil
}

var (
	MigrateProtocolInfo = []DmProtocolInfo{{ID: 2, Name: "thingsboard协议", Protocol: "mqtt", ProtocolType: "thingsboard", EtcdKey: "dg.thingsboard.rpc"}}
)
