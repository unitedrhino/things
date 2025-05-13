package schemaDataRepo

import (
	"context"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/tdengine/schemaDataRepo"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	schema "gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/zeromicro/go-zero/core/stores/kv"
)

const (
	PropertyType = "property_type"
)

type DeviceDataRepo struct {
	//getProductSchemaModel schema.GetSchemaModel
	//getDeviceSchemaModel  schema.GetSchemaModel
	kv               kv.Store
	db               *stores.DB
	asyncEventInsert *stores.AsyncInsert[Event]
	asyncProperty    map[string]*stores.AsyncInsert[map[string]interface{}]

	asyncPropertyTimestamp      *stores.AsyncInsert[PropertyTimestamp]
	asyncPropertyTimestampArray *stores.AsyncInsert[PropertyTimestampArray]
	asyncPropertyBoolArray      *stores.AsyncInsert[PropertyBoolArray]
	asyncPropertyStringArray    *stores.AsyncInsert[PropertyStringArray]
	asyncPropertyIntArray       *stores.AsyncInsert[PropertyIntArray]
	asyncPropertyFloatArray     *stores.AsyncInsert[PropertyFloatArray]
	asyncPropertyBool           *stores.AsyncInsert[PropertyBool]
	asyncPropertyString         *stores.AsyncInsert[PropertyString]
	asyncPropertyInt            *stores.AsyncInsert[PropertyInt]
	asyncPropertyFloat          *stores.AsyncInsert[PropertyFloat]
	asyncPropertyStruct         *stores.AsyncInsert[PropertyStruct]
	asyncPropertyStructArray    *stores.AsyncInsert[PropertyStructArray]
	asyncPropertyEnum           *stores.AsyncInsert[PropertyEnum]
	asyncPropertyEnumArray      *stores.AsyncInsert[PropertyEnumArray]
}

func (d *DeviceDataRepo) VersionUpdate(ctx context.Context, version string, dc *caches.Cache[dm.DeviceInfo, devices.Core]) error {
	return nil
}

func (d *DeviceDataRepo) UpdateDevice(ctx context.Context, dev devices.Core, t *schema.Model, affiliation devices.Affiliation) error {
	return nil
}

func NewDeviceDataRepo(dataSource conf.TSDB, getProductSchemaModel schema.GetSchemaModel, getDeviceSchemaModel schema.GetSchemaModel, kv kv.Store) msgThing.SchemaDataRepo {
	if dataSource.DBType == conf.Tdengine {
		return schemaDataRepo.NewDeviceDataRepo(dataSource, getProductSchemaModel, getDeviceSchemaModel, kv)
	}
	stores.InitTsConn(dataSource)
	db := stores.GetTsConn(context.Background())
	asyncProperty := make(map[string]*stores.AsyncInsert[map[string]interface{}])
	for _, tb := range TableNames {
		asyncProperty[tb] = stores.NewAsyncInsert[map[string]interface{}](db, tb)
	}
	return &DeviceDataRepo{db: db, asyncProperty: asyncProperty,
		asyncEventInsert:            stores.NewAsyncInsert[Event](db, ""),
		asyncPropertyTimestamp:      stores.NewAsyncInsert[PropertyTimestamp](db, ""),
		asyncPropertyTimestampArray: stores.NewAsyncInsert[PropertyTimestampArray](db, ""),
		asyncPropertyBoolArray:      stores.NewAsyncInsert[PropertyBoolArray](db, ""),
		asyncPropertyStringArray:    stores.NewAsyncInsert[PropertyStringArray](db, ""),
		asyncPropertyIntArray:       stores.NewAsyncInsert[PropertyIntArray](db, ""),
		asyncPropertyFloatArray:     stores.NewAsyncInsert[PropertyFloatArray](db, ""),
		asyncPropertyBool:           stores.NewAsyncInsert[PropertyBool](db, ""),
		asyncPropertyString:         stores.NewAsyncInsert[PropertyString](db, ""),
		asyncPropertyInt:            stores.NewAsyncInsert[PropertyInt](db, ""),
		asyncPropertyFloat:          stores.NewAsyncInsert[PropertyFloat](db, ""),
		asyncPropertyStruct:         stores.NewAsyncInsert[PropertyStruct](db, ""),
		asyncPropertyStructArray:    stores.NewAsyncInsert[PropertyStructArray](db, ""),
		asyncPropertyEnum:           stores.NewAsyncInsert[PropertyEnum](db, ""),
		asyncPropertyEnumArray:      stores.NewAsyncInsert[PropertyEnumArray](db, ""),
		kv:                          kv}
}

func (d *DeviceDataRepo) Init(ctx context.Context) error {
	err := d.db.AutoMigrate(
		Event{},
		PropertyTimestamp{},
		PropertyTimestampArray{},
		PropertyBoolArray{},
		PropertyStringArray{},
		PropertyIntArray{},
		PropertyFloatArray{},
		PropertyBool{},
		PropertyString{},
		PropertyInt{},
		PropertyFloat{},
		PropertyStruct{},
		PropertyStructArray{},
	)
	if err != nil {
		return err
	}
	return nil
}
