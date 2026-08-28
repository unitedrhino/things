package schema

import "testing"

func TestModelShouldRecordHistory(t *testing.T) {
	tests := []struct {
		name       string
		model      *Model
		dataType   AffordanceType
		identifier string
		want       bool
	}{
		{
			name:       "模型缺失时保持历史兼容并记录",
			model:      nil,
			dataType:   AffordanceTypeAction,
			identifier: "ActionA",
			want:       true,
		},
		{
			name: "动作配置不记录历史时跳过",
			model: &Model{Action: ActionMap{
				"ActionA": {CommonParam: CommonParam{RecordMode: RecordModeNone}},
			}},
			dataType:   AffordanceTypeAction,
			identifier: "ActionA",
			want:       false,
		},
		{
			name: "事件配置不记录历史时跳过",
			model: &Model{Event: EventMap{
				"EventA": {CommonParam: CommonParam{RecordMode: RecordModeNone}},
			}},
			dataType:   AffordanceTypeEvent,
			identifier: "EventA",
			want:       false,
		},
		{
			name: "默认记录模式继续记录",
			model: &Model{Event: EventMap{
				"EventA": {CommonParam: CommonParam{}},
			}},
			dataType:   AffordanceTypeEvent,
			identifier: "EventA",
			want:       true,
		},
		{
			name:       "未知标识符保持历史兼容并记录",
			model:      &Model{Action: ActionMap{}},
			dataType:   AffordanceTypeAction,
			identifier: "MissingAction",
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.model.ShouldRecordHistory(tt.dataType, tt.identifier); got != tt.want {
				t.Fatalf("ShouldRecordHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}
