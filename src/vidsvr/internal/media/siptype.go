package media

import (
	"encoding/xml"
	db "github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"sync"
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
	XMLName  xml.Name            `xml:"Response"`
	CmdType  string              `xml:"CmdType"`
	SN       int                 `xml:"SN"`
	DeviceID string              `xml:"DeviceID"`
	SumNum   int                 `xml:"SumNum"`
	Item     []db.VidmgrChannels `xml:"DeviceList>Item"`
}

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

// ActiveDevices 记录当前活跃设备，请求播放时设备必须处于活跃状态
type ActiveDevices struct {
	sync.Map
}

var (
	_activeDevices ActiveDevices
	_serverDevices *db.VidmgrDevices
	SipInfo        *db.VidmgrSipInfo
	_recordList    *sync.Map
)

type recordList struct {
	channelid string
	resp      chan Records
	num       int
	data      [][]int64
	l         *sync.Mutex
	s, e      int64
}

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

// Get Get
func (a *ActiveDevices) Get(key string) (db.VidmgrDevices, bool) {
	if v, ok := a.Load(key); ok {
		return v.(db.VidmgrDevices), ok
	}
	return db.VidmgrDevices{}, false
}
