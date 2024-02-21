package relationDB

import (
	"gitee.com/i-Things/share/stores"
	"time"
)

// 升级任务表
type DmOtaTask struct {
	ID          int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	ProductID   string `gorm:"column:product_id;type:varchar(100);NOT NULL"`         // 产品id
	FirmwareID  int64  `gorm:"column:firmware_id;type:bigint;NOT NULL"`              // 固件id
	TaskUid     string `gorm:"column:task_uid;type:varchar(64)"`                     // 任务编号
	Type        int64  `gorm:"column:type;type:smallint;default:1;NOT NULL"`         // 升级范围1全部设备2定向升级
	UpgradeType int64  `gorm:"column:upgrade_type;type:smallint;default:1;NOT NULL"` // 升级策略:1静态升级2动态升级
	AutoRepeat  int64  `gorm:"column:auto_repeat;type:smallint;default:1;NOT NULL"`  // 是否自动重试,1:不,2自动重试
	Status      int64  `gorm:"column:status;type:smallint;default:1;NOT NULL"`       // 升级状态:1未升级2升级中3完成4已取消
	DeviceList  string `gorm:"column:device_list;type:json;NOT NULL"`                // 指定升级设备
	VersionList string `gorm:"column:version_list;type:json;NOT NULL"`               // 指定待升级版本
	stores.Time
}

func (m *DmOtaTask) TableName() string {
	return "dm_ota_task"
}

// 升级包附件列表
type DmOtaFirmwareFile struct {
	ID         int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	Name       string `gorm:"column:name;type:varchar(64)"`                // 附件名称
	FirmwareID int64  `gorm:"column:firmware_id;type:bigint;NOT NULL"`     // 固件id
	Size       int64  `gorm:"column:size;type:bigint;NOT NULL"`            // 文件大小单位bit
	Storage    string `gorm:"column:storage;type:varchar(15);NOT NULL"`    // 存储平台:minio/aliyun
	Host       string `gorm:"column:host;type:varchar(100);NOT NULL"`      // host
	FilePath   string `gorm:"column:file_path;type:varchar(100);NOT NULL"` // 文件路径,拿来下载文件
	Signature  string `gorm:"column:signature;type:char(32);NOT NULL"`     // 签名值
	stores.Time
}

func (m *DmOtaFirmwareFile) TableName() string {
	return "dm_ota_firmware_file"
}

// ota升级记录
type DmOtaTaskDevices struct {
	ID            int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	FirmwareID    int64  `gorm:"column:firmware_id;type:bigint;NOT NULL"`             // 固件id
	TaskUid       string `gorm:"column:task_uid;type:varchar(64);NOT NULL"`           // 任务批次
	ProductID     string `gorm:"column:product_id;type:varchar(100);NOT NULL"`        // 产品id
	DeviceName    string `gorm:"column:device_name;type:varchar(100);NOT NULL"`       // 设备编号
	Version       string `gorm:"column:version;type:varchar(64)"`                     // 当前版本
	TargetVersion string `gorm:"column:target_version;type:varchar(64);NOT NULL"`     // 升级包的版本
	Status        int64  `gorm:"column:status;type:integer;default:101"`              // 升级状态:101待确认 201/202/203待推送 301已推送 401升级中 501升级成功 601升级失败 701已取消
	RetryCount    int64  `gorm:"column:retry_count;type:smallint;default:0;NOT NULL"` // 重试次数,计划最多20次
	Step          int64  `gorm:"column:step;type:integer;default:0;NOT NULL"`         // OTA升级进度。1~100的整数升级进度百分比,-1升级失败,-2下载失败,-3校验失败,-4烧写失败
	Desc          string `gorm:"column:desc;type:varchar(200)"`                       // 状态详情
	stores.Time
}

func (m *DmOtaTaskDevices) TableName() string {
	return "dm_ota_task_devices"
}

// 产品固件升级包信息表
type DmOtaFirmware struct {
	ID         int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	ProductID  string `gorm:"column:product_id;type:varchar(100);NOT NULL"`    // 产品id
	Version    string `gorm:"column:version;type:varchar(64)"`                 // 固件版本
	SrcVersion string `gorm:"column:src_version;type:varchar(64)"`             // 待升级版本号
	Module     string `gorm:"column:module;type:varchar(64)"`                  // 模块名称
	Name       string `gorm:"column:name;type:varchar(64)"`                    // 固件名称
	Desc       string `gorm:"column:desc;type:varchar(200)"`                   // 描述
	Status     int64  `gorm:"column:status;type:bigint;NOT NULL"`              //升级包状态，-1：不需要验证，0：未验证，1：已验证，2：验证中，3：验证失败
	TotalSize  int64  `gorm:"column:total_size;type:bigint;NOT NULL"`          // 升级包总大小
	IsDiff     int64  `gorm:"column:is_diff;type:smallint;default:1;NOT NULL"` // 是否差分包,1:整包,2:差分
	SignMethod string `gorm:"column:sign_method;type:varchar(20);NOT NULL"`    // 签名方式:MD5/SHA256
	Extra      string `gorm:"column:extra;type:varchar(256)"`                  // 自定义推送参数
	stores.Time
}

func (m *DmOtaFirmware) TableName() string {
	return "dm_ota_firmware"
}

// DMOTAjob 表示OTA升级任务的信息
type DmOtaJob struct {
	ID                 int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
	FirmwareId         int64  `gorm:"column:firmware_id"`          // 升级包ID，升级包的唯一标识符。
	ProductId          string `gorm:"column:product_id"`           // 升级包所属产品的ProductId。
	JobType            int    `gorm:"column:job_type"`             // 升级包所属产品的JobType。
	JobStatus          int    `gorm:"column:job_status"`           // 批次状态,PLANNED：计划中,IN_PROGRESS,执行中，COMPLETED，已完成，CANCELED，已经取消
	UpgradeType        int    `gorm:"column:upgrade_type"`         // 升级策略，0-静态，1-动态
	TargetSelection    string `gorm:"column:target_selection"`     // 升级范围。可选值：0-ALL、1-SPECIFIC、2-GRAY、3-GROUP。
	SrcVersion         string `gorm:"column:src_version"`          // 待升级版本号列表。最多可传入10个版本号。用逗号分隔多个版本号
	ScheduleTime       int64  `gorm:"column:schedule_time"`        // 指定发起OTA升级的时间，单位为毫秒。
	RetryInterval      int    `gorm:"column:retry_interval"`       // 设备升级失败后，自动重试的时间间隔，单位为分钟。
	RetryCount         int    `gorm:"column:retry_count"`          // 自动重试次数。1/2/5次
	TimeoutInMinutes   int    `gorm:"column:timeout_in_minutes"`   // 设备升级超时时间，单位为分钟。
	MaximumPerMinute   int    `gorm:"column:maximum_per_minute"`   // 每分钟最多向多少个设备推送升级包下载URL。
	GrayPercent        string `gorm:"column:gray_percent"`         // 灰度比例。设置灰度升级时有效。
	TargetDeviceName   string `gorm:"column:target_device_name"`   // 定向升级的设备名称列表。最多可传入200个设备名称。以逗号分隔
	ScheduleFinishTime int64  `gorm:"column:schedule_finish_time"` // 指定结束升级的时间，单位为毫秒。
	OverwriteMode      int    `gorm:"column:overwrite_mode"`       // 是否覆盖之前的升级任务。取值：0（不覆盖）、1（覆盖）。
	DnListFileUrl      string `gorm:"column:dn_list_file_url"`     // 定向升级设备列表文件的URL。发起定向升级任务时使用。
	NeedPush           int    `gorm:"column:need_push"`            // 物联网平台是否主动向设备推送升级任务。
	NeedConfirm        int    `gorm:"column:need_confirm"`         // 是否需要App确认OTA升级。
	GroupId            string `gorm:"column:group_id"`             // 分组ID。仅发起分组升级任务时使用。
	GroupType          string `gorm:"column:group_type"`           // 分组类型，仅可取值LINK_PLATFORM。仅发起分组升级任务时使用。
	DownloadProtocol   string `gorm:"column:download_protocol"`    // 升级包下载协议。可选：0-HTTPS、1-MQTT。
	MultiModuleMode    int    `gorm:"column:multi_module_mode"`    // 设备是否支持多模块同时升级。
	stores.Time
}

func (m *DmOtaJob) TableName() string {
	return "dm_ota_job"
}

type DmOtaUpgradeTask struct {
	ID          int64     `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
	DestVersion string    `gorm:"column:dest_version" ` // 升级的目标OTA升级包版本
	DeviceName  string    `gorm:"column:device_name" `  // 设备名称
	FirmwareId  int64     `gorm:"column:firmware_id"`   // 升级包ID
	JobId       int64     `gorm:"column:job_id" `       // 升级批次ID
	ProductId   string    `gorm:"column:product_id" `   // 设备所属产品的ProductKey
	ProductName string    `gorm:"column:product_name"`  // 设备所属产品的名称
	Step        int64     `gorm:"column:step"`          // 当前的升级进度
	Module      string    `gorm:"column:module"`        //模块名称
	SrcVersion  string    `gorm:"column:src_version"`   // 设备的原固件版本
	TaskDesc    string    `gorm:"column:task_desc" `    // 升级作业描述信息
	TaskStatus  int       `gorm:"column:task_status"`   // 设备升级状态
	Timeout     string    `gorm:"column:timeout"`       // 设备升级超时时间，单位是分钟
	UtcCreate   time.Time `gorm:"column:utc_create"`    // 升级作业创建时的时间，UTC格式
	UtcModified time.Time `gorm:"column:utc_modified"`  // 升级作业最后一次修改时的时间，UTC格式
	stores.Time
}

func (m *DmOtaUpgradeTask) TableName() string {
	return "dm_ota_upgrade_task"
}

type DmOtaModule struct {
	ID            int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
	ModuleName    string `gorm:"column:module_name"`    // OTA模块名称，产品下唯一且不可修改。仅支持英文字母、数字、英文句号（.）、短划线（-）和下划线（_）。长度限制为1~64个字符。
	ProductId     string `gorm:"column:product_id"`     // OTA模块所属产品的ProductId。
	ModuleVersion string `gorm:"column:module_version"` //模块版本
	AliasName     string `gorm:"column:alias_name"`     // OTA模块别名。支持中文、英文字母、数字、英文句号（.）、短划线（-）和下划线（_），长度限制为1~64个字符。
	Desc          string `gorm:"column:desc"`           // OTA模块的描述信息，支持最多100个字符。
	stores.Time
}

func (m *DmOtaModule) TableName() string {
	return "dm_ota_module"
}
