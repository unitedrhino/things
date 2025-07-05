package schemaDataRepo

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceGroup"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/tdengine/schemaDataRepo"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	schema "gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"strings"
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

func NewDeviceDataRepo(dataSource conf.TSDB, getProductSchemaModel schema.GetSchemaModel, getDeviceSchemaModel schema.GetSchemaModel, kv kv.Store, g []*deviceGroup.GroupDetail) msgThing.SchemaDataRepo {
	if dataSource.DBType == conf.Tdengine {
		return schemaDataRepo.NewDeviceDataRepo(dataSource, getProductSchemaModel, getDeviceSchemaModel, kv, g)
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
	var NeedInitColumn bool
	if !d.db.Migrator().HasTable(&Event{}) {
		//需要初始化表
		NeedInitColumn = true
	}
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
		PropertyEnum{},
		PropertyEnumArray{},
	)
	if err != nil {
		return err
	}
	if NeedInitColumn && stores.GetTsDBType() == conf.Pgsql {
		d.db.Exec("SELECT create_hypertable('dm_time_model_event','ts', chunk_time_interval => interval '1 day'    );")
		for _, tb := range TableNames {
			d.db.Exec(fmt.Sprintf("SELECT create_hypertable('%s','ts', chunk_time_interval => interval '1 day');", tb))
			d.db.Exec(fmt.Sprintf("SELECT add_dimension('%s', by_hash('device_name', 2));", tb))
			d.db.Exec(fmt.Sprintf("SELECT add_dimension('%s', by_hash('identifier', 2));", tb))
			d.db.Exec(fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS idx_unique ON %s( product_id, device_name, identifier,ts);", tb))
			if strings.HasSuffix(tb, "bool") {
				d.db.Exec(fmt.Sprintf(viewBoolTemplate, tb, "day", "day", tb))
				d.db.Exec(fmt.Sprintf(viewBoolTemplate, tb, "hour", "hour", tb))
			} else {
				d.db.Exec(fmt.Sprintf(viewTemplate, tb, "day", "day", tb))
				d.db.Exec(fmt.Sprintf(viewTemplate, tb, "hour", "hour", tb))
			}
		}
	}

	return nil
}

const (
	viewTemplate = `CREATE MATERIALIZED VIEW if not exists %s_%s(product_id,device_name,identifier,ts,first_ts,first_param,last_ts,last_param, max_ts,max_param,min_ts,min_param, sum_param, count_param,avg_param )
			WITH (timescaledb.continuous) AS
			SELECT product_id,device_name,identifier,time_bucket('1%s', ts) as ts_window,
				(ARRAY_AGG(ts ORDER BY ts ASC))[1]    AS first_ts,
				(ARRAY_AGG(param ORDER BY ts ASC))[1] AS first_param,
				(ARRAY_AGG(ts ORDER BY ts desc))[1]    AS last_ts,
				(ARRAY_AGG(param ORDER BY ts desc))[1] AS last_param,
       (ARRAY_AGG(ts ORDER BY param desc))[1]    AS max_ts,
       (ARRAY_AGG(param ORDER BY param desc))[1] AS max_param,
       (ARRAY_AGG(ts ORDER BY param ASC))[1]    AS min_ts,
       (ARRAY_AGG(param ORDER BY param ASC))[1] AS min_param, sum(param),count(param),avg(param)
			FROM %s
			GROUP BY product_id,device_name,identifier,ts_window;`
	viewBoolTemplate = `CREATE MATERIALIZED VIEW if not exists %s_%s(product_id,device_name,identifier,ts,first_ts,first_param,last_ts,last_param, count_param,avg_param )
			WITH (timescaledb.continuous) AS
			SELECT product_id,device_name,identifier,time_bucket('1%s', ts) as ts_window,
				(ARRAY_AGG(ts ORDER BY ts ASC))[1]    AS first_ts,
				(ARRAY_AGG(param ORDER BY ts ASC))[1] AS first_param,
				(ARRAY_AGG(ts ORDER BY ts desc))[1]    AS last_ts,
				(ARRAY_AGG(param ORDER BY ts desc))[1] AS last_param,count(param),avg(param)
			FROM %s
			GROUP BY product_id,device_name,identifier,ts_window;`
)
