// 本文件验证按属性清理 TDengine 数据时的子表展开规则。
package schemaDataRepo

import (
	"reflect"
	"testing"

	"gitee.com/unitedrhino/things/share/domain/schema"
)

func TestGetPropertyTableNamesForCleanup(t *testing.T) {
	store := &SchemaStore{}
	tests := []struct {
		name     string
		property *schema.Property
		want     []string
	}{
		{
			name: "scalar struct",
			property: &schema.Property{
				CommonParam: schema.CommonParam{Identifier: "hardware"},
				Define:      schema.Define{Type: schema.DataTypeStruct},
			},
			want: []string{"`device_property_product_device_hardware`"},
		},
		{
			name: "array",
			property: &schema.Property{
				CommonParam: schema.CommonParam{Identifier: "voltage"},
				Define: schema.Define{
					Type:      schema.DataTypeArray,
					Max:       "2",
					ArrayInfo: &schema.Define{Type: schema.DataTypeFloat},
				},
			},
			want: []string{
				"`device_property_product_device_voltage_0`",
				"`device_property_product_device_voltage_1`",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := store.GetPropertyTableNames("product", "device", test.property)
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("GetPropertyTableNames() = %v, want %v", got, test.want)
			}
		})
	}
}
