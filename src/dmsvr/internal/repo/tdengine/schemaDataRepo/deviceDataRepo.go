package schemaDataRepo

import (
	"github.com/i-Things/things/shared/clients"
	schema "github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/stores"
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
	stores.SchemaStore
	kv kv.Store
}

func NewDeviceDataRepo(dataSource string, getSchemaModel schema.GetSchemaModel, kv kv.Store) *DeviceDataRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &DeviceDataRepo{t: td, getSchemaModel: getSchemaModel, kv: kv}
}
