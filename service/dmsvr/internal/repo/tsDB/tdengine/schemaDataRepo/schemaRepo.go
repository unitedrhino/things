package schemaDataRepo

import (
	"context"
	"fmt"
	"os"

	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceGroup"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/cache"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/tdengine"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	schema "gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
)

const (
	PropertyType = "property_type"
)

type DeviceDataRepo struct {
	t                     *clients.Td
	getProductSchemaModel schema.GetSchemaModel
	getDeviceSchemaModel  schema.GetSchemaModel
	SchemaStore
	kv           kv.Store
	groupConfigs []*deviceGroup.GroupDetail
	cacheManager *cache.PropertyCacheManager
}

func NewDeviceDataRepo(dataSource conf.TSDB, getProductSchemaModel schema.GetSchemaModel,
	getDeviceSchemaModel schema.GetSchemaModel, kv kv.Store, g []*deviceGroup.GroupDetail) msgThing.SchemaDataRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &DeviceDataRepo{t: td, getProductSchemaModel: getProductSchemaModel,
		getDeviceSchemaModel: getDeviceSchemaModel, kv: kv, groupConfigs: g,
		cacheManager: cache.NewPropertyCacheManager(kv)}
}

func (d *DeviceDataRepo) VersionUpdate(ctx context.Context, version string, dc *caches.Cache[dm.DeviceInfo, devices.Core]) error {
	desc, err := d.t.DescTable(ctx, "model_device_property_bool")
	if err != nil {
		return err
	}
	if desc["group_ids"] != nil {
		tbs, err := d.t.STables(ctx, "_property_")
		if err != nil {
			return err
		}
		for _, tb := range tbs {
			d.t.ExecContext(ctx, fmt.Sprintf("ALTER STABLE `%s` DROP TAG `group_ids` ;", tb))
			d.t.ExecContext(ctx, fmt.Sprintf("ALTER STABLE `%s` DROP TAG `group_id_paths`;", tb))
			for _, g := range d.groupConfigs {
				_, err = d.t.ExecContext(ctx, fmt.Sprintf("ALTER STABLE `%s` ADD TAG `group_%s_ids`  BINARY(250) ;", tb, g.Value))
				if err != nil {
					continue
				}
				_, err = d.t.ExecContext(ctx, fmt.Sprintf("ALTER STABLE `%s` ADD TAG `group_%s_id_paths`  BINARY(250) ;", tb, g.Value))
				if err != nil {
					continue
				}
			}
		}
	}
	if desc["tenant_code"] != nil {
		return nil
	}
	{
		tbs, err := d.t.STables(ctx, "_property_")
		if err != nil {
			return err
		}
		for _, tb := range tbs {
			_, err = d.t.ExecContext(ctx, fmt.Sprintf("ALTER STABLE `%s` ADD TAG `tenant_code`  BINARY(50) ;", tb))
			if err != nil {
				continue
			}
			_, err = d.t.ExecContext(ctx, fmt.Sprintf("ALTER STABLE `%s` ADD TAG  `project_id` BIGINT ;", tb))
			if err != nil {
				continue
			}
			_, err = d.t.ExecContext(ctx, fmt.Sprintf("ALTER STABLE `%s` ADD TAG  `area_id` BIGINT  ;", tb))
			if err != nil {
				continue
			}
			_, err = d.t.ExecContext(ctx, fmt.Sprintf("ALTER STABLE `%s` ADD TAG `area_id_path`  BINARY(50) ;", tb))
			if err != nil {
				continue
			}
			for _, g := range d.groupConfigs {
				_, err = d.t.ExecContext(ctx, fmt.Sprintf("ALTER STABLE `%s` ADD TAG `group_%s_ids`  BINARY(250) ;", tb, g.Value))
				if err != nil {
					continue
				}
				_, err = d.t.ExecContext(ctx, fmt.Sprintf("ALTER STABLE `%s` ADD TAG `group_%s_id_paths`  BINARY(250) ;", tb, g.Value))
				if err != nil {
					continue
				}
			}
		}
	}
	{
		tbs, err := d.t.Tables(ctx, "device_property_")
		if err != nil {
			return err
		}
		for _, tb := range tbs {
			ts, err := d.t.TableTags(ctx, tb)
			if err != nil {
				logx.WithContext(ctx).Errorf("get tags err: %v", err)
				continue
			}
			dev := devices.Core{ProductID: ts["product_id"], DeviceName: ts["device_name"]}
			di, err := dc.GetData(ctx, dev)
			if err != nil {
				logx.WithContext(ctx).Error(err.Error())
				continue
			}
			err = tdengine.AlterTag(ctx, d.t, []string{tb}, tdengine.AffiliationToMap(devices.Affiliation{
				TenantCode:  di.TenantCode,
				ProjectID:   di.ProjectID,
				AreaID:      di.AreaID,
				AreaIDPath:  di.AreaIDPath,
				BelongGroup: tdengine.ToBelongGroup(di.BelongGroup),
			}, d.groupConfigs))
			if err != nil {
				logx.WithContext(ctx).Error(err.Error())
				continue
			}
		}
	}
	return nil
}

func (d *DeviceDataRepo) Init(ctx context.Context) error {
	{
		tags := tdengine.GenTagsDef(defaultTagDef, d.groupConfigs)
		sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s "+
			"(`ts` timestamp,`event_id` BINARY(50),`event_type` BINARY(20), `param` BINARY(5000)) "+
			"TAGS (%s);",
			d.GetEventStableName(), tags)
		if _, err := d.t.ExecContext(ctx, sql); err != nil {
			return errors.Database.AddDetail(err)
		}
	}
	ts := "`product_id` BINARY(50),`device_name` BINARY(50),`" + PropertyType + "` BINARY(50)," +
		" `tenant_code`  BINARY(50),`project_id` BIGINT,`area_id` BIGINT,`area_id_path`  BINARY(50)"
	tags := tdengine.GenTagsDef(ts, d.groupConfigs)
	genDeviceStable := func(tb string, def schema.Define) error {
		sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s (`ts` timestamp,`param` %s)"+
			" TAGS (%s);",
			tb, tdengine.GetTdType(def), tags)
		if _, err := d.t.ExecContext(ctx, sql); err != nil {
			return errors.Database.AddDetail(err)
		}
		return nil
	}
	execDeviceTables := []struct {
		tb  string
		def schema.Define
	}{
		{tb: d.GetDeviceStableBoolName(), def: schema.Define{Type: schema.DataTypeBool}},
		{tb: d.GetDeviceStableIntName(), def: schema.Define{Type: schema.DataTypeInt}},
		{tb: d.GetDeviceStableStringName(), def: schema.Define{Type: schema.DataTypeString}},
		{tb: d.GetDeviceStableTimestampName(), def: schema.Define{Type: schema.DataTypeTimestamp}},
		{tb: d.GetDeviceStableEnumName(), def: schema.Define{Type: schema.DataTypeEnum}},
		{tb: d.GetDeviceStableFloatName(), def: schema.Define{Type: schema.DataTypeFloat}},
	}
	for _, e := range execDeviceTables {
		err := genDeviceStable(e.tb, e.def)
		if err != nil {
			return err
		}
	}
	return nil
}
