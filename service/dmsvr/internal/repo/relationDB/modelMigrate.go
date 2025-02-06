package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/share/domain/protocols"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Migrate(c conf.Database) error {
	if c.IsInitTable == false {
		return nil
	}
	db := stores.GetCommonConn(ctxs.WithRoot(context.Background()))
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
		&DmDeviceSchema{},
		//&DmProductSchema{},
		&DmSchemaInfo{},
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

	err = versionUpdate(db)
	if err != nil {
		return err
	}

	if needInitColumn {
		return migrateTableColumn()
	}
	return err
}

func versionUpdate(db *gorm.DB) error {
	ctx := ctxs.WithRoot(context.Background())
	//m := db.Migrator()
	{
		old, err := NewProtocolInfoRepo(ctx).FindOneByFilter(ctx, ProtocolInfoFilter{Code: "iThings"})
		if err == nil { //旧版的需要更新为新版
			db.Transaction(func(tx *gorm.DB) error {
				db := NewProtocolInfoRepo(tx)
				err = db.Delete(ctx, old.ID)
				if err != nil {
					return err
				}
				if err := tx.CreateInBatches(&MigrateProtocolInfo, 100).Error; err != nil {
					return err
				}
				return nil
			})
		}
	}

	return nil
}

func migrateTableColumn() error {
	db := stores.GetCommonConn(context.TODO()).Clauses(clause.OnConflict{DoNothing: true})
	if err := db.CreateInBatches(&MigrateProtocolInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateProductCategory, 100).Error; err != nil {
		return err
	}
	//if err := db.CreateInBatches(&MigrateManufacturerInfo, 100).Error; err != nil {
	//	return err
	//}
	if err := db.Create(&DmProductID{ID: 100}).Error; err != nil {
		return err
	}
	db.Create(&DmGroupInfo{ID: 10}) //分组前几个是特殊ID,不能使用,给他占位了
	db.Delete(&DmGroupInfo{ID: 10})

	return nil
}

var (
	MigrateProtocolInfo = []DmProtocolInfo{
		{Name: "联犀MQTT协议", Code: protocols.ProtocolCodeUrMqtt, Type: protocols.TypeNormal, TransProtocol: protocols.ProtocolMqtt, EtcdKey: "dg.rpc"},
		{Name: "联犀Http协议", Code: protocols.ProtocolCodeUrHttp, Type: protocols.TypeNormal, TransProtocol: protocols.ProtocolMqtt, EtcdKey: "dg.rpc"}}
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
