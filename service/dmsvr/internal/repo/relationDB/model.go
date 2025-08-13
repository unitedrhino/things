package relationDB

import (
	"context"
	"database/sql"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/product"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/productCustom"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/protocol"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/protocols"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"time"
)

type DmExample struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
}

// 设备信息表
type DmDeviceInfo struct {
	ID              int64               `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	TenantCode      dataType.TenantCode `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL"`                                                                                // 租户编码
	ProjectID       dataType.ProjectID  `gorm:"column:project_id;index:project_id_area_id;type:bigint;default:0;NOT NULL"`                                                         // 项目ID(雪花ID)
	AreaID          dataType.AreaID     `gorm:"column:area_id;index:project_id_area_id;type:bigint;default:0;NOT NULL"`                                                            // 项目区域ID(雪花ID)
	AreaIDPath      dataType.AreaIDPath `gorm:"column:area_id_path;type:varchar(100);default:'';NOT NULL"`                                                                         // 项目区域ID路径(雪花ID)
	DeptID          dataType.DeptID     `gorm:"column:dept_id;type:bigint;default:0;NOT NULL"`                                                                                     // 部门ID
	DeptIDPath      dataType.DeptIDPath `gorm:"column:dept_id_path;type:varchar(100);default:'';NOT NULL"`                                                                         // 部门ID路径
	DeptUpdatedTime time.Time           `gorm:"column:dept_updated_time;default:null"`                                                                                             //部门更新时间
	ProductID       string              `gorm:"column:product_id;type:varchar(100);index:idx_dm_device_info_pd_dn;uniqueIndex:idx_dm_device_info_product_id_deviceName;NOT NULL"`  // 产品id
	DeviceName      string              `gorm:"column:device_name;index:idx_dm_device_info_pd_dn;uniqueIndex:idx_dm_device_info_product_id_deviceName;type:varchar(100);NOT NULL"` // 设备名称
	DeviceAlias     string              `gorm:"column:device_alias;type:varchar(100);NOT NULL"`                                                                                    // 设备别名
	Position        stores.Point        `gorm:"column:position;"`                                                                                                                  // 设备的位置(默认百度坐标系BD09)
	RatedPower      int64               `gorm:"column:rated_power;index;type:bigint;NOT NULL;default:0"`                                                                           // 额定功率:单位w/h
	Secret          string              `gorm:"column:secret;type:varchar(50);NOT NULL"`                                                                                           // 设备秘钥
	Cert            string              `gorm:"column:cert;type:varchar(512);NOT NULL"`                                                                                            // 设备证书
	Imei            string              `gorm:"column:imei;type:varchar(25);NOT NULL"`                                                                                             // IMEI号信息
	Mac             string              `gorm:"column:mac;type:varchar(17);NOT NULL"`                                                                                              // MAC号信息
	DeviceType      int64               `gorm:"column:device_type;index;type:smallint;default:1"`                                                                                  // 设备类型:1:设备,2:网关,3:子设备
	Version         string              `gorm:"column:version;index;type:varchar(64);NOT NULL"`                                                                                    // 固件版本
	//ModuleVersion  map[string]string `gorm:"column:module_version;type:json;serializer:json;NOT NULL;default:'{}'"`      // 所有模块的版本
	HardInfo           string                 `gorm:"column:hard_info;type:varchar(64);NOT NULL"`                               // 模组硬件型号
	SoftInfo           string                 `gorm:"column:soft_info;type:varchar(64);NOT NULL"`                               // 模组软件版本
	MobileOperator     def.MobileOperator     `gorm:"column:mobile_operator;type:smallint;default:10;NOT NULL"`                 // 移动运营商:1)移动 2)联通 3)电信 4)广电 10) 无
	Phone              sql.NullString         `gorm:"column:phone;type:varchar(20)"`                                            // 手机号
	Iccid              sql.NullString         `gorm:"column:iccid;type:varchar(20)"`                                            // SIM卡卡号
	Address            string                 `gorm:"column:address;type:varchar(512);default:''"`                              // 所在地址
	Adcode             string                 `gorm:"column:adcode;type:varchar(125);default:''"`                               // 地区编码
	Tags               map[string]string      `gorm:"column:tags;type:json;serializer:json;NOT NULL;default:'{}'"`              // 设备标签
	SchemaAlias        map[string]string      `gorm:"column:schema_alias;type:json;serializer:json;NOT NULL;default:'{}'"`      // 设备物模型别名,如果是结构体类型则key为xxx.xxx
	Rssi               int64                  `gorm:"column:rssi;type:bigint;default:0;NOT NULL"`                               // 设备信号（信号极好[-55— 0]，信号好[-70— -55]，信号一般[-85— -70]，信号差[-100— -85]）
	ProtocolConf       map[string]string      `gorm:"column:protocol_conf;type:json;serializer:json;NOT NULL;default:'{}'"`     // 主协议配置
	SubProtocolConf    map[string]string      `gorm:"column:sub_protocol_conf;type:json;serializer:json;NOT NULL;default:'{}'"` // 子模块自定义协议配置
	DeviceImg          string                 `gorm:"column:device_img;type:varchar(200);default:''"`                           // 设备图片
	File               string                 `gorm:"column:file;type:varchar(200);default:''"`                                 // 设备相关文件
	IsOnline           int64                  `gorm:"column:is_online;type:smallint;default:2;NOT NULL"`                        // 是否在线,1是2否
	FirstLogin         sql.NullTime           `gorm:"column:first_login"`                                                       // 激活时间
	LastLogin          sql.NullTime           `gorm:"column:last_login"`                                                        // 最后上线时间
	LastOffline        sql.NullTime           `gorm:"column:last_offline"`                                                      // 最后下线线时间
	FirstBind          sql.NullTime           `gorm:"column:first_bind"`                                                        // 首次绑定时间
	LastBind           sql.NullTime           `gorm:"column:last_bind"`                                                         // 最后一次绑定时间
	LogLevel           int64                  `gorm:"column:log_level;type:smallint;default:1;NOT NULL"`                        // 日志级别:1)关闭 2)错误 3)告警 4)信息 5)调试
	UserID             int64                  `gorm:"column:user_id;type:BIGINT;default:1"`                                     // 用户id
	Sort               int64                  `gorm:"column:sort;type:BIGINT;default:100"`                                      // 排序
	Status             def.DeviceStatus       `gorm:"column:status;index;type:smallint;default:1;NOT NULL"`                     // 设备状态 1-未激活，2-在线，3-离线 4-异常(频繁上下线,告警中) 5-禁用
	IsEnable           int64                  `gorm:"column:is_enable;index;type:smallint;default:1;"`                          // 是否启用: 1:是 2:否
	ExpTime            sql.NullTime           `gorm:"column:exp_time"`                                                          // 过期时间,为0不限制
	NeedConfirmJobID   int64                  `gorm:"column:need_confirm_job_id;type:smallint;default:0;"`                      // 需要app确认升级的任务ID,为0是没有
	NeedConfirmVersion string                 `gorm:"column:need_confirm_version;type:varchar(128);default:'';"`                // 待确认升级的版本
	LastIp             string                 `gorm:"column:last_ip;type:varchar(128);default:'';"`                             // 最后登录的ip地址
	LastLocalIp        string                 `gorm:"column:last_local_ip;type:varchar(128);default:'';"`                       // 局域网的ip地址
	BelongGroup        map[string]def.IDsInfo `gorm:"column:belong_group;type:json;serializer:json;default:'{}'"`
	stores.NoDelTime
	Desc        string                  `gorm:"column:desc;type:varchar(200)"`        // 描述
	Distributor stores.IDPathWithUpdate `gorm:"embedded;embeddedPrefix:distributor_"` // 代理的id,如果为空,则未参与分销
	DeletedTime stores.DeletedTime      `gorm:"column:deleted_time;default:0;uniqueIndex:idx_dm_device_info_product_id_deviceName"`
	ProductInfo *DmProductInfo          `gorm:"foreignKey:ProductID;references:ProductID"` // 添加外键
}

func (m *DmDeviceInfo) TableName() string {
	return "dm_device_info"
}

type ConnParam struct {
	ClientID string `json:"clientID"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// 设备信息表
type DmDeviceMsgCount struct {
	ID   int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	Type deviceLog.MsgType `gorm:"column:type;index;type:VARCHAR(50);uniqueIndex:idx_dm_device_msg_count_date_type;NOT NULL"` // 消息类型
	Num  int64             `gorm:"column:num;type:bigint;default:0"`                                                          //数量
	Date time.Time         `gorm:"column:date;NOT NULL;uniqueIndex:idx_dm_device_msg_count_date_type"`                        //统计的日期
	stores.OnlyTime
}

func (m *DmDeviceMsgCount) TableName() string {
	return "dm_device_msg_count"
}

var ClearDeviceInfo func(ctx context.Context, dev devices.Core) error

func (u *DmDeviceInfo) AfterSave(tx *gorm.DB) (err error) {
	if ClearDeviceInfo != nil && u.DeviceName != "" {
		err := ClearDeviceInfo(tx.Statement.Context, devices.Core{
			ProductID:  u.ProductID,
			DeviceName: u.DeviceName,
		})
		if err != nil {
			logx.WithContext(tx.Statement.Context).Error(err, u)
		}
	}
	return nil
}

func (u *DmDeviceInfo) AfterDelete(tx *gorm.DB) (err error) {
	if ClearDeviceInfo != nil && u.DeviceName != "" {
		err := ClearDeviceInfo(tx.Statement.Context, devices.Core{
			ProductID:  u.ProductID,
			DeviceName: u.DeviceName,
		})
		if err != nil {
			logx.WithContext(tx.Statement.Context).Error(err, u)
		}
	}
	return nil
}

type DmDeviceModuleVersion struct {
	ID         int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	ProductID  string `gorm:"column:product_id;type:varchar(100);index:idx_dm_device_module_version_pd_dn;uniqueIndex:idx_dm_device_module_version_product_id_deviceName;NOT NULL"`  // 产品id
	DeviceName string `gorm:"column:device_name;index:idx_dm_device_module_version_pd_dn;uniqueIndex:idx_dm_device_module_version_product_id_deviceName;type:varchar(100);NOT NULL"` // 设备名称
	ModuleCode string `gorm:"column:module_code;type:varchar(64);uniqueIndex:idx_dm_device_module_version_product_id_deviceName"`                                                    // 固件名称
	Version    string `gorm:"column:version;type:varchar(64);NOT NULL;uniqueIndex:idx_dm_device_module_version_product_id_deviceName"`                                               // 固件版本
	stores.NoDelTime
}

func (m *DmDeviceModuleVersion) TableName() string {
	return "dm_device_module_version"
}

// 用户配置表
type DmDeviceProfile struct {
	ID         int64               `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	TenantCode dataType.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL;uniqueIndex:idx_dm_device_profile_tc_un;"`                                   // 租户编码
	ProductID  string              `gorm:"column:product_id;type:varchar(100);index:idx_dm_device_profile_pd_dn;uniqueIndex:idx_dm_device_profile_tc_un;NOT NULL"`  // 产品id
	DeviceName string              `gorm:"column:device_name;index:idx_dm_device_profile_pd_dn;uniqueIndex:idx_dm_device_profile_tc_un;type:varchar(100);NOT NULL"` // 设备名称
	Code       string              `gorm:"column:code;type:VARCHAR(50);uniqueIndex:idx_dm_device_profile_tc_un;NOT NULL"`                                           //配置code
	Params     string              `gorm:"column:params;type:text;NOT NULL"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_dm_device_profile_tc_un;"`
}

func (m *DmDeviceProfile) TableName() string {
	return "dm_device_profile"
}

// 产品信息表
type DmProductInfo struct {
	ID               int64                 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	ProductID        string                `gorm:"column:product_id;type:varchar(100);uniqueIndex:idx_dm_product_info_pd;NOT NULL"` // 产品id
	ProductName      string                `gorm:"column:product_name;type:varchar(100);NOT NULL"`                                  // 产品名称
	ProductImg       string                `gorm:"column:product_img;type:varchar(200)"`                                            // 产品图片
	ProductType      int64                 `gorm:"column:product_type;type:smallint;default:1"`                                     // 产品状态:1:开发中,2:审核中,3:已发布
	AuthMode         int64                 `gorm:"column:auth_mode;type:smallint;default:1"`                                        // 认证方式:1:账密认证,2:秘钥认证
	DeviceType       int64                 `gorm:"column:device_type;index;type:smallint;default:1"`                                // 设备类型:1:设备,2:网关,3:子设备
	CategoryID       int64                 `gorm:"column:category_id;type:integer;default:2"`                                       // 产品品类 2:未分类
	NetType          int64                 `gorm:"column:net_type;type:smallint;default:1"`                                         // 通讯方式:1:其他,2:wi-fi,3:2G/3G/4G,4:5G,5:BLE,6:LoRaWAN
	ProtocolCode     string                `gorm:"column:protocol_code;type:varchar(100);default:'urMqtt'"`                         // 协议code,默认urMqtt  urMqtt,urHttp,wumei,aliyun,huaweiyun,tuya
	SubProtocolCode  string                `gorm:"column:sub_protocol_code;type:varchar(100);default:''"`                           //子协议,主协议和子协议传输类型必须不相同, 设备控制下发只会发送给主协议, 当设备是音视频设备但是控制协议需要单独走的时候就可以把主协议定义为普通协议,子协议定义为音视频协议,这样就能实现音视频走音视频协议,控制走子协议
	AutoRegister     def.AutoReg           `gorm:"column:auto_register;type:smallint;default:1"`                                    // 动态注册:1:关闭,2:打开,3:打开并自动创建设备
	OnlineHandle     int64                 `gorm:"column:online_handle;type:smallint;default:1"`                                    //在线处理:1: 自动 2: 永远在线
	DeviceSchemaMode int64                 `gorm:"column:device_schema_mode;type:smallint;default:1"`                               // 设备物模型模式:1:关闭,2:设备自动创建3: 设备自动创建及上报无定义自动创建
	BindLevel        product.BindLevel     `gorm:"column:bind_level;type:smallint;default:1"`                                       // 绑定级别: 1:强绑定(默认,只有用户解绑之后才能绑定) 2:中绑定(可以通过token强制绑定设备) 3:弱绑定(app可以内部解绑被绑定的设备)
	Secret           string                `gorm:"column:secret;type:varchar(50)"`                                                  // 动态注册产品秘钥
	Desc             string                `gorm:"column:desc;type:varchar(200)"`                                                   // 描述
	TrialTime        int64                 `gorm:"column:trial_time"`                                                               //试用时间(单位为天,为0不限制)
	Status           devices.ProductStatus `gorm:"column:status;type:smallint;default:1"`
	SceneMode        string                `gorm:"column:scene_mode;type:varchar(20);default:rw"`                   // 场景模式 读写类型: r(只读) rw(可读可写) none(不参与场景)
	Tags             map[string]string     `gorm:"column:tags;type:json;serializer:json;default:'{}'"`              // 产品标签
	ProtocolConf     map[string]string     `gorm:"column:protocol_conf;type:json;serializer:json;default:'{}'"`     // 自定义协议配置
	SubProtocolConf  map[string]string     `gorm:"column:sub_protocol_conf;type:json;serializer:json;default:'{}'"` // 子模块自定义协议配置
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:idx_dm_product_info_pd"`
	Category    *DmProductCategory `gorm:"foreignKey:ID;references:CategoryID"`       // 添加外键
	Config      *DmProductConfig   `gorm:"foreignKey:ProductID;references:ProductID"` // 添加外键
	Protocol    *DmProtocolInfo    `gorm:"foreignKey:Code;references:ProtocolCode"`
	SubProtocol *DmProtocolInfo    `gorm:"foreignKey:Code;references:SubProtocolCode"`
	//Devices []*DmDeviceInfo    `gorm:"foreignKey:ProductID;references:ProductID"` // 添加外键
	CustomUi map[string]*ProductCustomUi `gorm:"column:customUi;type:json;serializer:json;NOT NULL;default:'{}'"`
}

type ProductCustomUi struct {
	Type    string `json:"type"` //detail(设备详情) ,setNet(配置网络)
	Path    string `json:"path"`
	Version int64  `json:"version"`
}

func (m *DmProductInfo) TableName() string {
	return "dm_product_info"
}

type DmProductID struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
}

func (m *DmProductID) TableName() string {
	return "dm_product_id"
}

type DmProductConfig struct {
	ID        int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	ProductID string `gorm:"column:product_id;type:varchar(100);uniqueIndex:idx_dm_product_config_pd;NOT NULL"` // 产品id
	DevInit   `gorm:"embedded;embeddedPrefix:dev_init_"`                                                        //设备初始化配置
	stores.NoDelTime
	Info        *DmProductInfo     `gorm:"foreignKey:product_id;references:product_id"`
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:idx_dm_product_config_pd"`
}

func (m *DmProductConfig) TableName() string {
	return "dm_product_config"
}

type DevInit struct {
	TenantCode string `gorm:"column:tenant_code;type:VARCHAR(50);default:''"` //设备创建默认分配给的租户
	ProjectID  int64  `gorm:"column:project_id;type:bigint;default:0;"`       //设备创建默认分配给的项目
	AreaID     int64  `gorm:"column:area_id;type:bigint;default:0;"`          //设备创建默认分配给的区域ID
	DeptID     int64  `gorm:"column:dept_id;type:bigint;default:0;"`          //设备创建默认分配给的部门
}

type DmProductCategory struct {
	ID          int64                      `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	Name        string                     `gorm:"column:name;uniqueIndex:idx_dm_product_category_pn;type:varchar(100);NOT NULL"` // 产品品类名称
	Desc        string                     `gorm:"column:desc;type:varchar(200)"`                                                 // 描述
	HeadImg     string                     `gorm:"column:head_img;type:varchar(200)"`                                             // 图片
	ParentID    int64                      `gorm:"column:parent_id;type:bigint;NOT NULL"`                                         // 上级区域ID(雪花ID)
	IDPath      string                     `gorm:"column:id_path;type:varchar(100);NOT NULL"`                                     // 1-2-3-的格式记录顶级区域到当前区域的路径
	IsLeaf      int64                      `gorm:"column:is_leaf;type:bigint;default:1;NOT NULL"`                                 //是否是叶子节点
	DeviceCount int64                      `gorm:"column:device_count;type:bigint;default:0"`                                     //该产品品类下的设备数量统计
	Children    []*DmProductCategory       `gorm:"foreignKey:ParentID;references:ID;"`
	Schemas     []*DmProductCategorySchema `gorm:"foreignKey:ProductCategoryID;references:ID;"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_dm_product_category_pn"`
}

func (m *DmProductCategory) TableName() string {
	return "dm_product_category"
}

type DmProductCategorySchema struct {
	ID                int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	ProductCategoryID int64  `gorm:"column:product_category_id;uniqueIndex:idx_dm_product_category_schema_pn;type:bigint;NOT NULL"` // 产品品类id
	Identifier        string `gorm:"column:identifier;type:varchar(200);uniqueIndex:idx_dm_product_category_schema_pn"`             // 标识符的id
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_dm_product_category_schema_pn"`
}

func (m *DmProductCategorySchema) TableName() string {
	return "dm_product_category_schema"
}

// 自定义协议表
type DmProtocolInfo struct {
	ID                  int64                 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	Name                string                `gorm:"column:name;uniqueIndex:idx_dm_protocol_info_pn;type:varchar(100);NOT NULL"` // 协议名称
	Code                string                `gorm:"column:code;uniqueIndex:idx_dm_protocol_info_pc;type:varchar(100);NOT NULL"` // iThings,iThings-thingsboard,wumei,aliyun,huaweiyun,tuya
	Type                protocols.Type        `gorm:"column:type;type:varchar(100);default:normal"`                               //协议类型
	TransProtocol       string                `gorm:"column:trans_protocol;type:varchar(100);default:mqtt"`                       // 传输协议: mqtt,tcp,udp
	Desc                string                `gorm:"column:desc;type:varchar(2000)"`                                             // 描述
	IsEnableSyncProduct int64                 `gorm:"column:is_enable_sync_product;type:bigint;default:2"`                        //是否支持主动同步产品信息
	IsEnableSyncDevice  int64                 `gorm:"column:is_enable_sync_device;type:bigint;default:2"`                         //是否支持主动同步设备信息
	DeviceFields        protocol.ConfigFields `gorm:"column:device_fields;type:json;serializer:json;NOT NULL;default:'[]'"`       // 协议需要的设备配置字段列表,没有可以不传
	ProductFields       protocol.ConfigFields `gorm:"column:product_fields;type:json;serializer:json;NOT NULL;default:'[]'"`      // 协议需要的产品配置字段列表,没有可以不传
	ConfigFields        protocol.ConfigFields `gorm:"column:config_fields;type:json;serializer:json;NOT NULL;default:'[]'"`       // 需要配置的字段列表,没有可以不传
	ConfigInfos         protocol.ConfigInfos  `gorm:"column:config_infos;type:json;serializer:json;NOT NULL;default:'[]'"`        // 配置列表
	Endpoints           []string              `gorm:"column:endpoints;type:json;serializer:json;NOT NULL;default:'[]'"`           // 协议端点,如果填写了优先使用该字段
	EtcdKey             string                `gorm:"column:etcd_key;type:varchar(200);default:null"`                             // 服务etcd发现的key etcd key
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_dm_protocol_info_pc;uniqueIndex:idx_dm_protocol_info_pn"`
}

func (m *DmProtocolInfo) TableName() string {
	return "dm_protocol_info"
}

type DmProtocolService struct {
	ID     int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	Code   string `gorm:"column:code;uniqueIndex:idx_dm_protocol_service_pc;type:varchar(100);default:iThings"` //
	Ip     string `gorm:"column:ip;uniqueIndex:idx_dm_protocol_service_pc;type:varchar(100);not null"`          //
	Port   int64  `gorm:"column:port;uniqueIndex:idx_dm_protocol_service_pc;type:varchar(200)"`                 //
	Status int64  `gorm:"column:status;type:bigint;default:2"`
	stores.NoDelTime
}

func (m *DmProtocolService) TableName() string {
	return "dm_protocol_service"
}

// 协议插件
type DmProtocolScript struct {
	ID   int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	Name string `gorm:"column:name;uniqueIndex:idx_dm_protocol_script_name;type:varchar(100);not null"` //转换名称
	//ProductIDs    []string          `gorm:"column:product_ids;type:json;serializer:json;default:'[]'"` // 产品id
	//Devices       []devices.Core    `gorm:"column:devices;type:json;serializer:json;default:'[]'"`     //设备触发
	TriggerDir    protocol.TriggerDir   `gorm:"column:trigger_dir;not null"`                        //up down
	TriggerTimer  protocol.TriggerTimer `gorm:"column:trigger_timer;"`                              //收到前处理before after
	TriggerHandle devices.MsgHandle     `gorm:"column:trigger_handle;type:varchar(100);default:''"` //对应 mqtt topic的第一个 thing ota config 等等
	TriggerType   string                `gorm:"column:trigger_type;type:varchar(100);default:''"`   // 操作类型 从topic中提取 物模型下就是   property属性 event事件 action行为
	ScriptLang    int64                 `gorm:"column:script_lang;type:smallint;default:1"`         // 脚本语言类型 1:go
	Script        string                `gorm:"column:script;type:text"`                            // 协议转换脚本
	Desc          string                `gorm:"column:desc;type:varchar(200)"`                      // 描述
	Status        def.Bool              `gorm:"column:status;default:1"`                            //状态:是否启用
	stores.NoDelTime
	DeletedTime stores.DeletedTime        `gorm:"column:deleted_time;default:0;uniqueIndex:idx_dm_protocol_script_name"`
	Devices     []*DmProtocolScriptDevice `gorm:"foreignKey:ScriptID;references:ID"`
}

func (m *DmProtocolScript) TableName() string {
	return "dm_protocol_script"
}

// 协议插件
type DmProtocolScriptDevice struct {
	ID         int64               `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	TriggerSrc protocol.TriggerSrc `gorm:"column:trigger_src;"`                             //product:1 device:2
	ProductID  string              `gorm:"column:product_id;type:varchar(100);not null"`    // 产品id
	DeviceName string              `gorm:"column:device_name;type:varchar(100);default:''"` //设备
	ScriptID   int64               `gorm:"column:script_id;type:bigint"`
	Priority   int64               `gorm:"column:priority;default:1"` //执行优先级
	Status     def.Bool            `gorm:"column:status;default:1"`   //状态:是否启用
	Script     *DmProtocolScript   `gorm:"foreignKey:ID;references:ScriptID"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0"`
}

func (m *DmProtocolScriptDevice) TableName() string {
	return "dm_protocol_script_device"
}

// 产品自定义协议表
type DmProductCustom struct {
	ID              int64                        `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	ProductID       string                       `gorm:"column:product_id;uniqueIndex:idx_dm_product_custom_pn;type:varchar(100);NOT NULL"` // 产品id
	ScriptLang      int64                        `gorm:"column:script_lang;type:smallint;default:1"`                                        // 脚本语言类型 1:JavaScript 2:lua 3:python
	CustomTopics    []*productCustom.CustomTopic `gorm:"column:custom_topics;type:json"`                                                    // 自定义topic数组
	TransformScript string                       `gorm:"column:transform_script;type:text"`                                                 // 协议转换脚本
	LoginAuthScript string                       `gorm:"column:login_auth_script;type:text"`                                                // 登录认证脚本
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_dm_product_custom_pn"`
	ProductInfo *DmProductInfo     `gorm:"foreignKey:ProductID;references:ProductID"`
}

func (m *DmProductCustom) TableName() string {
	return "dm_product_custom"
}

// 产品物模型表
type DmProductSchema struct {
	ID         int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	ProductID  string `gorm:"column:product_id;uniqueIndex:idx_dm_product_schema_identifier;index:product_id_type;type:varchar(100);NOT NULL"` // 产品id
	Identifier string `gorm:"column:identifier;uniqueIndex:idx_dm_product_schema_identifier;type:varchar(100);NOT NULL"`                       // 标识符
	DmSchemaCore
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_dm_product_schema_identifier"`
	ProductInfo *DmProductInfo     `gorm:"foreignKey:ProductID;references:ProductID"`
}

func (m *DmProductSchema) TableName() string {
	return "dm_product_schema"
}

// 产品物模型表
type DmDeviceSchema struct {
	ID         int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	ProductID  string `gorm:"column:product_id;uniqueIndex:idx_dm_schema_info_identifier;index:product_id_type;type:varchar(100);NOT NULL"`      // 产品id
	DeviceName string `gorm:"column:device_name;uniqueIndex:idx_dm_schema_info_identifier;index:product_id_type;type:varchar(100);default:null"` // 产品id
	Identifier string `gorm:"column:identifier;uniqueIndex:idx_dm_schema_info_identifier;type:varchar(100);NOT NULL"`                            // 标识符
	DmSchemaCore
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_dm_schema_info_identifier"`
}

func (m *DmDeviceSchema) TableName() string {
	return "dm_schema_info"
}

// 产品物模型表
type DmSchemaInfo struct {
	ID         int64          `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	ProductID  string         `gorm:"column:product_id;uniqueIndex:idx_dm_schema_info_identifier;index:product_id_type;type:varchar(100);NOT NULL"`      // 产品id
	DeviceName sql.NullString `gorm:"column:device_name;uniqueIndex:idx_dm_schema_info_identifier;index:product_id_type;type:varchar(100);default:null"` // 产品id
	Identifier string         `gorm:"column:identifier;uniqueIndex:idx_dm_schema_info_identifier;type:varchar(100);NOT NULL"`                            // 标识符
	DmSchemaCore
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_dm_schema_info_identifier"`
}

func (m *DmSchemaInfo) TableName() string {
	return "dm_schema_info"
}

type DmSchemaCore struct {
	Tag               schema.Tag        `gorm:"column:tag;type:smallint;default:1"`                        // 物模型标签 1:自定义 2:可选 3:必选  必选不可删除
	Type              int64             `gorm:"column:type;index:product_id_type;type:smallint;default:1"` // 物模型类型 1:property属性 2:event事件 3:action行为
	ExtendConfig      string            `gorm:"column:extend_config;type:text"`                            //拓展参数
	Required          int64             `gorm:"column:required;type:smallint;default:2"`                   // 是否必须,1是 2否
	Name              string            `gorm:"column:name;type:varchar(100);NOT NULL"`                    // 功能名称
	Desc              string            `gorm:"column:desc;type:varchar(200)"`                             // 描述
	IsCanSceneLinkage int64             `gorm:"column:is_can_scene_linkage;type:smallint;default:1"`       // 是否放到场景联动中
	FuncGroup         int64             `gorm:"column:func_group;type:smallint;default:1"`                 // 功能分类: 1:普通功能 2:系统功能
	ControlMode       int64             `gorm:"column:control_mode;type:smallint;default:1"`               //控制模式: 1: 可以群控,可以单控  2:只能单控
	UserPerm          int64             `gorm:"column:user_auth;type:smallint;default:3"`                  //用户权限操作: 1:r(只读) 3:rw(可读可写)
	RecordMode        schema.RecordMode `gorm:"column:record_mode;type:smallint;default:1"`                //1(默认) 记录历史记录 2 只记录差异值 3 不记录历史记录
	IsPassword        int64             `gorm:"column:is_password;type:smallint;default:2"`                //是否是密码类型
	Affordance        string            `gorm:"column:affordance;type:json;NOT NULL"`                      // 各类型的自定义功能定义
	Order             int64             `gorm:"column:order;type:BIGINT;default:1;NOT NULL"`               // 左侧table排序序号
}

// 通用物模型表
type DmCommonSchema struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	DmSchemaCore
	Identifier string `gorm:"column:identifier;uniqueIndex:idx_dm_common_schema_identifier;type:varchar(100);NOT NULL"` // 标识符
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_dm_common_schema_identifier"`
}

func (m *DmCommonSchema) TableName() string {
	return "dm_common_schema"
}

// 设备分组信息表
type DmGroupInfo struct {
	ID         int64               `gorm:"column:id;primary_key;AUTO_INCREMENT;type:bigint"`                                                 // 分组ID
	TenantCode dataType.TenantCode `gorm:"column:tenant_code;uniqueIndex:idx_dm_group_info_tc_ac;default:default;type:VARCHAR(50);NOT NULL"` // 租户编码
	Purpose    string              `gorm:"column:purpose;uniqueIndex:idx_dm_group_info_tc_ac;type:varchar(100);default:'default'"`           //用途 不填默认为default
	ProjectID  dataType.ProjectID  `gorm:"column:project_id;uniqueIndex:idx_dm_group_info_tc_ac;type:bigint;default:2;NOT NULL"`             // 项目ID(雪花ID)
	AreaID     dataType.AreaID     `gorm:"column:area_id;uniqueIndex:idx_dm_group_info_tc_ac;type:bigint;default:2;NOT NULL"`                // 项目区域ID(雪花ID)
	ParentID   int64               `gorm:"column:parent_id;type:bigint;default:0;NOT NULL"`                                                  // 父组ID 0-根组
	IDPath     string              `gorm:"column:id_path;type:varchar(100);NOT NULL"`                                                        // 1-2-3-的格式记录顶级区域到当前区域的路径
	IsLeaf     int64               `gorm:"column:is_leaf;type:bigint;default:1;NOT NULL"`                                                    //是否是叶子节点
	ProductID  string              `gorm:"column:product_id;type:varchar(100);default:'';NOT NULL"`                                          // 产品id,为空则不限定分组内的产品类型
	Name       string              `gorm:"column:name;uniqueIndex:idx_dm_group_info_tc_ac;default:'';type:varchar(100);NOT NULL"`            // 分组名称
	Desc       string              `gorm:"column:desc;type:varchar(200);default:''"`                                                         // 描述
	Tags       map[string]string   `gorm:"column:tags;type:json;serializer:json;NOT NULL;default:'{}'"`                                      // 分组标签
	Files      map[string]string   `gorm:"column:files;type:json;serializer:json;NOT NULL;default:'{}'"`                                     // 分组文件

	stores.NoDelTime
	DeviceCount int64              `gorm:"column:device_count;type:bigint;default:0;"` //设备数量统计
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_dm_group_info_tc_ac"`
	ProductInfo *DmProductInfo     `gorm:"foreignKey:ProductID;references:ProductID"`
}

func (m *DmGroupInfo) TableName() string {
	return "dm_group_info"
}

// 分组与设备关系表
type DmGroupDevice struct {
	ID          int64               `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	TenantCode  dataType.TenantCode `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL"`                                                                 // 租户编码
	ProjectID   dataType.ProjectID  `gorm:"column:project_id;index;type:bigint;default:0;NOT NULL"`                                                             // 项目ID(雪花ID)
	AreaID      dataType.AreaID     `gorm:"column:area_id;index:project_id_area_id;type:bigint;default:0;NOT NULL"`                                             // 项目区域ID(雪花ID)
	GroupID     int64               `gorm:"column:group_id;uniqueIndex:idx_dm_group_device_group_id_product_id_device_name;type:bigint;NOT NULL"`               // 分组ID
	GroupIDPath string              `gorm:"column:group_id_path;type:varchar(100);default:''"`                                                                  // 1-2-3-的格式记录顶级区域到当前区域的路径
	Purpose     string              `gorm:"column:purpose;uniqueIndex:idx_dm_group_device_group_id_product_id_device_name;type:varchar(100);default:'default'"` //用途 不填默认为default
	ProductID   string              `gorm:"column:product_id;uniqueIndex:idx_dm_group_device_group_id_product_id_device_name;type:varchar(100);NOT NULL"`       // 产品id
	DeviceName  string              `gorm:"column:device_name;uniqueIndex:idx_dm_group_device_group_id_product_id_device_name;type:varchar(100);NOT NULL"`      // 设备名称
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_dm_group_device_group_id_product_id_device_name"`
	ProductInfo *DmProductInfo     `gorm:"foreignKey:ProductID;references:ProductID"`
	Device      *DmDeviceInfo      `gorm:"foreignKey:ProductID,DeviceName;references:ProductID,DeviceName"`
	Group       *DmGroupInfo       `gorm:"foreignKey:ID;references:GroupID"`
}

func (m *DmGroupDevice) TableName() string {
	return "dm_group_device"
}

// 网关与子设备关系表
type DmGatewayDevice struct {
	ID                int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	GatewayProductID  string `gorm:"column:gateway_product_id;type:varchar(100);NOT NULL"`                                          // 网关产品id
	GatewayDeviceName string `gorm:"column:gateway_device_name;type:varchar(100);NOT NULL"`                                         // 网关设备名称
	ProductID         string `gorm:"column:product_id;uniqueIndex:idx_dm_gateway_device_gpi_pdn_pi_dn;type:varchar(100);NOT NULL"`  // 子设备产品id
	DeviceName        string `gorm:"column:device_name;uniqueIndex:idx_dm_gateway_device_gpi_pdn_pi_dn;type:varchar(100);NOT NULL"` // 子设备名称
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_dm_gateway_device_gpi_pdn_pi_dn"`
	Device      *DmDeviceInfo      `gorm:"foreignKey:ProductID,DeviceName;references:ProductID,DeviceName"`
	Gateway     *DmDeviceInfo      `gorm:"foreignKey:ProductID,DeviceName;references:GatewayProductID,GatewayDeviceName"`
}

func (m *DmGatewayDevice) TableName() string {
	return "dm_gateway_device"
}

// 产品远程配置表
type DmProductRemoteConfig struct {
	ID        int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	ProductID string `gorm:"column:product_id;uniqueIndex:idx_dm_product_remote_config_pn;type:varchar(100);NOT NULL"` // 产品id
	Content   string `gorm:"column:content;type:json;NOT NULL"`                                                        // 配置内容
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_dm_product_remote_config_pn"`
	ProductInfo *DmProductInfo     `gorm:"foreignKey:ProductID;references:ProductID"`
}

func (m *DmProductRemoteConfig) TableName() string {
	return "dm_product_remote_config"
}

// 设备影子表
type DmDeviceShadow struct {
	ID                int64        `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
	ProductID         string       `gorm:"column:product_id;uniqueIndex:idx_dm_device_shadow_pi_dn_di;type:varchar(100);NOT NULL"`
	DeviceName        string       `gorm:"column:device_name;uniqueIndex:idx_dm_device_shadow_pi_dn_di;type:VARCHAR(100);NOT NULL"`
	DataID            string       `gorm:"column:data_id;uniqueIndex:idx_dm_device_shadow_pi_dn_di;type:VARCHAR(100);NOT NULL"`
	Value             string       `gorm:"column:value;type:VARCHAR(100);default:NULL"`
	UpdatedDeviceTime sql.NullTime `gorm:"column:updated_device_time;"`
	stores.OnlyTime
}

func (m *DmDeviceShadow) TableName() string {
	return "dm_device_shadow"
}

type DmManufacturerInfo struct {
	ID    int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	Name  string `gorm:"column:name;uniqueIndex:idx_dm_manufacturer_info_pn;type:varchar(100);NOT NULL"` // 协议名称
	Desc  string `gorm:"column:desc;type:varchar(200)"`                                                  // 描述
	Phone string `gorm:"column:phone;type:varchar(200);default:null"`                                    //联系电话
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_dm_manufacturer_info_pn"`
}

func (m *DmManufacturerInfo) TableName() string {
	return "dm_manufacturer_info"
}
