package deviceMsg

type TimeValue struct {
	Timestamp int64 `json:"timestamp,omitempty"` //毫秒时间戳
	Value     any   `json:"value"`               //值
}
type TimeParams struct {
	Timestamp int64          `json:"timestamp,omitempty"` //毫秒时间戳
	EventID   string         `json:"eventID,omitempty"`   //事件的 Id，在数据模板事件中定义。
	Params    map[string]any `json:"params,omitempty"`    //参数列表
}
