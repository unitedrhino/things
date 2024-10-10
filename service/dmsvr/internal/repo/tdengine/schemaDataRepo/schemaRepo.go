package schemaDataRepo

import (
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/conf"
	schema "gitee.com/unitedrhino/share/domain/schema"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"os"
)

const (
	PropertyType = "property_type"
)

type DeviceDataRepo struct {
	t              *clients.Td
	getSchemaModel schema.GetSchemaModel
	SchemaStore
	kv kv.Store
}

func NewDeviceDataRepo(dataSource conf.TSDB, getSchemaModel schema.GetSchemaModel, kv kv.Store) *DeviceDataRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &DeviceDataRepo{t: td, getSchemaModel: getSchemaModel, kv: kv}
}
