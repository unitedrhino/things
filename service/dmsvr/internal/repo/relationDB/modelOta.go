package relationDB

import (
	"database/sql"
	"gitee.com/i-Things/share/domain/deviceMsg/msgOta"
	"gitee.com/i-Things/share/stores"
)

//// 升级任务表
//type DmOtaFirmwareDevice struct {
//	ID          int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
//	ProductID   string `gorm:"column:product_id;type:varchar(100);NOT NULL"`         // 产品id
//	FirmwareID  int64  `gorm:"column:firmware_id;type:bigint;NOT NULL"`              // 固件id
//	TaskUid     string `gorm:"column:task_uid;type:varchar(64)"`                     // 任务编号
//	Type        int64  `gorm:"column:type;type:smallint;default:1;NOT NULL"`         // 升级范围1全部设备2定向升级
//	UpgradeType int64  `gorm:"column:upgrade_type;type:smallint;default:1;NOT NULL"` // 升级策略:1静态升级2动态升级
//	AutoRepeat  int64  `gorm:"column:auto_repeat;type:smallint;default:1;NOT NULL"`  // 是否自动重试,1:不,2自动重试
//	Msg      int64  `gorm:"column:status;type:smallint;default:1;NOT NULL"`       // 升级状态:1未升级2升级中3完成4已取消
//	DeviceList  string `gorm:"column:device_list;type:json;NOT NULL"`                // 指定升级设备
//	VersionList string `gorm:"column:version_list;type:json;NOT NULL"`               // 指定待升级版本
//	stores.Time
//}
//
//func (m *DmOtaFirmwareDevice) TableName() string {
//	return "dm_ota_task"
//}

// 升级包附件列表
type DmOtaFirmwareFile struct {
	ID         int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	Name       string `gorm:"column:name;type:varchar(64)"`                 // 附件名称
	FirmwareID int64  `gorm:"column:firmware_id;type:bigint;NOT NULL"`      // 固件id
	Size       int64  `gorm:"column:size;type:bigint;NOT NULL"`             // 文件大小单位bit
	FilePath   string `gorm:"column:file_path;type:varchar(2048);NOT NULL"` // 文件路径,拿来下载文件
	Signature  string `gorm:"column:signature;type:char(32);NOT NULL"`      // 签名值
	FileMd5    string `gorm:"column:file_md5;type:char(32);NOT NULL"`
	stores.Time
}

func (m *DmOtaFirmwareFile) TableName() string {
	return "dm_ota_firmware_file"
}

// ota升级记录
type DmOtaTaskDevice struct {
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

func (m *DmOtaTaskDevice) TableName() string {
	return "dm_ota_task_device"
}

// 产品固件升级包信息表
type DmOtaFirmwareInfo struct {
	ID             int64                `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	ProductID      string               `gorm:"column:product_id;type:varchar(100);uniqueIndex:tc_un;NOT NULL"` // 产品id
	Version        string               `gorm:"column:version;type:varchar(64);uniqueIndex:tc_un"`              // 固件版本
	SrcVersion     string               `gorm:"column:src_version;type:varchar(64)"`                            // 待升级版本号
	Name           string               `gorm:"column:name;type:varchar(64)"`                                   // 固件名称
	ModuleCode     string               `gorm:"column:module_code;type:varchar(64)"`                            // 固件名称
	Desc           string               `gorm:"column:desc;type:varchar(200)"`                                  // 描述
	Status         int64                `gorm:"column:status;type:bigint;NOT NULL"`                             //升级包状态，-1：不需要验证，0：未验证，1：已验证，2：验证中，3：验证失败
	TotalSize      int64                `gorm:"column:total_size;type:bigint;NOT NULL"`                         // 升级包总大小
	IsDiff         int64                `gorm:"column:is_diff;type:smallint;default:1;NOT NULL"`                // 是否差分包,1:整包,2:差分
	SignMethod     string               `gorm:"column:sign_method;type:varchar(20);NOT NULL"`                   // 签名方式:MD5/SHA256
	Extra          string               `gorm:"column:extra;type:varchar(256)"`                                 // 自定义推送参数
	IsNeedToVerify int64                `gorm:"column:is_need_to_verify;type:smallint;default:2;NOT NULL"`      // 是否需要验证
	Files          []*DmOtaFirmwareFile `gorm:"foreignKey:FirmwareID;references:ID"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:tc_un;"`
}

func (m *DmOtaFirmwareInfo) TableName() string {
	return "dm_ota_firmware_info"
}

// DMOTAjob 表示OTA升级任务的信息
type DmOtaFirmwareJob struct {
	ID          int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
	FirmwareID  int64  `gorm:"column:firmware_id"`                           // 升级包ID，升级包的唯一标识符。
	ProductID   string `gorm:"column:product_id;type:varchar(100);NOT NULL"` // 产品id
	Type        int64  `gorm:"column:type"`                                  // 升级包所属产品的JobType。 验证升级包:1  批量升级:2
	Status      int64  `gorm:"column:status"`                                // 批次状态,计划中:1  执行中:2  已完成:3  已取消:4
	UpgradeType int64  `gorm:"column:upgrade_type"`                          // 升级策略，1-静态，2-动态

	/*
		待升级版本号列表。
			发起全量升级（TargetSelection=ALL）和灰度升级（TargetSelection=GRAY）任务时，可以传入该参数。
			使用差分升级包发起全量升级和灰度升级任务时，该参数值需指定为差分升级包的待升级版本号（SrcVersion）。
			发起定向升级（TargetSelection=SPECIFIC）或分组升级（TargetSelection=GROUP）任务时，不能传入该参数。
			可以调用 QueryDeviceDetail ，查看设备 OTA 模块版本号（FirmwareVersion）。
			列表中不能有重复的版本号。
			最多可传入 10 个版本号。
	*/
	SrcVersions []string `gorm:"column:src_versions;type:json;serializer:json"` // 待升级版本号列表。最多可传入10个版本号。用逗号分隔多个版本号

	RetryInterval    int64 `gorm:"column:retry_interval"`     // 设备升级失败后，自动重试的时间间隔，单位为分钟。 动态升级 静态升级
	RetryCount       int64 `gorm:"column:retry_count"`        // 自动重试次数。1/2/5次 动态升级 静态升级
	TimeoutInMinutes int64 `gorm:"column:timeout_in_minutes"` // 设备升级超时时间，单位为分钟。 动态升级 静态升级
	MaximumPerMinute int64 `gorm:"column:maximum_per_minute"` // 每分钟最多向多少个设备推送升级包下载URL。 动态升级 静态升级

	/*
			是否覆盖之前的升级任务。取值：

			2（默认）：不覆盖。若设备已有升级任务，则只执行已有任务。
			1：覆盖。设备只执行新的升级任务。此时 MultiModuleMode 不能传入 true。
		动态升级 静态升级
	*/
	IsOverwriteMode int64 `gorm:"column:is_overwrite_mode;default:2"` // 是否覆盖之前的升级任务。取值：2（不覆盖）、1（覆盖）。
	/*
				物联网平台是否主动向设备推送升级任务。
			1（默认）：是。批次任务创建完成后，物联网平台主动将升级任务，直接推送给升级范围内的在线设备。
				此时，设备仍可主动向物联网平台发起请求，来获取 OTA 升级任务信息。
			2：否。设备必须通过向物联网平台发起请求，来获取 OTA 升级任务信息。
		动态升级
	*/
	IsNeedPush int64 `gorm:"column:is_need_push"` // 物联网平台是否主动向设备推送升级任务。

	/*
		如需自主控制设备 OTA 升级时，可配置此参数，通过手机 App 来控制，设备是否可进行 OTA 升级。手机 App 需您自行开发。
			2（默认）：否。直接按照 NeedPush 设置，获取 OTA 升级任务信息。
			1：是。设备无法获取 OTA 升级任务，需 App 侧确认 OTA 升级后，才能按照 NeedPush 设置，获取 OTA 升级任务信息。
	*/
	IsNeedConfirm   int64                   `gorm:"is_need_confirm"`
	TargetSelection int64                   `gorm:"column:target_selection;type:bigint;default:1"` //升级范围。 1：全量升级。 2：定向升级。 3：灰度升级。 4：分组升级 5: 区域升级
	TargetID        int64                   `gorm:"column:target_id;type:bigint;default:0"`
	Dynamic         DmOtaFirmwareJobDynamic `gorm:"embedded"`
	Static          DmOtaFirmwareJobStatic  `gorm:"embedded"`
	Firmware        *DmOtaFirmwareInfo      `gorm:"foreignKey:FirmwareID;references:ID"`
	Files           []*DmOtaFirmwareFile    `gorm:"foreignKey:FirmwareID;references:FirmwareID"`
	stores.Time
}

func (m *DmOtaFirmwareJob) TableName() string {
	return "dm_ota_firmware_job"
}

type DmOtaFirmwareJobDynamic struct {
	/*
			动态升级模式。取值范围：
			1（默认）：除了升级当前满足升级条件的设备，还将持续检查设备是否满足升级条件，对满足升级条件的设备进行升级。
			2：仅对后续上报新版本号的设备生效。
		动态升级
	*/
	DynamicMode int64 `gorm:"column:dynamic_mode"` //
}
type DmOtaFirmwareJobStatic struct {
	/*
			定向升级的设备名称列表。
			使用差分升级包进行定向升级时，要升级的设备的当前 OTA 模块版本号需与差分升级包的待升级版本号（SrcVersion）相同。
			可以调用 QueryDeviceDetail ，查看设备 OTA 模块版本号（FirmwareVersion）。
			列表中的设备所属的产品必须与升级包所属产品一致。
			列表中不能有重复的设备名称。
			最多可传入 200 个设备名称。
		静态升级
	*/
	TargetDeviceNames []string `gorm:"column:target_device_names;type:json;serializer:json"` // 定向升级的设备名称列表。最多可传入200个设备名称。以逗号分隔
	/*
			指定结束升级的时间。
			结束时间距发起时间（ScheduleTime）最少 1 小时，最多为 30 天。取值为 13 位毫秒值时间戳。
			不传入该参数，则表示不会强制结束升级。
		静态升级
	*/
	ScheduleFinishTime int64 `gorm:"column:schedule_finish_time"` // 指定结束升级的时间，单位为毫秒。
	/*
			指定发起 OTA 升级的时间。
			定时时间范围需为当前时间的 5 分钟后至 7 天内。取值为秒时间戳。
			不传入该参数，则表示立即升级。
		静态升级
	*/
	ScheduleTime int64 `gorm:"column:schedule_time"`
}

type DmOtaFirmwareDevice struct {
	ID              int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
	FirmwareID      int64               `gorm:"column:firmware_id;uniqueIndex:tc_un"`                    // 升级包ID
	SrcVersion      string              `gorm:"column:src_version;type:varchar(125)"`                    // 设备的原固件版本
	DestVersion     string              `gorm:"column:dest_version;type:varchar(125)"`                   // 设备的目标固件版本
	ProductID       string              `gorm:"column:product_id;type:varchar(20);uniqueIndex:tc_un"`    // 设备所属产品的productID
	DeviceName      string              `gorm:"column:device_name;type:varchar(100);uniqueIndex:tc_un" ` // 设备名称
	JobID           int64               `gorm:"column:job_id;type:BIGINT" `                              // 升级批次ID
	Step            int64               `gorm:"column:step;type:BIGINT"`                                 // 当前的升级进度  0-100%    -1：升级失败。-2：下载失败。-3：校验失败。-4：烧写失败。
	Detail          string              `gorm:"column:detail;type:varchar(256)"`                         //详情
	Status          msgOta.DeviceStatus `gorm:"column:status;type:BIGINT"`                               // 设备升级作业状态。1：待确认。 2：待推送。 3：已推送。  4：升级中。 5:升级成功 6: 升级失败. 7:已取消
	PushTime        sql.NullTime        `gorm:"column:push_time"`                                        //推送时间
	LastFailureTime sql.NullTime        `gorm:"column:last_failure_time"`                                //最后失败时间
	stores.NoDelTime
	RetryCount  int64                `gorm:"column:retry_count;default:0"` // 自动重试次数
	DeletedTime stores.DeletedTime   `gorm:"column:deleted_time;default:0;uniqueIndex:tc_un;"`
	Job         *DmOtaFirmwareJob    `gorm:"foreignKey:JobID;references:ID"`
	Firmware    *DmOtaFirmwareInfo   `gorm:"foreignKey:FirmwareID;references:ID"`
	Files       []*DmOtaFirmwareFile `gorm:"foreignKey:FirmwareID;references:FirmwareID"`
}

func (m *DmOtaFirmwareDevice) TableName() string {
	return "dm_ota_firmware_device"
}

// 产品固件升级包信息表
type DmOtaModuleInfo struct {
	ID        int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	ProductID string `gorm:"column:product_id;uniqueIndex:tc_un;type:varchar(20)" ` // 设备所属产品的productID
	Code      string `gorm:"column:code;uniqueIndex:tc_un;type:varchar(64)"`
	Name      string `gorm:"column:name;type:varchar(64)"`  // 固件名称
	Desc      string `gorm:"column:desc;type:varchar(200)"` // 描述
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:tc_un;"`
}

func (m *DmOtaModuleInfo) TableName() string {
	return "dm_ota_module_info"
}
