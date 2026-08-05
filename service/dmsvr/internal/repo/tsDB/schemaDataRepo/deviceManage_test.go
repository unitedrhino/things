// 本文件验证按属性清理普通时序库数据时的标识符展开规则。
package schemaDataRepo

import (
	"reflect"
	"testing"

	"gitee.com/unitedrhino/things/share/domain/schema"
)

func TestPropertyDBIdentifiers(t *testing.T) {
	tests := []struct {
		name     string
		property schema.Property
		want     []string
	}{
		{
			name: "scalar",
			property: schema.Property{
				CommonParam: schema.CommonParam{Identifier: "temperature"},
				Define:      schema.Define{Type: schema.DataTypeFloat},
			},
			want: []string{"temperature"},
		},
		{
			name: "array",
			property: schema.Property{
				CommonParam: schema.CommonParam{Identifier: "voltage"},
				Define: schema.Define{
					Type:      schema.DataTypeArray,
					Max:       "3",
					ArrayInfo: &schema.Define{Type: schema.DataTypeFloat},
				},
			},
			want: []string{"voltage_0", "voltage_1", "voltage_2"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := propertyDBIdentifiers(test.property); !reflect.DeepEqual(got, test.want) {
				t.Fatalf("propertyDBIdentifiers() = %v, want %v", got, test.want)
			}
		})
	}
}
