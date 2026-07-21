package deviceinteractlogic

import (
	"testing"

	"gitee.com/unitedrhino/things/share/domain/schema"
)

func TestFillLegacyEmptySec64(t *testing.T) {
	tests := []struct {
		name       string
		actionID   string
		model      *schema.Model
		params     map[string]any
		wantSec64  string
		wantExists bool
	}{
		{
			name:     "物模型定义sec64且请求缺失时补空值",
			actionID: "RF433Ctrl",
			model: &schema.Model{Action: schema.ActionMap{
				"RF433Ctrl": {In: map[string]*schema.Param{
					"sec64": {Identifier: "sec64"},
				}},
			}},
			params:     map[string]any{"keyID": int64(1)},
			wantSec64:  "",
			wantExists: true,
		},
		{
			name:     "请求已有sec64时保留原值",
			actionID: "RF433Ctrl",
			model: &schema.Model{Action: schema.ActionMap{
				"RF433Ctrl": {In: map[string]*schema.Param{
					"sec64": {Identifier: "sec64"},
				}},
			}},
			params:     map[string]any{"sec64": "signed-value"},
			wantSec64:  "signed-value",
			wantExists: true,
		},
		{
			name:     "动作未定义sec64时不添加字段",
			actionID: "OtherAction",
			model: &schema.Model{Action: schema.ActionMap{
				"OtherAction": {In: map[string]*schema.Param{}},
			}},
			params:     map[string]any{"value": int64(1)},
			wantExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fillLegacyEmptySec64(tt.model, tt.actionID, tt.params)
			got, exists := tt.params["sec64"]
			if exists != tt.wantExists {
				t.Fatalf("sec64 exists = %v, want %v", exists, tt.wantExists)
			}
			if !exists {
				return
			}
			gotSec64, ok := got.(string)
			if !ok {
				t.Fatalf("sec64 type = %T, want string", got)
			}
			if gotSec64 != tt.wantSec64 {
				t.Fatalf("sec64 = %q, want %q", gotSec64, tt.wantSec64)
			}
		})
	}
}
