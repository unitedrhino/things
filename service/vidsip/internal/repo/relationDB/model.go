package relationDB

import (
	"gitee.com/i-Things/core/shared/stores"
	sip2 "github.com/i-Things/things/service/vidsip/gosip/sip"
	"net"
)

/********************************** GB28181 数据 ***********************************/
type SipChannels struct {
	//ID int64 `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	ChannelID    string `gorm:"column:channel_id;primary_key;type:char(20);NOT NULL" xml:"DeviceID" json:"channelid"` // ChannelID 通道编码
	DeviceID     string `gorm:"column:device_id;type:char(20);NOT NULL" xml:"-" json:"deviceid"`                      // DeviceID 设备编号
	VidmgrID     string `gorm:"column:vidmgr_id;type:char(11);NOT NULL"`                                              //用于统计通道数量                                              // 流服务ID
	Stream       string `gorm:"column:stream;type:char(20)"`                                                          //用于申请新的流
	Memo         string `gorm:"column:memo" json:"memo"`                                                              // Memo 备注（用来标示通道信息）
	Name         string `gorm:"column:name" xml:"Name" json:"name"  `                                                 // Name 通道名称（设备端设置名称）
	Manufacturer string `gorm:"column:manufacturer" xml:"Manufacturer" json:"manufacturer"`
	Model        string `gorm:"column:model"  xml:"Model" json:"model"`
	Owner        string `gorm:"column:owner" xml:"Owner"  json:"owner"`
	CivilCode    string `gorm:"column:civilcode" xml:"CivilCode" json:"civilcode"`

	IsPlay bool `gorm:"column:isplay;type:smallint;default:0;NOT NULL"` //是否是播放状态

	Address     string `gorm:"column:address" xml:"Address"  json:"address"` // Address ip地址
	Parental    int32  `gorm:"column:parental" xml:"Parental"  json:"parental"`
	SafetyWay   int32  `gorm:"column:safetyway" xml:"SafetyWay"  json:"safetyway"`
	RegisterWay int32  `gorm:"column:registerway" xml:"RegisterWay"  json:"registerway"`
	Secrecy     int32  `gorm:"column:secrecy" xml:"Secrecy" json:"secrecy"`
	// Status 状态:  unactivated:0  on:1 在线  offline:2 下线   // Status 状态  on 在线
	Status string `gorm:"column:status;type:smallint;default:0" xml:"Status"  json:"status"`
	// Active 最后活跃时间
	LastLogin int64  `gorm:"column:last_login" json:"active" ` // 最后登录时间
	URIStr    string `gorm:"column:uri" json:"uri" `
	// 视频编码格式
	VF string `gorm:"column:vf" json:"vf"`
	// 视频高
	Height int32 `gorm:"column:height" json:"height"`
	// 视频宽
	Width int32 `gorm:"column:width" json:"width"`
	// 视频FPS
	FPS int32 `gorm:"column:fps" json:"fps"`
	//  pull 媒体服务器主动拉流，push 监控设备主动推流
	StreamType string `gorm:"column:streamtype" json:"streamtype"`
	// streamtype=pull时，拉流地址
	URL string `gorm:"column:url" json:"url"`
	stores.Time
	//----
	Taddr *sip2.Address `gorm:"-"`
}

func (m *SipChannels) TableName() string {
	return "vid_sip_channels"
}

type SipDevices struct {
	//ID int64 `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	// DeviceID 设备id
	DeviceID  string `gorm:"column:device_id;primary_key;type:char(20);NOT NULL" json:"deviceid"`
	VidmgrID  string `gorm:"column:vidmgr_id;type:char(11);NOT NULL"` // 流服务ID 用于统计设备数量
	MediaIP   int64  `gorm:"column:media_ip" `                        //rtp服务端的ip
	MediaPort int64  `gorm:"column:media_port" `                      //rtp服务端的端口
	// Name 设备名称
	Name string `gorm:"column:name" json:"name"`
	// Region 设备域
	Region string `gorm:"column:region" json:"region"`
	// Host Via 地址
	Host string `gorm:"column:host" json:"host"`
	// Port via 端口
	Port string `gorm:"column:port" json:"port"`
	// TransPort via transport
	TransPort string `gorm:"column:transport" json:"transport"`
	// Proto 协议
	Proto string `gorm:"column:proto" json:"proto"`
	// Rport via rport
	Rport string `gorm:"column:report" json:"report"`
	// RAddr via recevied
	RAddr string `gorm:"column:raddr" json:"raddr"`
	// Manufacturer 制造厂商
	Manufacturer string `gorm:"column:manufacturer" xml:"Manufacturer"  json:"manufacturer"`
	// 设备类型DVR，NVR
	DeviceType string `gorm:"column:devicetype" xml:"DeviceType"  json:"devicetype"`
	// Firmware 固件版本
	Firmware string `gorm:"column:firmware" json:"firmware"`
	// Model 型号
	Model  string `gorm:"column:model" json:"model"`
	URIStr string `gorm:"column:uri" json:"uri"`

	//  最后心跳检测时间
	LastLogin int64 `gorm:"column:last_login" json:"active"`
	// Regist 是否注册
	Regist bool `gorm:"column:regist" json:"regist"`
	// PWD 密码
	PWD string `gorm:"column:pwd" json:"pwd"`

	// Status 状态:  unactivated:0  on:1 在线  offline:2 下在
	Status int `gorm:"column:status;type:smallint;default:0"`
	// Source
	Source string `gorm:"column:source" json:"source"`
	stores.Time
	//---- 数据类型外部使用
	Sys     *GbSipInfo    `json:"sysinfo" gorm:"-"`
	Taddr   *sip2.Address `gorm:"-"`
	Tsource net.Addr      `gorm:"-"`
}

func (m *SipDevices) TableName() string {
	return "vid_sip_devices"
}

// 存放GBSIP通用信息(数据类型外部使用)
type GbSipInfo struct {
	Region  string `json:"region" mapstructure:"region"`          // Region 当前域
	CID     string `json:"cid"    mapstructure:"cid"`             // CID 通道id固定头部
	CNUM    int32  `json:"cnum" bson:"unum"  mapstructure:"unum"` // CNUM 当前通道数
	DID     string `json:"did" bson:"did"    mapstructure:"did"`  // DID 设备id固定头部
	DNUM    int32  `json:"dnum" bson:"dnum"  mapstructure:"dnum"` // DNUM 当前设备数
	LID     string `json:"lid" bson:"lid"    mapstructure:"lid"`  // LID 当前服务id
	SipIp   string `json:"-"`                                     //SIP服务IP
	SipPort int32  `json:"-"`                                     //SIP服务端口
	NetType string `json:"-"`                                     //SIP服务
}
