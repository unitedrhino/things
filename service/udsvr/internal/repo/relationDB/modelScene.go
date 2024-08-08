package relationDB

import (
	"database/sql"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"github.com/i-Things/things/service/udsvr/internal/domain/scene"
)

type UdSceneInfo struct {
	ID          int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`                       // id编号
	TenantCode  stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`                           // 租户编码
	ProjectID   stores.ProjectID  `gorm:"column:project_id;type:bigint;default:0;NOT NULL"`                       // 项目ID(雪花ID)
	AreaID      stores.AreaID     `gorm:"column:area_id;index:project_id_area_id;type:bigint;default:0;NOT NULL"` // 项目区域ID(雪花ID)
	FlowPath    []*scene.FlowInfo `gorm:"column:flow_path;type:json;serializer:json;"`                            //执行路径
	DeviceMode  scene.DeviceMode  `gorm:"column:tenant_code;type:VARCHAR(50);default:'single'"`                   //设备模式
	ProductID   string            `gorm:"column:product_id;index;type:VARCHAR(25);default:''"`                    //产品id
	DeviceName  string            `gorm:"column:device_name;type:VARCHAR(255);default:''"`                        //设备名
	DeviceAlias string            `gorm:"column:device_alias;type:VARCHAR(255);default:''"`                       //设备别名
	Tag         string            `gorm:"column:tag;type:VARCHAR(128);NOT NULL;default:normal"`                   //标签 admin: 管理员 normal: 普通
	HeadImg     string            `gorm:"column:head_img;type:VARCHAR(256);NOT NULL"`                             // 头像
	Logo        string            `gorm:"column:logo;type:VARCHAR(25);NOT NULL"`                                  // logo
	Name        string            `gorm:"column:name;type:varchar(100);NOT NULL"`                                 // 名称
	Desc        string            `gorm:"column:desc;type:varchar(200);NOT NULL"`                                 // 描述
	Type        scene.SceneType   `gorm:"column:type;type:VARCHAR(25);NOT NULL"`                                  //auto manual
	LastRunTime sql.NullTime      `gorm:"column:last_run_time;index;default:null"`
	Status      def.Bool          `gorm:"column:status;type:BIGINT;default:1"` //状态
	Body        string            `gorm:"column:body;type:VARCHAR(1024)"`      // 自定义数据
	UdSceneIf   `gorm:"embedded;embeddedPrefix:if_"`
	UdSceneWhen `gorm:"embedded;embeddedPrefix:when_"`
	UdSceneThen `gorm:"embedded;embeddedPrefix:then_"`
	stores.SoftTime
}

func (m *UdSceneInfo) TableName() string {
	return "ud_scene_info"
}

type UdSceneIf struct {
	Triggers []*UdSceneIfTrigger `gorm:"foreignKey:SceneID;references:ID"`
}

type UdSceneIfTrigger struct {
	ID              int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
	Type            scene.TriggerType `gorm:"column:type;type:VARCHAR(25);NOT NULL"`            //触发类型 device: 设备触发 timer: 定时触发
	SceneID         int64             `gorm:"column:scene_id;index;type:bigint"`                // 场景id编号
	AreaID          int64             `gorm:"column:area_id;type:bigint;default:0;"`            // 项目区域ID(雪花ID)
	Order           int64             `gorm:"column:order;type:BIGINT;default:1;NOT NULL"`      // 排序序号
	Status          int64             `gorm:"column:status;type:BIGINT;default:1"`              //状态 同步场景联动的status
	LastRunTime     sql.NullTime      `gorm:"column:last_run_time;index;default: NULL"`         //最后执行时间
	LastTriggerTime sql.NullTime      `gorm:"column:last_trigger_time;index;default: NULL"`     //最后触发时间

	Device    UdSceneTriggerDevice  `gorm:"embedded;embeddedPrefix:device_"`
	Timer     UdSceneTriggerTimer   `gorm:"embedded;embeddedPrefix:timer_"`
	Weather   UdSceneTriggerWeather `gorm:"embedded;embeddedPrefix:weather_"`
	SceneInfo *UdSceneInfo          `gorm:"foreignKey:ID;references:SceneID"`
	stores.SoftTime
}

func (m *UdSceneIfTrigger) TableName() string {
	return "ud_scene_if_trigger"
}

type UdSceneTriggerTimer struct {
	ExecType      scene.ExecType   `gorm:"column:exec_type;type:VARCHAR(25);"`  //执行方式
	ExecAdd       int64            `gorm:"column:exec_add;type:bigint;"`        //如果是日出日落模式,则为日出日落前后的秒数
	ExecAt        int64            `gorm:"column:exec_at;type:bigint;"`         //执行时间 从0点加起来的秒数 如 1点就是 1*60*60
	ExecRepeat    int64            `gorm:"column:exec_repeat;type:bigint;"`     //重复 二进制周日到周六 11111111 这个参数只有定时触发才有
	ExecLoopStart int64            `gorm:"column:exec_loop_start;type:bigint;"` //循环执行起始时间配置
	ExecLoopEnd   int64            `gorm:"column:exec_loop_end;type:bigint;"`
	ExecLoop      int64            `gorm:"column:exec_loop;type:bigint;"`
	RepeatType    scene.RepeatType `gorm:"column:repeat_type;type:VARCHAR(25);"`
}

type UdSceneTriggerWeather struct {
	Type     scene.WeatherType `gorm:"column:type;type:VARCHAR(25);default:''"`              //天气的类型
	TermType string            `gorm:"column:term_type;type:VARCHAR(255);default:''"`        //动态条件类型  eq: 相等  not:不相等  btw:在xx之间  gt: 大于  gte:大于等于 lt:小于  lte:小于等于   in:在xx值之间
	Values   []string          `gorm:"column:values;type:json;serializer:json;default:'[]'"` //条件值 参数根据动态条件类型会有多个参数
}

type UdSceneTriggerDevice struct {
	ProductID        string                  `gorm:"column:product_id;index;type:VARCHAR(25);default:''"`  //产品id
	SelectType       scene.SelectType        `gorm:"column:select_type;type:VARCHAR(25);default:''"`       //设备选择方式  all: 全部 fixed:指定的设备
	GroupID          int64                   `gorm:"column:group_id;index;type:bigint;default:0"`          //group类型传GroupID
	DeviceName       string                  `gorm:"column:device_name;type:VARCHAR(255);default:''"`      //选择的列表  选择的列表, fixed类型是设备名列表
	DeviceAlias      string                  `gorm:"column:device_alias;type:VARCHAR(255);default:''"`     //设备别名
	Type             scene.TriggerDeviceType `gorm:"column:type;type:VARCHAR(25);default:''"`              //触发类型  connected:上线 disConnected:下线 reportProperty:属性上报 reportEvent: 事件上报
	DataID           string                  `gorm:"column:data_id;type:VARCHAR(255);default:''"`          //选择为属性或事件时需要填该字段 属性的id及事件的id aa.bb.cc
	DataName         string                  `gorm:"column:data_name;type:VARCHAR(255);default:''"`        //对应的物模型定义,只读
	TermType         string                  `gorm:"column:term_type;type:VARCHAR(255);default:''"`        //动态条件类型  eq: 相等  not:不相等  btw:在xx之间  gt: 大于  gte:大于等于 lt:小于  lte:小于等于   in:在xx值之间
	Values           []string                `gorm:"column:values;type:json;serializer:json;default:'[]'"` //比较条件列表
	SchemaAffordance string                  `gorm:"column:schema_affordance;type:VARCHAR(500);default:''"`
	Body             string                  `gorm:"column:body;type:VARCHAR(1024)"` // 自定义数据
	StateKeep        UdStateKeep             `gorm:"embedded;embeddedPrefix:state_keep_"`
	FirstTriggerTime sql.NullTime            `gorm:"column:first_trigger_time;index;default: NULL"` //最后触发时间
}

// StateKeep 状态保持
type UdStateKeep struct {
	Type  scene.StateKeepType `gorm:"column:type;type:VARCHAR(25);default:''"` //持续时间: duration  重复次数 repeating
	Value int64               `gorm:"column:value;type:bigint;default:0"`      //持续的时间或重复的次数
}

type UdSceneWhen struct {
	ValidRanges   scene.WhenRanges `gorm:"column:validRanges;type:json;serializer:json"`
	InvalidRanges scene.WhenRanges `gorm:"column:invalidRanges;type:json;serializer:json"`
	Conditions    scene.Conditions `gorm:"column:conditions;type:json;serializer:json"`
}

type UdSceneThen struct {
	Actions []*UdSceneThenAction `gorm:"foreignKey:scene_id;references:ID"`
}

type UdSceneThenAction struct {
	ID         int64               `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
	TenantCode stores.TenantCode   `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`     // 租户编码
	SceneID    int64               `gorm:"column:scene_id;index;type:bigint"`                // 场景id编号
	Order      int64               `gorm:"column:order;type:BIGINT;default:1;NOT NULL"`      // 排序序号
	Type       scene.ActionType    `gorm:"column:type;type:VARCHAR(25);NOT NULL"`
	Delay      int64               `gorm:"column:delay;type:bigint"`
	Device     UdSceneActionDevice `gorm:"embedded;embeddedPrefix:device_"`
	Notify     *scene.ActionNotify `gorm:"column:notify;type:json;serializer:json"`
	Alarm      *scene.ActionAlarm  `gorm:"column:alarm;type:json;serializer:json"`
	Scene      UdSceneActionScene  `gorm:"embedded;embeddedPrefix:scene_"`
}

func (m *UdSceneThenAction) TableName() string {
	return "ud_scene_then_action"
}

type UdSceneActionDevice struct {
	//ProjectID        int64                  `gorm:"column:project_id;type:bigint;default:2;NOT NULL"`  // 项目ID(雪花ID)
	AreaID           int64                  `gorm:"column:area_id;type:bigint;default:2;"` // 项目区域ID(雪花ID)
	AreaName         string                 `gorm:"column:area_name;index;type:VARCHAR(100);default:''"`
	ProductID        string                 `gorm:"column:product_id;index;type:VARCHAR(25);default:''"`    //产品id
	ProductName      string                 `gorm:"column:product_name;index;type:VARCHAR(200);default:''"` //产品id
	SelectType       scene.SelectType       `gorm:"column:select_type;type:VARCHAR(25);default:''"`         //设备选择方式
	DeviceName       string                 `gorm:"column:device_name;type:VARCHAR(255);"`                  //选择的列表  选择的列表, fixed类型是设备名列表
	DeviceAlias      string                 `gorm:"column:device_alias;type:VARCHAR(255);"`                 //设备别名
	GroupID          int64                  `gorm:"column:group_id;index;type:bigint"`                      //group类型传GroupID
	Type             scene.ActionDeviceType `gorm:"column:type;type:VARCHAR(25);default:''"`
	DataID           string                 `gorm:"column:data_id;index;type:VARCHAR(100);default:''"`
	Value            string                 `gorm:"column:value;index;type:VARCHAR(500);default:''"`
	DataName         string                 `gorm:"column:data_name;type:VARCHAR(500);default:''"` //对应的物模型定义,只读
	SchemaAffordance string                 `gorm:"column:schema_affordance;type:VARCHAR(500);default:''"`
	Values           scene.DeviceValues     `gorm:"column:values;type:json;serializer:json"`
	Body             string                 `gorm:"column:body;type:VARCHAR(1024)"` // 自定义数据
}

type UdSceneActionScene struct {
	SceneID int64 `gorm:"column:scene_id;index;type:bigint"` // 场景id编号
	AreaID  int64 `gorm:"column:area_id;type:bigint;"`       // 项目区域ID(雪花ID)

}

type UdSceneLog struct {
	ID         int64              `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`    // id编号
	TenantCode stores.TenantCode  `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`        // 租户编码
	ProjectID  stores.ProjectID   `gorm:"column:project_id;type:bigint;default:0;NOT NULL"`    // 项目ID(雪花ID)
	AreaID     stores.AreaID      `gorm:"column:area_id;index;type:bigint;default:0;NOT NULL"` // 项目区域ID(雪花ID)
	SceneID    int64              `gorm:"column:scene_id;index;type:bigint;NOT NULL"`          // 项目区域ID(雪花ID)
	Type       scene.SceneType    `gorm:"column:type;type:VARCHAR(25);NOT NULL"`               //auto manual
	Trigger    *scene.LogTrigger  `gorm:"column:trigger;type:json;serializer:json"`
	Actions    []*scene.LogAction `gorm:"column:actions;type:json;serializer:json"`
	Status     def.Bool           `gorm:"column:status;type:BIGINT;default:1"` //状态
	SceneInfo  *UdSceneInfo       `gorm:"foreignKey:ID;references:SceneID"`
	stores.OnlyTime
}

func (m *UdSceneLog) TableName() string {
	return "ud_scene_log"
}
