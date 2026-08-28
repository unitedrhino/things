package schema

// ShouldRecordHistory 判断指定物模型功能是否应写入历史记录。
// 模型或标识符缺失时保持既有行为，默认继续记录，避免因缓存或旧数据不完整误丢日志。
func (m *Model) ShouldRecordHistory(dataType AffordanceType, identifier string) bool {
	if m == nil {
		return true
	}

	switch dataType {
	case AffordanceTypeProperty:
		property := m.Property[identifier]
		return property == nil || property.RecordMode != RecordModeNone
	case AffordanceTypeEvent:
		event := m.Event[identifier]
		return event == nil || event.RecordMode != RecordModeNone
	case AffordanceTypeAction:
		action := m.Action[identifier]
		return action == nil || action.RecordMode != RecordModeNone
	default:
		return true
	}
}
