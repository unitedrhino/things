package schemaDataRepo

import (
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"os"
)

const (
	PROPERTY_TYPE = "property_type"
)

type SchemaDataRepo struct {
	t              *clients.Td
	getSchemaModel schema.GetSchemaModel
	stores.SchemaStore
	kv kv.Store
}

func NewSchemaDataRepo(dataSource string, getSchemaModel schema.GetSchemaModel, kv kv.Store) *SchemaDataRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &SchemaDataRepo{t: td, getSchemaModel: getSchemaModel, kv: kv}
}
