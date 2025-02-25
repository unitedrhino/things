package protocols

type Type = string //协议类型

const (
	TypeMedia  = "media"  //音视频
	TypeNormal = "normal" //普通设备接入
)

type Trans = string //传输层协议定义

const (

	//音视频协议接入
	ProtocolRtspRtmp = "rtspRtmp" //固定地址,rtsp和rtmp协议接入
	ProtocolGB28181  = "gb28181"  // GB/T 28181 国标视频接入
	ProtocolOnvif    = "onvif"    //

	//普通设备接入
	ProtocolCloud Trans = "cloud" //云云对接
	ProtocolMqtt  Trans = "mqtt"  //device
	ProtocolOther Trans = "other" //其他
	ProtocolHttp  Trans = "http"  //device
	ProtocolTcp   Trans = "tcp"
	ProtocolUdp   Trans = "udp"
)

// 协议编码
const (
	ProtocolCodeUrMqtt       = "urMqtt"
	ProtocolCodeUrMedia      = "urMedia" //联犀音视频协议
	ProtocolCodeUrHttp       = "urHttp"  //联犀音视频协议
	ProtocolCodeAliyunCloud  = "aliyunCloud"
	ProtocolCodeAliyunMqtt   = "aliyunMqtt"
	ProtocolCodeTencentCloud = "tencentCloud"
	ProtocolCodeModbus       = "modbus"
)
