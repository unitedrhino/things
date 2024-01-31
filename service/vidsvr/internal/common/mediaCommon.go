package common

import (
	"fmt"
	"github.com/i-Things/things/service/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/vidsvr/internal/types"
)

// 设置流服务为默认配置
func SetDefaultConfig(hostip string, hostport int64, config *types.ServerConfig) {
	//hostport
	config.ApiDebug = "1"
	config.HookEnable = "1"
	config.HookOnFlowReport = fmt.Sprintf("http://%s:%d/api/v1/zlmedia/hooks/onFlowReport", hostip, hostport)
	config.HookOnHttpAccess = fmt.Sprintf("http://%s:%d/api/v1/zlmedia/hooks/onHttpAccess", hostip, hostport)
	config.HookOnPlay = fmt.Sprintf("http://%s:%d/api/v1/zlmedia/hooks/onPlay", hostip, hostport)
	config.HookOnPublish = fmt.Sprintf("http://%s:%d/api/v1/zlmedia/hooks/onPublish", hostip, hostport)
	config.HookOnRecordMp4 = fmt.Sprintf("http://%s:%d/api/v1/zlmedia/hooks/onRecordMp4", hostip, hostport)
	config.HookOnRecordTs = fmt.Sprintf("http://%s:%d/api/v1/zlmedia/hooks/onRecordTs", hostip, hostport)
	config.HookOnRtpServerTimeout = fmt.Sprintf("https://%s:%d/api/v1/zlmedia/hooks/onRtpServerTimeout", hostip, hostport)
	config.HookOnRtspAuth = fmt.Sprintf("http://%s:%d/api/v1/zlmedia/hooks/onRtspAuth", hostip, hostport)
	config.HookOnRtspRealm = fmt.Sprintf("http://%s:%d/api/v1/zlmedia/hooks/onRtspRealm", hostip, hostport)
	config.HookOnSendRtpStopped = fmt.Sprintf("http://%s:%d/api/v1/zlmedia/hooks/onSendRtpStopped", hostip, hostport)
	config.HookOnServerExited = fmt.Sprintf("http://%s:%d/api/v1/zlmedia/hooks/onServerExited", hostip, hostport)
	config.HookOnServerKeepalive = fmt.Sprintf("http://%s:%d/api/v1/zlmedia/hooks/onServerKeepalive", hostip, hostport)
	config.HookOnServerStarted = fmt.Sprintf("http://%s:%d/api/v1/zlmedia/hooks/onServerStarted", hostip, hostport)
	config.HookOnShellLogin = fmt.Sprintf("http://%s:%d/api/v1/zlmedia/hooks/onShellLogin", hostip, hostport)
	config.HookOnStreamChanged = fmt.Sprintf("http://%s:%d/api/v1/zlmedia/hooks/onStreamChanged", hostip, hostport)
	config.HookOnStreamNoneReader = fmt.Sprintf("http://%s:%d/api/v1/zlmedia/hooks/onStreamNoneReader", hostip, hostport)
	config.HookOnStreamNotFound = fmt.Sprintf("http://%s:%d/api/v1/zlmedia/hooks/onStreamNotFound", hostip, hostport)
}

func SetProtocol(schema string, streamInfo *relationDB.VidmgrStream) {
	switch schema {
	case "rtmp":
		streamInfo.OnRtmp = true
	case "rtsp":
		streamInfo.OnRtsp = true
	case "ts":
		streamInfo.OnTs = true
	case "fmp4":
		streamInfo.OnFmp4 = true
	case "hls":
		streamInfo.OnHls = true
	case "hls.fmp4":
		streamInfo.OnHlsFmp4 = true
	default:
	}
}

func UnSetProtocol(schema string, streamInfo *relationDB.VidmgrStream) {
	switch schema {
	case "rtmp":
		streamInfo.OnRtmp = false
	case "rtsp":
		streamInfo.OnRtsp = false
	case "ts":
		streamInfo.OnTs = false
	case "fmp4":
		streamInfo.OnFmp4 = false
	case "hls":
		streamInfo.OnHls = false
	case "hls.fmp4":
		streamInfo.OnHlsFmp4 = false
	default:
	}
}

//func CheckProtocol(streamInfo *relationDB.VidmgrStream) bool {
//	return streamInfo.OnRtmp || streamInfo.OnRtsp || streamInfo.OnTs || streamInfo.OnFmp4 || streamInfo.OnHls || streamInfo.OnHlsFmp4
//}
