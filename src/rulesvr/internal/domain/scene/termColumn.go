package scene

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
)

type TermColumnType string

const (
	TermColumnTypeProperty   TermColumnType = "property"
	TermColumnTypeEvent      TermColumnType = "event"
	TermColumnTypeReportTime TermColumnType = "reportTime"
	TermColumnTypeSysTime    TermColumnType = "sysTime"
)

type TermColumn struct {
	Type   TermColumnType `json:"type"`   //字段类型 property:属性 event:事件 reportTime:上报时间 sysTime:系统时间
	DataID []string       `json:"dataID"` //属性的id及事件的id
}

func (t *TermColumn) Validate() error {
	if t == nil {
		return nil
	}
	if !utils.SliceIn(t.Type, TermColumnTypeProperty, TermColumnTypeEvent, TermColumnTypeReportTime, TermColumnTypeSysTime) {
		return errors.Parameter.AddMsg("触发条件中的字段名类型不支持:" + string(t.Type))
	}
	//todo 这里需要校验属性是否存在
	return nil
}
