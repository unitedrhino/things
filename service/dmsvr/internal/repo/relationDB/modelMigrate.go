package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/share/domain/protocols"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var needInitProductConfig bool
var NeedInitColumn bool

func Migrate(c conf.Database) error {
	if c.IsInitTable == false {
		return nil
	}
	ctx := ctxs.WithRoot(context.Background())
	db := stores.GetCommonConn(ctx)

	if !db.Migrator().HasTable(&DmProtocolInfo{}) {
		//需要初始化表
		NeedInitColumn = true
	} else if !db.Migrator().HasTable(&DmProductConfig{}) {
		needInitProductConfig = true
	}

	err := db.AutoMigrate(
		&DmProductConfig{},
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
		&DmProtocolScript{},
		&DmProtocolScriptDevice{},
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

	if NeedInitColumn {
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
				err = NewProductInfoRepo(tx).UpdateWithField(ctx, ProductFilter{ProtocolCode: "iThings"}, map[string]any{
					"protocol_code": protocols.ProtocolCodeUrMqtt,
				})
				return err
			})
		}
	}

	if needInitProductConfig { //1.3->1.4 升级
		err := func() error {
			pis, err := NewProductInfoRepo(ctx).FindByFilter(ctx, ProductFilter{}, nil)
			if err != nil {
				return err
			}
			var dbs []*DmProductConfig
			for _, p := range pis {
				dbs = append(dbs, &DmProductConfig{ProductID: p.ProductID})
			}
			return NewProductConfigRepo(ctx).MultiInsert(ctx, dbs)
		}()
		if err != nil {
			logx.Error(err)
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
	if err := db.CreateInBatches([]DmProductID{{ID: 100}}, 100).Error; err != nil {
		logx.Error(err)
	}
	if err := db.CreateInBatches([]DmGroupInfo{{ID: 1}, {ID: 2}, {ID: 3}}, 100).Error; err != nil {
		logx.Error(err)
	}

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
