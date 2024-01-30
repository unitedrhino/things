package media

import (
	"encoding/xml"
	sip2 "github.com/i-Things/things/src/vidsip/gosip/sip"
	db "github.com/i-Things/things/src/vidsip/internal/repo/relationDB"
	"sync"
	"time"
)

const (
	// NotifyMethodUserActive 设备活跃状态通知
	NotifyMethodDevicesActive = "devices.active"
	// NotifyMethodUserRegister 设备注册通知
	NotifyMethodDevicesRegister = "devices.regiester"
	// NotifyMethodDeviceActive 通道活跃通知
	NotifyMethodChannelsActive = "channels.active"
	// NotifyMethodRecordStop 视频录制结束
	NotifyMethodRecordStop = "records.stop"
)

const (
	DeviceStatusON  = "ON"
	DeviceStatusOFF = "OFF"
	defaultLimit    = 20
	defaultSort     = "-addtime"
)

// RecordItem 目录详情
type RecordItem struct {
	// DeviceID 设备编号
	DeviceID string `xml:"DeviceID" bson:"DeviceID" json:"DeviceID"`
	// Name 设备名称
	Name      string `xml:"Name" bson:"Name" json:"Name"`
	FilePath  string `xml:"FilePath" bson:"FilePath" json:"FilePath"`
	Address   string `xml:"Address" bson:"Address" json:"Address"`
	StartTime string `xml:"StartTime" bson:"StartTime" json:"StartTime"`
	EndTime   string `xml:"EndTime" bson:"EndTime" json:"EndTime"`
	Secrecy   int    `xml:"Secrecy" bson:"Secrecy" json:"Secrecy"`
	Type      string `xml:"Type" bson:"Type" json:"Type"`
}

// MessageRecordInfoResponse 目录列表
type MessageRecordInfoResponse struct {
	CmdType  string       `xml:"CmdType"`
	SN       int          `xml:"SN"`
	DeviceID string       `xml:"DeviceID"`
	SumNum   int          `xml:"SumNum"`
	Item     []RecordItem `xml:"RecordList>Item"`
}

// MessageNotify 心跳包xml结构
type MessageNotify struct {
	CmdType  string `xml:"CmdType"`
	SN       int    `xml:"SN"`
	DeviceID string `xml:"DeviceID"`
	Status   string `xml:"Status"`
	Info     string `xml:"Info"`
}

// MessageReceive 接收到的请求数据最外层，主要用来判断数据类型
type MessageReceive struct {
	CmdType string `xml:"CmdType"`
	SN      int    `xml:"SN"`
}

// MessageDeviceInfoResponse 主设备明细返回结构
type MessageDeviceInfoResponse struct {
	CmdType      string `xml:"CmdType"`
	SN           int    `xml:"SN"`
	DeviceID     string `xml:"DeviceID"`
	DeviceType   string `xml:"DeviceType"`
	Manufacturer string `xml:"Manufacturer"`
	Model        string `xml:"Model"`
	Firmware     string `xml:"Firmware"`
}

// MessageDeviceListResponse 设备明细列表返回结构
type MessageDeviceListResponse struct {
	XMLName  xml.Name         `xml:"Response"`
	CmdType  string           `xml:"CmdType"`
	SN       int              `xml:"SN"`
	DeviceID string           `xml:"DeviceID"`
	SumNum   int              `xml:"SumNum"`
	Item     []db.SipChannels `xml:"DeviceList>Item"`
}

//type MessageChannel struct {
//	ChannelID    string `xml:"DeviceID"`
//	DeviceID     string `xml:"-"`
//	Name         string `xml:"Name" `
//	Manufacturer string `xml:"Manufacturer" `
//	Model        string `xml:"Model" `
//	Owner        string `xml:"Owner"`
//	CivilCode    string `xml:"CivilCode"`
//	Address      string `xml:"Address"` // Address ip地址
//	Parental     int32  `xml:"Parental"`
//	SafetyWay    int32  `xml:"SafetyWay"`
//	RegisterWay  int32  `xml:"RegisterWay"`
//	Secrecy      int32  `xml:"Secrecy"`
//	Status       string `xml:"Status"` // Status 状态  on 在线
//}

var deviceStatusMap = map[string]string{
	"ON":     DeviceStatusON,
	"OK":     DeviceStatusON,
	"ONLINE": DeviceStatusON,
	"OFFILE": DeviceStatusOFF,
	"OFF":    DeviceStatusOFF,
}

// Notify 消息通知结构
type Notify struct {
	Method string      `json:"method"`
	Data   interface{} `json:"data"`
}

type RecordList struct {
	channelid string
	resp      chan Records
	num       int
	data      [][]int64
	l         *sync.Mutex
	s, e      int64
}

//type GbSipDevice struct {
//	RandomStr    *utils.SnowFlake
//	DeviceID     string // DeviceID 设备id
//	Name         string // Name 设备名称
//	Region       string // Region 设备域
//	Host         string // Host Via 地址
//	Port         string // Port via 端口
//	TransPort    string // TransPort via transport
//	Proto        string // Proto 协议
//	Rport        string // Rport via rport
//	RAddr        string // RAddr via recevied
//	Manufacturer string // Manufacturer 制造厂商
//	DeviceType   string // 设备类型DVR，NVR
//	Firmware     string // Firmware 固件版本
//	Model        string // Model 型号
//	URIStr       string
//	PWD          string // PWD 密码
//	Source       string
//	addr         *sip.Address
//	source       net.Addr
//	Sys          *GbSipInfo
//}

// ActiveDevices 记录当前活跃设备，请求播放时设备必须处于活跃状态
type ActiveDevices struct {
	sync.Map
}

func (a *ActiveDevices) Get(key string) (db.SipDevices, bool) {
	if v, ok := a.Load(key); ok {
		return v.(db.SipDevices), ok
	}
	return db.SipDevices{}, false
}

type M map[string]interface{}

// 请求是的直播还还是历史
type Stream struct {
	Type      int    //流类型  0直播 1录像回放
	ChannelID string //通道ID
	DeviceID  string //设备ID
	//Stream     string //ssrc
	//ChnURIStr  string //如：sip:33020000081318000006@330200000
	//DevNetType string //设备SIP(udp或者tcp类型)
	//DevSource  string //设备SIP端口:如192.168.10.117:5060
	//MediaPort  int    //流服务的端口  10000
	//MediaIP    string //流服务的IP    192.168.100.125
	//// m3u8播放地址
	//HTTP string `json:"http" gorm:"column:http"`
	//// rtmp 播放地址
	//RTMP string `json:"rtmp" gorm:"column:rtmp"`
	//// rtsp 播放地址
	//RTSP string `json:"rtsp" gorm:"column:rtsp"`
	//// flv 播放地址
	//WSFLV string `json:"wsflv" gorm:"column:wsflv"`
	//pull 媒体服务器主动拉流，push 监控设备主动推流
	StreamType string `json:"streamtype"`
	// ---
	S, E time.Time
	ssrc string // 国标ssrc 10进制字符串
	Ext  int64  // 流等待过期时间
	Resp *sip2.Response

	// header callid
	CallID string `json:"callid" `
	CseqNo uint32 `json:"cseqno" `
	Msg    string `json:"msg"`
	Ttag   M
	Ftag   M
	// 0正常 1关闭 -1 尚未开始
	Status int `json:"status"`
}

///********************************** GB28181 数据 ***********************************/
//type SipChannel struct {
//	ID int64 `gorm:"column:id;primary_key;AUTO_INCREMENT"`
//	// ChannelID 通道编码
//	ChannelID string `gorm:"column:channel_id;type:char(24);NOT NULL"`
//	// DeviceID 设备编号
//	DeviceID string `gorm:"column:device_id;type:char(24);NOT NULL"`
//	// Memo 备注（用来标示通道信息）
//	Memo string `gorm:"column:memo"`
//	// Name 通道名称（设备端设置名称）
//	Name         string `gorm:"column:name"`
//	Manufacturer string `gorm:"column:manufacturer"`
//	Model        string `gorm:"column:model"`
//	Owner        string `gorm:"column:owner"`
//	CivilCode    string `gorm:"column:civilcode"`
//	// Address ip地址
//	Address     string `gorm:"column:address"`
//	Parental    int32  `gorm:"column:parental"`
//	SafetyWay   int32  `gorm:"column:safetyway"`
//	RegisterWay int32  `gorm:"column:registerway"`
//	Secrecy     int32  `gorm:"column:secrecy"`
//	// Status 状态  on 在线
//	Status string `gorm:"column:status"`
//	// Active 最后活跃时间
//	Active int64  `gorm:"column:active"`
//	URIStr string `gorm:"column:uri"`
//	// 视频编码格式
//	VF string `gorm:"column:vf"`
//	// 视频高
//	Height int32 `gorm:"column:height"`
//	// 视频宽
//	Width int32 `gorm:"column:width"`
//	// 视频FPS
//	FPS int32 `gorm:"column:fps"`
//	//  pull 媒体服务器主动拉流，push 监控设备主动推流
//	StreamType string `gorm:"column:streamtype"`
//	// streamtype=pull时，拉流地址
//	URL  string       `gorm:"column:url"`
//	addr *sip.Address `gorm:"-"`
//	stores.Time
//	Owener    string
//	LastLogin time.Time `gorm:"column:last_login"` // 最后登录时间
//}

//func SipDeviceToDB(s *GbSipDevice, db *db.SipDevices) {
//	db.Region = s.Region
//	db.Host = s.Host
//	db.Port = s.Port
//	db.TransPort = s.TransPort
//	db.Proto = s.Proto
//	db.Rport = s.Rport
//	db.RAddr = s.RAddr
//	db.Manufacturer = s.Manufacturer
//	db.DeviceType = s.DeviceType
//	db.Firmware = s.Firmware
//	db.Model = s.Model
//	db.URIStr = s.URIStr
//	db.Source = s.Source
//	//db.Source = s.Source
//	//db.Addr = s.Addr
//}

// Records Records
type Records struct {
	// 存在录像的天数
	DayTotal int          `json:"daynum"`
	TimeNum  int          `json:"timenum"`
	Data     []RecordDate `json:"list"`
}
type RecordDate struct {
	// 日期
	Date string `json:"date"`
	// 时间段
	Items []RecordInfo `json:"items"`
}

// RecordInfo RecordInfo
type RecordInfo struct {
	Start int64 `json:"start" bson:"start"`
	End   int64 `json:"end" bson:"end"`
}
