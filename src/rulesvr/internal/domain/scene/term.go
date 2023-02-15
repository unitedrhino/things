// Package scene 触发条件
package scene

type Term struct {
	Column   string  `json:"column"`   //触发类型 device: 设备触发 timer: 定时触发 manual:手动触发
	Value    string  `json:"value"`    //条件值
	Type     string  `json:"type"`     //多个条件关联类型  or  and
	TermType string  `json:"termType"` //动态条件类型  eq: 相等  not:不相等  btw:在xx之间  gt: 大于  gte:大于等于 lt:小于  lte:小于等于   in:在xx值之间
	Terms    []*Term `json:"terms"`    //嵌套条件
}
