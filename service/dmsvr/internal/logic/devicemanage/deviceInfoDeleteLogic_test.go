// 本文件验证设备转让和解绑清理数据时的系统物模型保留规则。
package devicemanagelogic

import (
	"testing"

	"gitee.com/unitedrhino/things/share/domain/schema"
)

func TestFilterNonSystemProperties(t *testing.T) {
	model := &schema.Model{
		Property: schema.PropertyMap{
			"normal": {
				CommonParam: schema.CommonParam{Identifier: "normal", FuncGroup: schema.FuncGroupNormal},
			},
			"system": {
				CommonParam: schema.CommonParam{Identifier: "system", FuncGroup: schema.FuncGroupSystem},
			},
			"legacy": {
				CommonParam: schema.CommonParam{Identifier: "legacy"},
			},
			"nil": nil,
		},
	}

	got := filterNonSystemProperties(model)
	if len(got) != 2 {
		t.Fatalf("filterNonSystemProperties() length = %d, want 2", len(got))
	}

	identifiers := map[string]bool{}
	for _, property := range got {
		identifiers[property.Identifier] = true
		if property.FuncGroup == schema.FuncGroupSystem {
			t.Fatalf("filterNonSystemProperties() returned system property %q", property.Identifier)
		}
	}
	if !identifiers["normal"] || !identifiers["legacy"] {
		t.Fatalf("filterNonSystemProperties() identifiers = %v, want normal and legacy", identifiers)
	}
}

func TestFilterNonSystemPropertiesNilModel(t *testing.T) {
	if got := filterNonSystemProperties(nil); got != nil {
		t.Fatalf("filterNonSystemProperties(nil) = %v, want nil", got)
	}
}
