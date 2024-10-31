package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
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
		&DmDeviceMsgCount{},
		&DmManufacturerInfo{},
		&DmProtocolService{},
		&DmOtaModuleInfo{},
		&DmProductInfo{},
		&DmProductCategory{},
		&DmProductCategorySchema{},
		&DmProtocolInfo{},
		&DmDeviceInfo{},
		&DmProductCustom{},
		&DmSchemaInfo{},
		&DmProductSchema{},
		&DmGroupInfo{},
		&DmGroupDevice{},
		&DmCommonSchema{},
		&DmGatewayDevice{},
		&DmProductRemoteConfig{},
		&DmOtaFirmwareInfo{},
		&DmOtaFirmwareFile{},
		&DmDeviceShadow{},
		&DmDeviceProfile{},
		&DmOtaFirmwareJob{},
		&DmOtaFirmwareDevice{},
		&DmUserDeviceCollect{},
		&DmUserDeviceShare{},
		&DmProductID{},
	)
	if err != nil {
		return err
	}

	{ //版本升级兼容
		m := db.Migrator()
		if m.HasColumn(&DmGatewayDevice{}, "tenant_code") {
			err := db.Migrator().DropColumn(&DmGatewayDevice{}, "tenant_code")
			if err != nil {
				return err
			}
		}

	}
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
	if err := db.CreateInBatches(&MigrateProductCategory, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateManufacturerInfo, 100).Error; err != nil {
		return err
	}
	if err := db.Create(&DmProductID{ID: 100}).Error; err != nil {
		return err
	}

	return nil
}

var (
	MigrateProtocolInfo    = []DmProtocolInfo{{ID: 3, Name: "iThings标准协议", Code: "iThings", TransProtocol: "mqtt", EtcdKey: "dg.rpc"}}
	MigrateProductCategory = []DmProductCategory{
		{ID: 3, Name: "照明设备", ParentID: def.RootNode, IDPath: "3-"},
		{ID: 4, Name: "空调设备", ParentID: def.RootNode, IDPath: "4-"},
		{ID: 5, Name: "风扇设备", ParentID: def.RootNode, IDPath: "5-"},
		{ID: 6, Name: "传感器设备", ParentID: def.RootNode, IDPath: "6-"},
	}
	MigrateManufacturerInfo = []DmManufacturerInfo{
		{
			ID:    1,
			Name:  "iThings",
			Desc:  "欢迎加入",
			Phone: "166666666",
		},
	}
)
