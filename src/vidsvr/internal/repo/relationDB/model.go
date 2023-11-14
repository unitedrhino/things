package relationDB

import (
	"database/sql"
	"github.com/i-Things/things/shared/stores"
)

type DmExample struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
}

// 流服务信息
type VidmgrInfo struct {
	VidmgrID     string            `gorm:"column:id;type:char(11);primary_key;NOT NULL"`                // 服务id
	VidmgrName   string            `gorm:"column:name;type:varchar(100);NOT NULL"`                      // 服务名称
	VidmgrIpV4   int64             `gorm:"column:ipv4;type:bigint"`                                     // 服务IP
	VidmgrPort   int64             `gorm:"column:port;type:bigint"`                                     // 服务端口
	VidmgrType   int64             `gorm:"column:type;type:smallint;default:1"`                         // 服务类型:1:zlmediakit,2:srs,3:monibuca
	VidmgrStatus int64             `gorm:"column:status;type:smallint;default:0;NOT NULL"`              //服务状态: 0：未激活 1：在线  2:离线
	VidmgrSecret string            `gorm:"column:secret;type:varchar(50)"`                              // 服务秘钥
	FirstLogin   sql.NullTime      `gorm:"column:first_login"`                                          // 激活后首次登录时间
	LastLogin    sql.NullTime      `gorm:"column:last_login"`                                           // 最后登录时间
	Desc         string            `gorm:"column:desc;type:varchar(200)"`                               // 描述
	Tags         map[string]string `gorm:"column:tags;type:json;serializer:json;NOT NULL;default:'{}'"` // 产品标签
	stores.Time
}

func (m *VidmgrInfo) TableName() string {
	return "vid_mgr_info"
}

// 流服务配置表
type VidmgrConfig struct {
	GeneralMediaServerId           string `gorm:"column:general_mediaServerId;primary_key;NOT NULL"`
	ApiApiDebug                    string `gorm:"column:api_apiDebug"`
	ApiDefaultSnap                 string `gorm:"column:api_defaultSnap"`
	ApiSecret                      string `gorm:"column:api_secret"`
	ApiSnapRoot                    string `gorm:"column:api_snapRoot"`
	ClusterOriginUrl               string `gorm:"column:cluster_origin_url"`
	ClusterRetryCount              string `gorm:"column:cluster_retry_count"`
	ClusterTimeoutSec              string `gorm:"column:cluster_timeout_sec"`
	FfmpegBin                      string `gorm:"column:ffmpeg_bin"`
	FfmpegCmd                      string `gorm:"column:ffmpeg_cmd"`
	FfmpegLog                      string `gorm:"column:ffmpeg_log"`
	FfmpegRestartSec               string `gorm:"column:ffmpeg_restart_sec"`
	FfmpegSnap                     string `gorm:"column:ffmpeg_snap"`
	GeneralCheckNvidiaDev          string `gorm:"column:general_check_nvidia_dev"`
	GeneralEnableVhost             string `gorm:"column:general_enableVhost"`
	GeneralEnableFfmpegLog         string `gorm:"column:general_enable_ffmpeg_log"`
	GeneralFlowThreshold           string `gorm:"column:general_flowThreshold"`
	GeneralMaxStreamWaitMS         string `gorm:"column:general_maxStreamWaitMS"`
	GeneralMergeWriteMS            string `gorm:"column:general_mergeWriteMS"`
	GeneralResetWhenRePlay         string `gorm:"column:general_resetWhenRePlay"`
	GeneralStreamNoneReaderDelayMS string `gorm:"column:general_streamNoneReaderDelayMS"`
	GeneralUnreadyFrameCache       string `gorm:"column:general_unready_frame_cache"`
	GeneralWaitAddTrackMs          string `gorm:"column:general_wait_add_track_ms"`
	GeneralWaitTrackReadyMs        string `gorm:"column:general_wait_track_ready_ms"`
	HlsBroadcastRecordTs           string `gorm:"column:hls_broadcastRecordTs"`
	HlsDeleteDelaySec              string `gorm:"column:hls_deleteDelaySec"`
	HlsFileBufSize                 string `gorm:"column:hls_fileBufSize"`
	HlsSegDur                      string `gorm:"column:hls_segDur"`
	HlsSegKeep                     string `gorm:"column:hls_segKeep"`
	HlsSegNum                      string `gorm:"column:hls_segNum"`
	HlsSegRetain                   string `gorm:"column:hls_segRetain"`
	HookAliveInterval              string `gorm:"column:hook_alive_interval"`
	HookEnable                     string `gorm:"column:hook_enable"`
	HookOnFlowReport               string `gorm:"column:hook_on_flow_report"`
	HookOnHttpAccess               string `gorm:"column:hook_on_http_access"`
	HookOnPlay                     string `gorm:"column:hook_on_play"`
	HookOnPublish                  string `gorm:"column:hook_on_publish"`
	HookOnRecordMp4                string `gorm:"column:hook_on_record_mp4"`
	HookOnRecordTs                 string `gorm:"column:hook_on_record_ts"`
	HookOnRtpServerTimeout         string `gorm:"column:hook_on_rtp_server_timeout"`
	HookOnRtspAuth                 string `gorm:"column:hook_on_rtsp_auth"`
	HookOnRtspRealm                string `gorm:"column:hook_on_rtsp_realm"`
	HookOnSendRtpStopped           string `gorm:"column:hook_on_send_rtp_stopped"`
	HookOnServerExited             string `gorm:"column:hook_on_server_exited"`
	HookOnServerKeepalive          string `gorm:"column:hook_on_server_keepalive"`
	HookOnServerStarted            string `gorm:"column:hook_on_server_started"`
	HookOnShellLogin               string `gorm:"column:hook_on_shell_login"`
	HookOnStreamChanged            string `gorm:"column:hook_on_stream_changed"`
	HookOnStreamNoneReader         string `gorm:"column:hook_on_stream_none_reader"`
	HookOnStreamNotFound           string `gorm:"column:hook_on_stream_not_found"`
	HookRetry                      string `gorm:"column:hook_retry"`
	HookRetryDelay                 string `gorm:"column:hook_retry_delay"`
	HookStreamChangedSchemas       string `gorm:"column:hook_stream_changed_schemas"`
	HookTimeoutSec                 string `gorm:"column:hook_timeoutSec"`
	HttpAllowCrossDomains          string `gorm:"column:http_allow_cross_domains"`
	HttpAllowIpRange               string `gorm:"column:http_allow_ip_range"`
	HttpCharSet                    string `gorm:"column:http_charSet"`
	HttpDirMenu                    string `gorm:"column:http_dirMenu"`
	HttpForbidCacheSuffix          string `gorm:"column:http_forbidCacheSuffix"`
	HttpForwardedIpHeader          string `gorm:"column:http_forwarded_ip_header"`
	HttpKeepAliveSecond            string `gorm:"column:http_keepAliveSecond"`
	HttpMaxReqSize                 string `gorm:"column:http_maxReqSize"`
	HttpNotFound                   string `gorm:"column:http_notFound"`
	HttpPort                       string `gorm:"column:http_port"`
	HttpRootPath                   string `gorm:"column:http_rootPath"`
	HttpSendBufSize                string `gorm:"column:http_sendBufSize"`
	HttpSslport                    string `gorm:"column:http_sslport"`
	HttpVirtualPath                string `gorm:"column:http_virtualPath"`
	MulticastAddrMax               string `gorm:"column:multicast_addrMax"`
	MulticastAddrMin               string `gorm:"column:multicast_addrMin"`
	MulticastUdpTTL                string `gorm:"column:multicast_udpTTL"`
	ProtocolAddMuteAudio           string `gorm:"column:protocol_add_mute_audio"`
	ProtocolAutoClose              string `gorm:"column:protocol_auto_close"`
	ProtocolContinuePushMs         string `gorm:"column:protocol_continue_push_ms"`
	ProtocolEnableAudio            string `gorm:"column:protocol_enable_audio"`
	ProtocolEnableFmp4             string `gorm:"column:protocol_enable_fmp4"`
	ProtocolEnableHls              string `gorm:"column:protocol_enable_hls"`
	ProtocolEnableHlsFmp4          string `gorm:"column:protocol_enable_hls_fmp4"`
	ProtocolEnableMp4              string `gorm:"column:protocol_enable_mp4"`
	ProtocolEnableRtmp             string `gorm:"column:protocol_enable_rtmp"`
	ProtocolEnableRtsp             string `gorm:"column:protocol_enable_rtsp"`
	ProtocolEnableTs               string `gorm:"column:protocol_enable_ts"`
	ProtocolFmp4Demand             string `gorm:"column:protocol_fmp4_demand"`
	ProtocolHlsDemand              string `gorm:"column:protocol_hls_demand"`
	ProtocolHlsSavePath            string `gorm:"column:protocol_hls_save_path"`
	ProtocolModifyStamp            string `gorm:"column:protocol_modify_stamp"`
	ProtocolMp4AsPlayer            string `gorm:"column:protocol_mp4_as_player"`
	ProtocolMp4MaxSecond           string `gorm:"column:protocol_mp4_max_second"`
	ProtocolMp4SavePath            string `gorm:"column:protocol_mp4_save_path"`
	ProtocolRtmpDemand             string `gorm:"column:protocol_rtmp_demand"`
	ProtocolRtspDemand             string `gorm:"column:protocol_rtsp_demand"`
	ProtocolTsDemand               string `gorm:"column:protocol_ts_demand"`
	RecordAppName                  string `gorm:"column:record_appName"`
	RecordFastStart                string `gorm:"column:record_fastStart"`
	RecordFileBufSize              string `gorm:"column:record_fileBufSize"`
	RecordFileRepeat               string `gorm:"column:record_fileRepeat"`
	RecordSampleMS                 string `gorm:"column:record_sampleMS"`
	RtcExternIP                    string `gorm:"column:rtc_externIP"`
	RtcPort                        string `gorm:"column:rtc_port"`
	RtcPreferredCodecA             string `gorm:"column:rtc_preferredCodecA"`
	RtcPreferredCodecV             string `gorm:"column:rtc_preferredCodecV"`
	RtcRembBitRate                 string `gorm:"column:rtc_rembBitRate"`
	RtcTcpPort                     string `gorm:"column:rtc_tcpPort"`
	RtcTimeoutSec                  string `gorm:"column:rtc_timeoutSec"`
	RtmpHandshakeSecond            string `gorm:"column:rtmp_handshakeSecond"`
	RtmpKeepAliveSecond            string `gorm:"column:rtmp_keepAliveSecond"`
	RtmpPort                       string `gorm:"column:rtmp_port"`
	RtmpSslport                    string `gorm:"column:rtmp_sslport"`
	RtpAudioMtuSize                string `gorm:"column:rtp_audioMtuSize"`
	RtpH264StapA                   string `gorm:"column:rtp_h264_stap_a"`
	RtpLowLatency                  string `gorm:"column:rtp_lowLatency"`
	RtpRtpMaxSize                  string `gorm:"column:rtp_rtpMaxSize"`
	RtpVideoMtuSize                string `gorm:"column:rtp_videoMtuSize"`
	RtpProxyDumpDir                string `gorm:"column:rtp_proxy_dumpDir"`
	RtpProxyGopCache               string `gorm:"column:rtp_proxy_gop_cache"`
	RtpProxyH264Pt                 string `gorm:"column:rtp_proxy_h264_pt"`
	RtpProxyH265Pt                 string `gorm:"column:rtp_proxy_h265_pt"`
	RtpProxyOpusPt                 string `gorm:"column:rtp_proxy_opus_pt"`
	RtpProxyPort                   string `gorm:"column:rtp_proxy_port"`
	RtpProxyPortRange              string `gorm:"column:rtp_proxy_port_range"`
	RtpProxyPsPt                   string `gorm:"column:rtp_proxy_ps_pt"`
	RtpProxyTimeoutSec             string `gorm:"column:rtp_proxy_timeoutSec"`
	RtspAuthBasic                  string `gorm:"column:rtsp_authBasic"`
	RtspDirectProxy                string `gorm:"column:rtsp_directProxy"`
	RtspHandshakeSecond            string `gorm:"column:rtsp_handshakeSecond"`
	RtspKeepAliveSecond            string `gorm:"column:rtsp_keepAliveSecond"`
	RtspLowLatency                 string `gorm:"column:rtsp_lowLatency"`
	RtspPort                       string `gorm:"column:rtsp_port"`
	RtspRtpTransportType           string `gorm:"column:rtsp_rtpTransportType"`
	RtspSslport                    string `gorm:"column:rtsp_sslport"`
	ShellMaxReqSize                string `gorm:"column:shell_maxReqSize"`
	ShellPort                      string `gorm:"column:shell_port"`
	SrtLatencyMul                  string `gorm:"column:srt_latencyMul"`
	SrtPktBufSize                  string `gorm:"column:srt_pktBufSize"`
	SrtPort                        string `gorm:"column:srt_port"`
	SrtTimeoutSec                  string `gorm:"column:srt_timeoutSec"`
}

// 流服务激活之后创建该表
func (m *VidmgrConfig) TableName() string {
	return "vid_mgr_config"
}

// 视频流信息
type VidstreamInfo struct {
	StreamID       string            `gorm:"column:id;type:bigint;primary_key;NOT NULL"` // 视频流的id
	StreamName     string            `gorm:"column:name;type:varchar(100);NOT NULL"`     // 服务名称
	NetType        int64             `gorm:"column:net_type;type:smallint;NOT NULL"`     // 服务名称
	DevType        int64             `gorm:"column:dev_type;type:smallint;default:1"`
	DevStreamType  int64             `gorm:"column:dev_streamtype;type:smallint;default:1"`
	ChannelID      string            `gorm:"column:channel_id;type:varchar(32)"`
	ChannelName    string            `gorm:"column:channel_name;type:varchar(32)"`
	LowNetType     int64             `gorm:"column:low_nettype;type:smallint;default:1"`
	IsShareChannel bool              `gorm:"column:share_channel;type:bit(1);default:0;NOT NULL"`
	IsAutoPush     bool              `gorm:"column:auto_push;type:bit(1);default:0;NOT NULL"`
	IsAutoRecord   bool              `gorm:"column:auto_record;type:bit(1);default:0;NOT NULL"`
	IsPTZ          bool              `gorm:"column:is_ptz;type:bit(1);default:0;NOT NULL"`
	IsOnline       bool              `gorm:"column:is_online;type:bit(1);default:0;NOT NULL"`
	VidmgrInfo     *VidmgrInfo       `gorm:"foreignKey:VidmgrID;references:VidmgrID"`                     // 添加外键
	Desc           string            `gorm:"column:desc;type:varchar(200)"`                               // 描述
	Tags           map[string]string `gorm:"column:tags;type:json;serializer:json;NOT NULL;default:'{}'"` // 产品标签
	stores.Time
}

func (m *VidstreamInfo) TableName() string {
	return "vid_stream_info"
}
