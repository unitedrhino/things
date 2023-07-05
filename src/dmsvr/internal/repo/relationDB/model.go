package relationDB

import (
	"database/sql"
	"github.com/i-Things/things/shared/store"
	"time"
)

// 产品信息表

type DmProductInfo struct {
	ProductID    string `gorm:"column:productID;type:char(11);NOT NULL"`       // 产品id
	ProductName  string `gorm:"column:productName;type:varchar(100);NOT NULL"` // 产品名称
	ProductImg   string `gorm:"column:productImg;type:varchar(200)"`           // 产品图片
	ProductType  int64  `gorm:"column:productType;type:tinyint(1);default:1"`  // 产品状态:1:开发中,2:审核中,3:已发布
	AuthMode     int64  `gorm:"column:authMode;type:tinyint(1);default:1"`     // 认证方式:1:账密认证,2:秘钥认证
	DeviceType   int64  `gorm:"column:deviceType;type:tinyint(1);default:1"`   // 设备类型:1:设备,2:网关,3:子设备
	CategoryID   int64  `gorm:"column:categoryID;type:int(10);default:1"`      // 产品品类
	NetType      int64  `gorm:"column:netType;type:tinyint(1);default:1"`      // 通讯方式:1:其他,2:wi-fi,3:2G/3G/4G,4:5G,5:BLE,6:LoRaWAN
	DataProto    int64  `gorm:"column:dataProto;type:tinyint(1);default:1"`    // 数据协议:1:自定义,2:数据模板
	AutoRegister int64  `gorm:"column:autoRegister;type:tinyint(1);default:1"` // 动态注册:1:关闭,2:打开,3:打开并自动创建设备
	Secret       string `gorm:"column:secret;type:varchar(50)"`                // 动态注册产品秘钥
	Desc         string `gorm:"column:desc;type:varchar(200)"`                 // 描述
	DevStatus    string `gorm:"column:devStatus;type:varchar(20);NOT NULL"`    // 产品状态
	Tags         string `gorm:"column:tags;type:json;NOT NULL"`                // 产品标签
	store.Time
}

func (m *DmProductInfo) TableName() string {
	return "dm_product_info"
}

// 产品自定义协议表
type DmProductCustom struct {
	ID              int64          `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"`
	ProductID       string         `gorm:"column:productID;type:char(11);NOT NULL"`     // 产品id
	ScriptLang      int            `gorm:"column:scriptLang;type:tinyint(1);default:1"` // 脚本语言类型 1:JavaScript 2:lua 3:python
	CustomTopic     sql.NullString `gorm:"column:customTopic;type:json"`                // 自定义topic数组
	TransformScript sql.NullString `gorm:"column:transformScript;type:text"`            // 协议转换脚本
	LoginAuthScript sql.NullString `gorm:"column:loginAuthScript;type:text"`            // 登录认证脚本
	store.Time
}

func (m *DmProductCustom) TableName() string {
	return "dm_product_custom"
}

// 产品物模型表
type DmProductSchema struct {
	ID         int64  `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"`
	ProductID  string `gorm:"column:productID;type:char(11);NOT NULL"`      // 产品id
	Tag        int    `gorm:"column:tag;type:tinyint(1);default:1"`         // 物模型标签 1:自定义 2:可选 3:必选  必选不可删除
	Type       int    `gorm:"column:type;type:tinyint(1);default:1"`        // 物模型类型 1:property属性 2:event事件 3:action行为
	Identifier string `gorm:"column:identifier;type:varchar(100);NOT NULL"` // 标识符
	Name       string `gorm:"column:name;type:varchar(100);NOT NULL"`       // 功能名称
	Desc       string `gorm:"column:desc;type:varchar(200)"`                // 描述
	Required   int    `gorm:"column:required;type:tinyint(1);default:2"`    // 是否必须,1是 2否
	Affordance string `gorm:"column:affordance;type:json;NOT NULL"`         // 各类型的自定义功能定义
	store.Time
}

func (m *DmProductSchema) TableName() string {
	return "dm_product_schema"
}

// 设备信息表
type DmDeviceInfo struct {
	ID             int64       `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"`
	ProjectID      int64       `gorm:"column:projectID;type:bigint(20);default:0;NOT NULL"`      // 项目ID(雪花ID)
	AreaID         int64       `gorm:"column:areaID;type:bigint(20);default:0;NOT NULL"`         // 项目区域ID(雪花ID)
	ProductID      string      `gorm:"column:productID;type:char(11);NOT NULL"`                  // 产品id
	DeviceName     string      `gorm:"column:deviceName;type:varchar(100);NOT NULL"`             // 设备名称
	DeviceAlias    string      `gorm:"column:deviceAlias;type:varchar(100);NOT NULL"`            // 设备别名
	Position       store.Point `gorm:"column:position;type:point;NOT NULL"`                      // 设备的位置(默认百度坐标系BD09)
	Secret         string      `gorm:"column:secret;type:varchar(50);NOT NULL"`                  // 设备秘钥
	Cert           string      `gorm:"column:cert;type:varchar(512);NOT NULL"`                   // 设备证书
	Imei           string      `gorm:"column:imei;type:varchar(15);NOT NULL"`                    // IMEI号信息
	Mac            string      `gorm:"column:mac;type:varchar(17);NOT NULL"`                     // MAC号信息
	Version        string      `gorm:"column:version;type:varchar(64);NOT NULL"`                 // 固件版本
	HardInfo       string      `gorm:"column:hardInfo;type:varchar(64);NOT NULL"`                // 模组硬件型号
	SoftInfo       string      `gorm:"column:softInfo;type:varchar(64);NOT NULL"`                // 模组软件版本
	MobileOperator int         `gorm:"column:mobileOperator;type:tinyint(1);default:1;NOT NULL"` // 移动运营商:1)移动 2)联通 3)电信 4)广电
	Phone          string      `gorm:"column:phone;type:varchar(20)"`                            // 手机号
	Iccid          string      `gorm:"column:iccid;type:varchar(20)"`                            // SIM卡卡号
	Address        string      `gorm:"column:address;type:varchar(512);NOT NULL"`                // 所在地址
	Tags           string      `gorm:"column:tags;type:json;NOT NULL"`                           // 设备标签
	UserID         int64       `gorm:"column:userID;type:bigint(20);NOT NULL"`                   // 所属用户id
	IsOnline       int         `gorm:"column:isOnline;type:tinyint(1);default:2;NOT NULL"`       // 是否在线,1是2否
	FirstLogin     time.Time   `gorm:"column:firstLogin;type:datetime"`                          // 激活时间
	LastLogin      time.Time   `gorm:"column:lastLogin;type:datetime"`                           // 最后上线时间
	LogLevel       int         `gorm:"column:logLevel;type:tinyint(1);default:1;NOT NULL"`       // 日志级别:1)关闭 2)错误 3)告警 4)信息 5)调试
	store.Time
}

func (m *DmDeviceInfo) TableName() string {
	return "dm_device_info"
}

// 设备分组信息表
type DmGroupInfo struct {
	GroupID   int64  `gorm:"column:groupID;type:bigint(20)"`                      // 分组ID
	ParentID  int64  `gorm:"column:parentID;type:bigint(20);default:0;NOT NULL"`  // 父组ID 0-根组
	ProjectID int64  `gorm:"column:projectID;type:bigint(20);default:0;NOT NULL"` // 项目ID(雪花ID)
	ProductID string `gorm:"column:productID;type:char(11);NOT NULL"`             // 产品id,为空则不限定分组内的产品类型
	GroupName string `gorm:"column:groupName;type:varchar(100);NOT NULL"`         // 分组名称
	Desc      string `gorm:"column:desc;type:varchar(200)"`                       // 描述
	Tags      string `gorm:"column:tags;type:json;NOT NULL"`                      // 分组标签
	store.Time
}

func (m *DmGroupInfo) TableName() string {
	return "dm_group_info"
}

// 分组与设备关系表
type DmGroupDevice struct {
	ID         int64  `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"`
	GroupID    int64  `gorm:"column:groupID;type:bigint(20);NOT NULL"`             // 分组ID
	ProjectID  int64  `gorm:"column:projectID;type:bigint(20);default:0;NOT NULL"` // 项目ID(雪花ID)
	ProductID  string `gorm:"column:productID;type:char(11);NOT NULL"`             // 产品id
	DeviceName string `gorm:"column:deviceName;type:varchar(100);NOT NULL"`        // 设备名称
	store.Time
}

func (m *DmGroupDevice) TableName() string {
	return "dm_group_device"
}

// 网关与子设备关系表
type DmGatewayDevice struct {
	ID                int64  `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"`
	GatewayProductID  string `gorm:"column:gatewayProductID;type:char(11);NOT NULL"`      // 网关产品id
	GatewayDeviceName string `gorm:"column:gatewayDeviceName;type:varchar(100);NOT NULL"` // 网关设备名称
	ProductID         string `gorm:"column:productID;type:char(11);NOT NULL"`             // 子设备产品id
	DeviceName        string `gorm:"column:deviceName;type:varchar(100);NOT NULL"`        // 子设备名称
	store.Time
}

func (m *DmGatewayDevice) TableName() string {
	return "dm_gateway_device"
}

// 产品远程配置表
type DmProductRemoteConfig struct {
	ID        int64  `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"`
	ProductID string `gorm:"column:productID;type:char(11);NOT NULL"` // 产品id
	Content   string `gorm:"column:content;type:json;NOT NULL"`       // 配置内容
	store.Time
}

func (m *DmProductRemoteConfig) TableName() string {
	return "dm_product_remote_config"
}
