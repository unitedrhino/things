package schemaDataRepo

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/conf"
	schema "gitee.com/unitedrhino/share/domain/schema"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"os"
)

const (
	PropertyType = "property_type"
)

type DeviceDataRepo struct {
	t                     *clients.Td
	getProductSchemaModel schema.GetSchemaModel
	getDeviceSchemaModel  schema.GetSchemaModel
	SchemaStore
	kv kv.Store
}

func NewDeviceDataRepo(dataSource conf.TSDB, getProductSchemaModel schema.GetSchemaModel, getDeviceSchemaModel schema.GetSchemaModel, kv kv.Store) *DeviceDataRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &DeviceDataRepo{t: td, getProductSchemaModel: getProductSchemaModel, getDeviceSchemaModel: getDeviceSchemaModel, kv: kv}
}

func (d *DeviceDataRepo) Init(ctx context.Context) error {
	{
		sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s "+
			"(`ts` timestamp,`event_id` BINARY(50),`event_type` BINARY(20), `param` BINARY(5000)) "+
			"TAGS (`product_id` BINARY(50),`device_name` BINARY(50));",
			d.GetEventStableName())
		if _, err := d.t.ExecContext(ctx, sql); err != nil {
			return errors.Database.AddDetail(err)
		}
	}

	genDeviceStable := func(tb string, def schema.Define) error {
		sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s (`ts` timestamp,`param` %s)"+
			" TAGS (`product_id` BINARY(50),`device_name` BINARY(50),`"+PropertyType+"` BINARY(50));",
			tb, stores.GetTdType(def))
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
