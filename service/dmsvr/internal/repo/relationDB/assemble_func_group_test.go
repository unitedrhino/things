// 本文件验证物模型属性在数据库结构与领域结构互转时保留功能分类。
package relationDB

import (
	"testing"

	"gitee.com/unitedrhino/things/share/domain/schema"
)

func TestToPropertyPoPreservesFuncGroup(t *testing.T) {
	property := &schema.Property{
		CommonParam: schema.CommonParam{
			Identifier: "HWInfo",
			Name:       "硬件信息",
			Tag:        schema.TagOptional,
			FuncGroup:  schema.FuncGroupSystem,
		},
		Mode:   schema.PropertyModeR,
		Define: schema.Define{Type: schema.DataTypeString},
	}

	core := ToPropertyPo(property)
	if core.FuncGroup != schema.FuncGroupSystem {
		t.Fatalf("ToPropertyPo() FuncGroup = %d, want %d", core.FuncGroup, schema.FuncGroupSystem)
	}
}

func TestToPropertyDoPreservesFuncGroup(t *testing.T) {
	core := &DmSchemaCore{
		Type:      schema.AffordanceTypeProperty,
		Tag:       schema.TagOptional,
		Name:      "硬件信息",
		Required:  2,
		FuncGroup: schema.FuncGroupSystem,
		Affordance: `{
			"mode":"r",
			"define":{"type":"string"}
		}`,
	}

	property := ToPropertyDo("HWInfo", core)
	if property.FuncGroup != schema.FuncGroupSystem {
		t.Fatalf("ToPropertyDo() FuncGroup = %d, want %d", property.FuncGroup, schema.FuncGroupSystem)
	}
}
