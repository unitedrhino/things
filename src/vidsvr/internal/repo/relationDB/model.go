package relationDB

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/i-Things/things/shared/stores"
	"time"
)

// 流服务信息
type VidmgrInfo struct {
	VidmgrID     string    `gorm:"column:vidmgr_id;type:char(11);primary_key;NOT NULL"` // 服务id
	VidmgrName   string    `gorm:"column:name;type:varchar(100);NOT NULL"`              // 服务名称
	VidmgrIpV4   int64     `gorm:"column:ipv4;type:bigint"`                             // 服务IP
	VidmgrPort   int64     `gorm:"column:port;type:bigint"`                             // 服务端口
	VidmgrType   int64     `gorm:"column:type;type:smallint;default:1"`                 // 服务类型:1:zlmediakit,2:srs,3:monibuca
	VidmgrStatus int64     `gorm:"column:status;type:smallint;default:0;NOT NULL"`      //服务状态: 0：未激活 1：在线  2:离线
	VidmgrSecret string    `gorm:"column:secret;type:varchar(50)"`                      // 服务秘钥
	FirstLogin   time.Time `gorm:"column:first_login"`                                  // 激活后首次登录时间
	LastLogin    time.Time `gorm:"column:last_login"`                                   // 最后登录时间
	IsOpenGbSip  bool      `gorm:"column:open_gbsip;type:smallint;default:1"`           // 国标服务是否开启
	RtpPort      int64     `gorm:"column:rtpport;type:bigint"`                          // 国标服务RTP端口(10000)
	MediasvrType int64     `gorm:"column:mediasvr_type;type:smallint;default:2"`        // 流服务部署类型:1,docker部署  2,独立主机
	//使用vid.yaml配置代替
	Desc string            `gorm:"column:desc;type:varchar(200)"`                               // 描述
	Tags map[string]string `gorm:"column:tags;type:json;serializer:json;NOT NULL;default:'{}'"` // 产品标签
	stores.Time
}

func (m *VidmgrInfo) TableName() string {
	return "vid_mgr_info"
}

/********************************** GB28181 数据 ***********************************/
type StreamTrack struct {
	Channels    int64   `json:"channels"`
	CodecId     int64   `json:"codec_id"`
	CodecIdName string  `json:"codec_id_name"`
	CodecType   int64   `json:"codec_type"`
	Ready       bool    `json:"ready"`
	Loss        float64 `json:"loss"`
	SampleBit   int64   `json:"sample_bit"`
	SampleRate  int64   `json:"sample_rate"`
	Fps         int64   `json:"fps"`
	Height      int64   `json:"height"`
	Width       int64   `json:"width"`
}

// type BList []*BStruct
type STracks []*StreamTrack

func (b STracks) Value() (driver.Value, error) {
	d, err := json.Marshal(b)
	return string(d), err
}

// 注意，这里的接收器是指针类型，否则无法把数据从数据库读到结构体
func (b *STracks) Scan(v interface{}) error {
	return json.Unmarshal(v.([]byte), b)
}

/*
 *Protocol为视频协议
 *当前协议支持类型有 rtmp/rtsp/ts/fmp4/hls/hls.fmp4/
 *分别用一个bit位来表示一个协议
 *对应关系:
 *          bit位        	5	          4   	  3   	  2      1      0
 *                          hls.fmp4      hls     fmp4    ts     rtsp   rtmp
 */
// 视频流信息
type VidmgrStream struct {
	StreamID   int64  `gorm:"column:stream_id;type:bigint;primary_key;AUTO_INCREMENT"` // 视频流的id(主键唯一)
	VidmgrID   string `gorm:"column:vidmgr_id;type:char(11);NOT NULL"`                 // 流服务ID  外键
	StreamName string `gorm:"column:name;type:varchar(63)"`                            // 视频流名称

	App    string `gorm:"column:app;type:varchar(31);NOT NULL"`
	Stream string `gorm:"column:stream;type:varchar(31);NOT NULL"`
	Vhost  string `gorm:"column:vhost;type:varchar(31);NOT NULL"`

	Identifier string `gorm:"column:identifier;type:varchar(31)"`
	LocalIP    int64  `gorm:"column:local_ip;type:bigint"`
	LocalPort  int64  `gorm:"column:local_port;type:bigint"`
	PeerIP     int64  `gorm:"column:peer_ip;type:bigint"`
	PeerPort   int64  `gorm:"column:peer_port;type:bigint"`
	//产生源类型，包括 unknown = 0,rtmp_push=1,rtsp_push=2,rtp_push=3,pull=4,ffmpeg_pull=5,mp4_vod=6,device_chn=7,rtc_push=8
	OriginType int64  `gorm:"column:origin_type;type:smallint"` // 源类型
	PullKey    string `gorm:"column:key;type:varchar(50)"`      //仅PULL当为拉流代理模式时会有Key 其它时间为空
	RtpType    int32  `gorm:"column:rtp_type;type:smallint"`    //仅PULL当为拉流代理模式时会有RtpType 其它时间为空

	OriginStr        string `gorm:"column:origin_str;type:char(15)"`
	OriginUrl        string `gorm:"column:origin_url;type:char(63)"`         //产生源的url
	ReaderCount      int64  `gorm:"column:reader_count;type:smallint"`       // 本协议观看人数
	TotalReaderCount int64  `gorm:"column:total_reader_count;type:smallint"` //观看总人数，包括hls/rtsp/rtmp/http-flv/ws-flv/rtc
	//流通道信息
	Tracks STracks `json:"tracks" gorm:"type:json;column:tracks"`
	//
	IsRecordingMp4 bool `gorm:"column:is_recording_mp4;type:bool;default:0;NOT NULL"`
	IsRecordingHLS bool `gorm:"column:is_recording_hls;type:bool;default:0;NOT NULL"`
	IsShareChannel bool `gorm:"column:share_channel;type:bool;default:0;NOT NULL"`
	IsAutoPush     bool `gorm:"column:auto_push;type:bool;default:0;NOT NULL"`
	IsAutoRecord   bool `gorm:"column:auto_record;type:bool;default:0;NOT NULL"`
	IsPTZ          bool `gorm:"column:is_ptz;type:bool;default:0;NOT NULL"`
	//正常流程有注册和注销过程，注册后，该流进行更新；并上线，注销后就设置标志位进行下线。
	//还需要有一个定时器用来检测异常断开的情况超时时间10S
	IsOnline bool `gorm:"column:is_online;type:bool;default:0;NOT NULL"`

	/*Protocol 为可支持的协议类型*/
	//Protocol uint32 `gorm:"column:protocol;type:uint;default:0;NOT NULL"`
	//当前协议支持类型有 rtmp/rtsp/ts/fmp4/hls/hls.fmp4/
	OnRtmp    bool `gorm:"column:on_rtmp;type:bool;default:0;NOT NULL"`
	OnRtsp    bool `gorm:"column:on_rtsp;type:bool;default:0;NOT NULL"`
	OnTs      bool `gorm:"column:on_ts;type:bool;default:0;NOT NULL"`
	OnHls     bool `gorm:"column:on_hls;type:bool;default:0;NOT NULL"`
	OnFmp4    bool `gorm:"column:on_fmp4;type:bool;default:0;NOT NULL"`
	OnHlsFmp4 bool `gorm:"column:on_hls_fmp4;type:bool;default:0;NOT NULL"`

	FirstLogin time.Time         `gorm:"column:first_login"`                                          // 最早登录时间
	LastLogin  time.Time         `gorm:"column:last_login"`                                           // 最后登录时间
	Desc       string            `gorm:"column:desc;type:varchar(200)"`                               // 描述
	Tags       map[string]string `gorm:"column:tags;type:json;serializer:json;NOT NULL;default:'{}'"` // 产品标签
	stores.Time
	VidmgrInfo *VidmgrInfo `gorm:"foreignKey:VidmgrID;references:VidmgrID"` // 添加外键
}

func (m *VidmgrStream) TableName() string {
	return "vid_mgr_stream"
}

// 流服务配置表
type VidmgrConfig struct {
	//ConfigID                       int64  `gorm:"column:config_id;type:bigint;primary_key;AUTO_INCREMENT"` // 视频流的id(主键唯一)
	VidmgrID                       string `gorm:"column:vidmgr_id;type:char(11);primary_key;NOT NULL"` //generalMediaserverID
	ApiDebug                       string `gorm:"column:apiDebug;char(1)"`
	ApiDefaultSnap                 string `gorm:"column:defaultSnap"`
	ApiSecret                      string `gorm:"column:secret"`
	ApiSnapRoot                    string `gorm:"column:snapRoot"`
	ClusterOriginUrl               string `gorm:"column:cluster_origin_url"`
	ClusterRetryCount              string `gorm:"column:cluster_retry_count"`
	ClusterTimeoutSec              string `gorm:"column:cluster_timeout_sec"`
	FfmpegBin                      string `gorm:"column:ffmpeg_bin"`
	FfmpegCmd                      string `gorm:"column:ffmpeg_cmd"`
	FfmpegLog                      string `gorm:"column:ffmpeg_log"`
	FfmpegRestartSec               string `gorm:"column:ffmpeg_restart_sec"`
	FfmpegSnap                     string `gorm:"column:ffmpeg_snap"`
	GeneralCheckNvidiaDev          string `gorm:"column:check_nvidia_dev"`
	GeneralEnableVhost             string `gorm:"column:enableVhost"`
	GeneralEnableFfmpegLog         string `gorm:"column:enable_ffmpeg_log"`
	GeneralFlowThreshold           string `gorm:"column:flowThreshold"`
	GeneralMaxStreamWaitMS         string `gorm:"column:maxStreamWaitMS"`
	GeneralMergeWriteMS            string `gorm:"column:mergeWriteMS"`
	GeneralResetWhenRePlay         string `gorm:"column:resetWhenRePlay"`
	GeneralStreamNoneReaderDelayMS string `gorm:"column:streamNoneReaderDelayMS"`
	GeneralUnreadyFrameCache       string `gorm:"column:unready_frame_cache"`
	GeneralWaitAddTrackMs          string `gorm:"column:wait_add_track_ms"`
	GeneralWaitTrackReadyMs        string `gorm:"column:wait_track_ready_ms"`
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
	stores.Time
	//VidmgrInfo *VidmgrInfo `gorm:"foreignKey:VidmgrID;references:VidmgrID"` // 添加外键
}

// 流服务激活之后创建该表
func (m *VidmgrConfig) TableName() string {
	return "vid_mgr_config"
}
