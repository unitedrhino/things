package initServer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/deviceAuth"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/logic"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
)

type ServerHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
	PiDB *relationDB.VidmgrInfoRepo
}

type ServerConfig struct {
	ApiDebug                       string `json:"api.apiDebug,omitempty"`
	ApiDefaultSnap                 string `json:"api.defaultSnap,omitempty"`
	ApiSecret                      string `json:"api.secret,omitempty"`
	ApiSnapRoot                    string `json:"api.snapRoot,omitempty"`
	ClusterOriginUrl               string `json:"cluster.origin_url,omitempty"`
	ClusterRetryCount              string `json:"cluster.retry_count,omitempty"`
	ClusterTimeoutSec              string `json:"cluster.timeout_sec,omitempty"`
	FfmpegBin                      string `json:"ffmpeg.bin,omitempty"`
	FfmpegCmd                      string `json:"ffmpeg.cmd,omitempty"`
	FfmpegLog                      string `json:"ffmpeg.log,omitempty"`
	FfmpegRestartSec               string `json:"ffmpeg.restart_sec,omitempty"`
	FfmpegSnap                     string `json:"ffmpeg.snap,omitempty"`
	GeneralCheckNvidiaDev          string `json:"general.check_nvidia_dev,omitempty"`
	GeneralEnableVhost             string `json:"general.enableVhost,omitempty"`
	GeneralEnableFfmpegLog         string `json:"general.enable_ffmpeg_log,omitempty"`
	GeneralFlowThreshold           string `json:"general.flowThreshold,omitempty"`
	GeneralMaxStreamWaitMS         string `json:"general.maxStreamWaitMS,omitempty"`
	GeneralMediaServerId           string `json:"general.mediaServerId,omitempty"`
	GeneralMergeWriteMS            string `json:"general.mergeWriteMS,omitempty"`
	GeneralResetWhenRePlay         string `json:"general.resetWhenRePlay,omitempty"`
	GeneralStreamNoneReaderDelayMS string `json:"general.streamNoneReaderDelayMS,omitempty"`
	GeneralUnreadyFrameCache       string `json:"general.unready_frame_cache,omitempty"`
	GeneralWaitAddTrackMs          string `json:"general.wait_add_track_ms,omitempty"`
	GeneralWaitTrackReadyMs        string `json:"general.wait_track_ready_ms,omitempty"`
	HlsBroadcastRecordTs           string `json:"hls.broadcastRecordTs,omitempty"`
	HlsDeleteDelaySec              string `json:"hls.deleteDelaySec,omitempty"`
	HlsFileBufSize                 string `json:"hls.fileBufSize,omitempty"`
	HlsSegDur                      string `json:"hls.segDur,omitempty"`
	HlsSegKeep                     string `json:"hls.segKeep,omitempty"`
	HlsSegNum                      string `json:"hls.segNum,omitempty"`
	HlsSegRetain                   string `json:"hls.segRetain,omitempty"`
	HookAliveInterval              string `json:"hook.alive_interval,omitempty"`
	HookEnable                     string `json:"hook.enable,omitempty"`
	HookOnFlowReport               string `json:"hook.on_flow_report,omitempty"`
	HookOnHttpAccess               string `json:"hook.on_http_access,omitempty"`
	HookOnPlay                     string `json:"hook.on_play,omitempty"`
	HookOnPublish                  string `json:"hook.on_publish,omitempty"`
	HookOnRecordMp4                string `json:"hook.on_record_mp4,omitempty"`
	HookOnRecordTs                 string `json:"hook.on_record_ts,omitempty"`
	HookOnRtpServerTimeout         string `json:"hook.on_rtp_server_timeout,omitempty"`
	HookOnRtspAuth                 string `json:"hook.on_rtsp_auth,omitempty"`
	HookOnRtspRealm                string `json:"hook.on_rtsp_realm,omitempty"`
	HookOnSendRtpStopped           string `json:"hook.on_send_rtp_stopped,omitempty"`
	HookOnServerExited             string `json:"hook.on_server_exited,omitempty"`
	HookOnServerKeepalive          string `json:"hook.on_server_keepalive,omitempty"`
	HookOnServerStarted            string `json:"hook.on_server_started,omitempty"`
	HookOnShellLogin               string `json:"hook.on_shell_login,omitempty"`
	HookOnStreamChanged            string `json:"hook.on_stream_changed,omitempty"`
	HookOnStreamNoneReader         string `json:"hook.on_stream_none_reader,omitempty"`
	HookOnStreamNotFound           string `json:"hook.on_stream_not_found,omitempty"`
	HookRetry                      string `json:"hook.retry,omitempty"`
	HookRetryDelay                 string `json:"hook.retry_delay,omitempty"`
	HookStreamChangedSchemas       string `json:"hook.stream_changed_schemas,omitempty"`
	HookTimeoutSec                 string `json:"hook.timeoutSec,omitempty"`
	HttpAllowCrossDomains          string `json:"http.allow_cross_domains,omitempty"`
	HttpAllowIpRange               string `json:"http.allow_ip_range,omitempty"`
	HttpCharSet                    string `json:"http.charSet,omitempty"`
	HttpDirMenu                    string `json:"http.dirMenu,omitempty"`
	HttpForbidCacheSuffix          string `json:"http.forbidCacheSuffix,omitempty"`
	HttpForwardedIpHeader          string `json:"http.forwarded_ip_header,omitempty"`
	HttpKeepAliveSecond            string `json:"http.keepAliveSecond,omitempty"`
	HttpMaxReqSize                 string `json:"http.maxReqSize,omitempty"`
	HttpNotFound                   string `json:"http.notFound,omitempty"`
	HttpPort                       string `json:"http.port,omitempty"`
	HttpRootPath                   string `json:"http.rootPath,omitempty"`
	HttpSendBufSize                string `json:"http.sendBufSize,omitempty"`
	HttpSslport                    string `json:"http.sslport,omitempty"`
	HttpVirtualPath                string `json:"http.virtualPath,omitempty"`
	MulticastAddrMax               string `json:"multicast.addrMax,omitempty"`
	MulticastAddrMin               string `json:"multicast.addrMin,omitempty"`
	MulticastUdpTTL                string `json:"multicast.udpTTL,omitempty"`
	ProtocolAddMuteAudio           string `json:"protocol.add_mute_audio,omitempty"`
	ProtocolAutoClose              string `json:"protocol.auto_close,omitempty"`
	ProtocolContinuePushMs         string `json:"protocol.continue_push_ms,omitempty"`
	ProtocolEnableAudio            string `json:"protocol.enable_audio,omitempty"`
	ProtocolEnableFmp4             string `json:"protocol.enable_fmp4,omitempty"`
	ProtocolEnableHls              string `json:"protocol.enable_hls,omitempty"`
	ProtocolEnableHlsFmp4          string `json:"protocol.enable_hls_fmp4,omitempty"`
	ProtocolEnableMp4              string `json:"protocol.enable_mp4,omitempty"`
	ProtocolEnableRtmp             string `json:"protocol.enable_rtmp,omitempty"`
	ProtocolEnableRtsp             string `json:"protocol.enable_rtsp,omitempty"`
	ProtocolEnableTs               string `json:"protocol.enable_ts,omitempty"`
	ProtocolFmp4Demand             string `json:"protocol.fmp4_demand,omitempty"`
	ProtocolHlsDemand              string `json:"protocol.hls_demand,omitempty"`
	ProtocolHlsSavePath            string `json:"protocol.hls_save_path,omitempty"`
	ProtocolModifyStamp            string `json:"protocol.modify_stamp,omitempty"`
	ProtocolMp4AsPlayer            string `json:"protocol.mp4_as_player,omitempty"`
	ProtocolMp4MaxSecond           string `json:"protocol.mp4_max_second,omitempty"`
	ProtocolMp4SavePath            string `json:"protocol.mp4_save_path,omitempty"`
	ProtocolRtmpDemand             string `json:"protocol.rtmp_demand,omitempty"`
	ProtocolRtspDemand             string `json:"protocol.rtsp_demand,omitempty"`
	ProtocolTsDemand               string `json:"protocol.ts_demand,omitempty"`
	RecordAppName                  string `json:"record.appName,omitempty"`
	RecordFastStart                string `json:"record.fastStart,omitempty"`
	RecordFileBufSize              string `json:"record.fileBufSize,omitempty"`
	RecordFileRepeat               string `json:"record.fileRepeat,omitempty"`
	RecordSampleMS                 string `json:"record.sampleMS,omitempty"`
	RtcExternIP                    string `json:"rtc.externIP,omitempty"`
	RtcPort                        string `json:"rtc.port,omitempty"`
	RtcPreferredCodecA             string `json:"rtc.preferredCodecA,omitempty"`
	RtcPreferredCodecV             string `json:"rtc.preferredCodecV,omitempty"`
	RtcRembBitRate                 string `json:"rtc.rembBitRate,omitempty"`
	RtcTcpPort                     string `json:"rtc.tcpPort,omitempty"`
	RtcTimeoutSec                  string `json:"rtc.timeoutSec,omitempty"`
	RtmpHandshakeSecond            string `json:"rtmp.handshakeSecond,omitempty"`
	RtmpKeepAliveSecond            string `json:"rtmp.keepAliveSecond,omitempty"`
	RtmpPort                       string `json:"rtmp.port,omitempty"`
	RtmpSslport                    string `json:"rtmp.sslport,omitempty"`
	RtpAudioMtuSize                string `json:"rtp.audioMtuSize,omitempty"`
	RtpH264StapA                   string `json:"rtp.h264_stap_a,omitempty"`
	RtpLowLatency                  string `json:"rtp.lowLatency,omitempty"`
	RtpRtpMaxSize                  string `json:"rtp.rtpMaxSize,omitempty"`
	RtpVideoMtuSize                string `json:"rtp.videoMtuSize,omitempty"`
	RtpProxyDumpDir                string `json:"rtp_proxy.dumpDir,omitempty"`
	RtpProxyGopCache               string `json:"rtp_proxy.gop_cache,omitempty"`
	RtpProxyH264Pt                 string `json:"rtp_proxy.h264_pt,omitempty"`
	RtpProxyH265Pt                 string `json:"rtp_proxy.h265_pt,omitempty"`
	RtpProxyOpusPt                 string `json:"rtp_proxy.opus_pt,omitempty"`
	RtpProxyPort                   string `json:"rtp_proxy.port,omitempty"`
	RtpProxyPortRange              string `json:"rtp_proxy.port_range,omitempty"`
	RtpProxyPsPt                   string `json:"rtp_proxy.ps_pt,omitempty"`
	RtpProxyTimeoutSec             string `json:"rtp_proxy.timeoutSec,omitempty"`
	RtspAuthBasic                  string `json:"rtsp.authBasic,omitempty"`
	RtspDirectProxy                string `json:"rtsp.directProxy,omitempty"`
	RtspHandshakeSecond            string `json:"rtsp.handshakeSecond,omitempty"`
	RtspKeepAliveSecond            string `json:"rtsp.keepAliveSecond,omitempty"`
	RtspLowLatency                 string `json:"rtsp.lowLatency,omitempty"`
	RtspPort                       string `json:"rtsp.port,omitempty"`
	RtspRtpTransportType           string `json:"rtsp.rtpTransportType,omitempty"`
	RtspSslport                    string `json:"rtsp.sslport,omitempty"`
	ShellMaxReqSize                string `json:"shell.maxReqSize,omitempty"`
	ShellPort                      string `json:"shell.port,omitempty"`
	SrtLatencyMul                  string `json:"srt.latencyMul,omitempty"`
	SrtPktBufSize                  string `json:"srt.pktBufSize,omitempty"`
	SrtPort                        string `json:"srt.port,omitempty"`
	SrtTimeoutSec                  string `json:"srt.timeoutSec,omitempty"`
}

func NewServerHandle(ctx context.Context, svcCtx *svc.ServiceContext) *ServerHandle {
	return &ServerHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		PiDB:   relationDB.NewVidmgrInfoRepo(ctx),
	}
}

func (l *ServerHandle) ActionCheck() error {
	//l.Infof("ActionCheck req:%v", in)
	fmt.Println("[****] func (l *ServerHandle) ActionCheck() error ")
	var (
		c      = l.svcCtx.Config
		filter = relationDB.VidmgrFilter{VidmgrIpV4: utils.InetAtoN(c.Mediakit.Host), VidmgrPort: c.Mediakit.Port}
	)
	size, err := l.PiDB.CountByFilter(l.ctx, filter)
	if err != nil {
		fmt.Errorf("MediaServer init data countfilter error")
		return err
	}
	if size > 0 {
		//update
		page := vid.PageInfo{}
		di, err := l.PiDB.FindByFilter(l.ctx, filter, logic.ToPageInfoWithDefault(&page, &def.PageInfo{
			Page: 1, Size: 20,
			Orders: []def.OrderBy{{"created_time", def.OrderDesc}, {"vidmgr_id", def.OrderDesc}},
		}))
		if err != nil {
			fmt.Errorf("MediaServer init data find filter error")
			return err
		}
		if di[0].VidmgrSecret != c.Mediakit.Secret {
			di[0].VidmgrSecret = c.Mediakit.Secret
			err = l.PiDB.Update(l.ctx, di[0])
		}
	} else {
		//create
		dbDocker := &relationDB.VidmgrInfo{
			VidmgrID:     deviceAuth.GetStrProductID(l.svcCtx.VidmgrID.GetSnowflakeId()),
			VidmgrName:   "default Docker",
			VidmgrIpV4:   utils.InetAtoN(c.Mediakit.Host),
			VidmgrPort:   c.Mediakit.Port,
			VidmgrSecret: c.Mediakit.Secret,
			VidmgrStatus: 2, //默认设置离线状态
			VidmgrType:   1, //ZLmediakit
			MediasvrType: 1, //docker模式
			Desc:         "",
			Tags:         map[string]string{},
		}
		err = l.PiDB.Insert(l.ctx, dbDocker)
		if err != nil {
			l.Errorf("%s.Insert err=%+v", utils.FuncName(), err)
			return err
		}
	}
	//config dockerServer
	config := new(ServerConfig)
	SetDefaultConfig(l.svcCtx, config)
	byte4, err := json.Marshal(config)
	var tdata map[string]interface{}
	err = json.Unmarshal(byte4, &tdata)
	tdata["secret"] = c.Mediakit.Secret
	byte4, err = json.Marshal(tdata)
	if err != nil {
		er := errors.Fmt(err)
		fmt.Print("%s map string phares failed  err=%+v", utils.FuncName(), er)
		return er
	}
	preUrl := fmt.Sprintf("http://%s:%s/index/api/setServerConfig", c.Mediakit.Host, c.Mediakit.Port)
	request, error := http.NewRequest("POST", preUrl, bytes.NewBuffer(byte4))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		fmt.Println(err)
		return err
	}
	defer response.Body.Close()
	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	//body, _ := ioutil.ReadAll(response.Body)
	body, err := io.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))

	return nil
}

func SetDefaultConfig(svcCtx *svc.ServiceContext, config *ServerConfig) {
	HostIP := svcCtx.Config.Mediakit.Host
	HPort := svcCtx.Config.Mediakit.Port
	//HPort
	config.ApiDebug = "1"
	config.HookEnable = "1"
	config.HookOnFlowReport = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onFlowReport", HostIP, HPort)
	config.HookOnHttpAccess = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onHttpAccess", HostIP, HPort)
	config.HookOnPlay = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onPlay", HostIP, HPort)
	config.HookOnPublish = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onPublish", HostIP, HPort)
	config.HookOnRecordMp4 = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onRecordMp4", HostIP, HPort)
	config.HookOnRecordTs = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onRecordTs", HostIP, HPort)
	config.HookOnRtpServerTimeout = fmt.Sprintf("https://%s:%d/api/v1/things/vidmgr/hooks/onRtpServerTimeout", HostIP, HPort)
	config.HookOnRtspAuth = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onRtspAuth", HostIP, HPort)
	config.HookOnRtspRealm = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onRtspRealm", HostIP, HPort)
	config.HookOnSendRtpStopped = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onSendRtpStopped", HostIP, HPort)
	config.HookOnServerExited = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onServerExited", HostIP, HPort)
	config.HookOnServerKeepalive = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onServerKeepalive", HostIP, HPort)
	config.HookOnServerStarted = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onServerStarted", HostIP, HPort)
	config.HookOnShellLogin = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onShellLogin", HostIP, HPort)
	config.HookOnStreamChanged = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onStreamChanged", HostIP, HPort)
	config.HookOnStreamNoneReader = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onStreamNoneReader", HostIP, HPort)
	config.HookOnStreamNotFound = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onStreamNotFound", HostIP, HPort)
}
