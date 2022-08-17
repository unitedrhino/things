package schemaDataRepo

import (
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/store"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

const (
	PROPERTY_TYPE = "property_type"
)

type SchemaDataRepo struct {
	t              *clients.Td
	getSchemaModel schema.GetSchemaModel
	store.SchemaStore
}

func NewSchemaDataRepo(dataSource string, getSchemaModel schema.GetSchemaModel) *SchemaDataRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &SchemaDataRepo{t: td, getSchemaModel: getSchemaModel}
}
