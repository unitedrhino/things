package scene

import (
	"gitee.com/i-Things/share/def"
	"time"
)

type Log struct {
	AreaID      int64        `json:"areaID,string"`
	SceneID     int64        `json:"sceneID"`
	Type        SceneType    `json:"type"`
	Status      def.Bool     `json:"status"`
	CreatedTime time.Time    `json:"createdTime"`
	Trigger     *LogTrigger  `json:"trigger,omitempty"`
	Actions     []*LogAction `json:"actions"`
}
type LogAction struct {
	Type     ActionType       `json:"type"` //执行器类型 notify: 通知 delay:延迟  device:设备输出  alarm: 告警
	Device   *LogActionDevice `json:"device"`
	Alarm    *LogActionAlarm  `json:"alarm"`
	Status   int64            `json:"status"`
	Code     int64            `json:"code"`
	Msg      string           `json:"msg"`
	MsgToken string           `json:"msgToken"`
}

type LogActionAlarm struct {
	Mode ActionAlarmMode `json:"mode"` //告警模式  trigger: 触发告警  relieve: 解除告警
}

type LogActionDevice struct {
	ProductID   string                  `json:"productID"`             //产品id
	ProductName string                  `json:"productName"`           //产品名称--填写产品ID的时候会自动补充
	DeviceName  string                  `json:"deviceName"`            //选择的设备列表 指定设备的时候才需要填写(如果设备换到其他区域里,这里删除该设备)
	DeviceAlias string                  `json:"deviceAlias,omitempty"` //设备别名,只读
	Values      []*LogActionDeviceValue `json:"values"`                //传的值
}

type LogActionDeviceValue struct {
	DataID   string `json:"dataID"`   // 属性的id及事件的id,不填则取values里面的
	DataName string `json:"dataName"` //对应的物模型定义,只读
	Value    string `json:"value"`    //传的值
}

type LogTrigger struct {
	Type   TriggerType       `json:"type"`
	Device *LogTriggerDevice `json:"device,omitempty"`
}

type LogTriggerDevice struct {
	ProductID   string            `json:"productID,omitempty"`   //产品id
	DeviceName  string            `json:"deviceName,omitempty"`  //选择的列表  选择的列表, fixedDevice类型是设备名列表
	DeviceAlias string            `json:"deviceAlias,omitempty"` //设备别名,只读
	Type        TriggerDeviceType `json:"type,omitempty"`        //触发类型  connected:上线 disConnected:下线 reportProperty:属性上报 reportEvent: 事件上报
	DataID      string            `json:"dataID"`                //选择为属性或事件时需要填该字段 属性的id及事件的id aa.bb.cc
	DataName    string            `json:"dataName"`              //对应的物模型定义,只读
	Value       string            `json:"value"`                 //触发的值
}

func NewLog(scene *Info) *Log {
	if scene == nil {
		return nil
	}
	var log = Log{Type: scene.Type, AreaID: scene.AreaID, SceneID: scene.ID, Status: def.True, CreatedTime: time.Now()}
	if len(scene.If.Triggers) == 0 {
		return &log
	}
	st := scene.If.Triggers[0]
	log.Trigger = &LogTrigger{
		Type: st.Type,
	}
	if st.Type == TriggerTypeDevice && st.Device != nil {
		dev := st.Device
		log.Trigger.Device = &LogTriggerDevice{
			ProductID:   dev.ProductID,
			DeviceName:  dev.DeviceName,
			DeviceAlias: dev.DeviceAlias,
			Type:        dev.Type,
			DataID:      dev.DataID,
			DataName:    dev.DataName,
			Value:       dev.Param,
		}
	}
	return &log
}
