package vidmgrinfomanagelogic

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/i-Things/things/src/vidsvr/pb/vid"
	"time"
)

/*
根据用户的输入生成对应的数据库数据
*/
func ConvVidmgrPbToPo(in *vid.VidmgrInfo) (*relationDB.VidmgrInfo, error) {
	pi := &relationDB.VidmgrInfo{
		VidmgrID:     in.VidmgrID,
		VidmgrName:   in.VidmgrName,
		VidmgrIpV4:   utils.InetAtoN(in.VidmgrIpV4),
		VidmgrPort:   in.VidmgrPort,
		VidmgrStatus: in.VidmgrStatus,
		VidmgrSecret: in.VidmgrSecret,
		VidmgrType:   in.VidmgrType,
		MediasvrType: 2, //设置为独立主机
		Desc:         in.Desc.GetValue(),
	}
	if in.Tags == nil {
		in.Tags = map[string]string{}
	}
	pi.Tags = in.Tags
	return pi, nil
}

func ToVidmgrInfo(ctx context.Context, pi *relationDB.VidmgrInfo, svcCtx *svc.ServiceContext) *vid.VidmgrInfo {

	if pi.VidmgrType == def.Unknown {
		pi.VidmgrType = def.VidmgrTypeZLMedia //当前默认仅支持zlmediakit
	}
	dpi := &vid.VidmgrInfo{
		VidmgrID:     pi.VidmgrID,   //服务id
		VidmgrName:   pi.VidmgrName, //服务名
		VidmgrIpV4:   utils.InetNtoA(pi.VidmgrIpV4),
		VidmgrPort:   pi.VidmgrPort,
		VidmgrType:   pi.VidmgrType,                         //流服务器类型:1:zlmediakit,2:srs,3:monibuca
		VidmgrStatus: pi.VidmgrStatus,                       //服务状态: 1：离线 2：在线  0：未激活
		VidmgrSecret: pi.VidmgrSecret,                       //流服务器注秘钥 只读
		Desc:         &wrappers.StringValue{Value: pi.Desc}, //描述
		CreatedTime:  pi.CreatedTime.Unix(),                 //创建时间
		LastLogin:    pi.LastLogin.Unix(),                   //最后登录时间
		FirstLogin:   pi.FirstLogin.Unix(),                  //首次登录时间
		Tags:         pi.Tags,                               //产品tags
	}

	return dpi
}

func setPoByPb(old *relationDB.VidmgrInfo, data *vid.VidmgrInfo) error {
	if data.VidmgrName != "" {
		old.VidmgrName = data.VidmgrName
	}
	if data.VidmgrIpV4 != "" {
		old.VidmgrIpV4 = utils.InetAtoN(data.VidmgrIpV4)
	}
	if data.VidmgrPort != 0 {
		old.VidmgrPort = data.VidmgrPort
	}
	if data.VidmgrType != 0 {
		old.VidmgrType = data.VidmgrType
	}
	if data.VidmgrSecret != "" {
		old.VidmgrSecret = data.VidmgrSecret
	}
	if data.FirstLogin != 0 {
		old.FirstLogin = time.Unix(data.FirstLogin, 0)
	}
	if data.LastLogin != 0 {
		old.LastLogin = time.Unix(data.LastLogin, 0)
	}
	if data.MediasvrType != 0 {
		old.MediasvrType = data.MediasvrType
	}
	if data.VidmgrStatus != 0 {
		old.VidmgrStatus = data.VidmgrStatus
	}
	if data.Desc != nil {
		old.Desc = data.Desc.GetValue()
	}
	if data.Tags != nil {
		old.Tags = data.Tags
	}
	return nil
}

func ToVidmgrConfigRpc(pi *types.ServerConfig) *relationDB.VidmgrConfig {
	dpi := &relationDB.VidmgrConfig{
		VidmgrID:                       pi.GeneralMediaServerId,
		ApiDebug:                       pi.ApiDebug,
		ApiDefaultSnap:                 pi.ApiDefaultSnap,
		ApiSecret:                      pi.ApiSecret,
		ApiSnapRoot:                    pi.ApiSnapRoot,
		ClusterOriginUrl:               pi.ClusterOriginUrl,
		ClusterRetryCount:              pi.ClusterRetryCount,
		ClusterTimeoutSec:              pi.ClusterTimeoutSec,
		FfmpegBin:                      pi.FfmpegBin,
		FfmpegCmd:                      pi.FfmpegCmd,
		FfmpegLog:                      pi.FfmpegLog,
		FfmpegRestartSec:               pi.FfmpegRestartSec,
		FfmpegSnap:                     pi.FfmpegSnap,
		GeneralCheckNvidiaDev:          pi.GeneralCheckNvidiaDev,
		GeneralEnableVhost:             pi.GeneralEnableVhost,
		GeneralEnableFfmpegLog:         pi.GeneralEnableFfmpegLog,
		GeneralFlowThreshold:           pi.GeneralFlowThreshold,
		GeneralMaxStreamWaitMS:         pi.GeneralMaxStreamWaitMS,
		GeneralMergeWriteMS:            pi.GeneralMergeWriteMS,
		GeneralResetWhenRePlay:         pi.GeneralResetWhenRePlay,
		GeneralStreamNoneReaderDelayMS: pi.GeneralStreamNoneReaderDelayMS,
		GeneralUnreadyFrameCache:       pi.GeneralUnreadyFrameCache,
		GeneralWaitAddTrackMs:          pi.GeneralWaitAddTrackMs,
		GeneralWaitTrackReadyMs:        pi.GeneralWaitTrackReadyMs,
		HlsBroadcastRecordTs:           pi.HlsBroadcastRecordTs,
		HlsDeleteDelaySec:              pi.HlsDeleteDelaySec,
		HlsFileBufSize:                 pi.HlsFileBufSize,
		HlsSegDur:                      pi.HlsSegDur,
		HlsSegKeep:                     pi.HlsSegKeep,
		HlsSegNum:                      pi.HlsSegNum,
		HlsSegRetain:                   pi.HlsSegRetain,
		HookAliveInterval:              pi.HookAliveInterval,
		HookEnable:                     pi.HookEnable,
		HookOnFlowReport:               pi.HookOnFlowReport,
		HookOnHttpAccess:               pi.HookOnHttpAccess,
		HookOnPlay:                     pi.HookOnPlay,
		HookOnPublish:                  pi.HookOnPublish,
		HookOnRecordMp4:                pi.HookOnRecordMp4,
		HookOnRecordTs:                 pi.HookOnRecordTs,
		HookOnRtpServerTimeout:         pi.HookOnRtpServerTimeout,
		HookOnRtspAuth:                 pi.HookOnRtspAuth,
		HookOnRtspRealm:                pi.HookOnRtspRealm,
		HookOnSendRtpStopped:           pi.HookOnSendRtpStopped,
		HookOnServerExited:             pi.HookOnServerExited,
		HookOnServerKeepalive:          pi.HookOnServerKeepalive,
		HookOnServerStarted:            pi.HookOnServerStarted,
		HookOnShellLogin:               pi.HookOnShellLogin,
		HookOnStreamChanged:            pi.HookOnStreamChanged,
		HookOnStreamNoneReader:         pi.HookOnStreamNoneReader,
		HookOnStreamNotFound:           pi.HookOnStreamNotFound,
		HookRetry:                      pi.HookRetry,
		HookRetryDelay:                 pi.HookRetryDelay,
		HookStreamChangedSchemas:       pi.HookStreamChangedSchemas,
		HookTimeoutSec:                 pi.HookTimeoutSec,
		HttpAllowCrossDomains:          pi.HttpAllowCrossDomains,
		HttpAllowIpRange:               pi.HttpAllowIpRange,
		HttpCharSet:                    pi.HttpCharSet,
		HttpDirMenu:                    pi.HttpDirMenu,
		HttpForbidCacheSuffix:          pi.HttpForbidCacheSuffix,
		HttpForwardedIpHeader:          pi.HttpForwardedIpHeader,
		HttpKeepAliveSecond:            pi.HttpKeepAliveSecond,
		HttpMaxReqSize:                 pi.HttpMaxReqSize,
		HttpNotFound:                   pi.HttpNotFound,
		HttpPort:                       pi.HttpPort,
		HttpRootPath:                   pi.HttpRootPath,
		HttpSendBufSize:                pi.HttpSendBufSize,
		HttpSslport:                    pi.HttpSslport,
		HttpVirtualPath:                pi.HttpVirtualPath,
		MulticastAddrMax:               pi.MulticastAddrMax,
		MulticastAddrMin:               pi.MulticastAddrMin,
		MulticastUdpTTL:                pi.MulticastUdpTTL,
		ProtocolAddMuteAudio:           pi.ProtocolAddMuteAudio,
		ProtocolAutoClose:              pi.ProtocolAutoClose,
		ProtocolContinuePushMs:         pi.ProtocolContinuePushMs,
		ProtocolEnableAudio:            pi.ProtocolEnableAudio,
		ProtocolEnableFmp4:             pi.ProtocolEnableFmp4,
		ProtocolEnableHls:              pi.ProtocolEnableHls,
		ProtocolEnableHlsFmp4:          pi.ProtocolEnableHlsFmp4,
		ProtocolEnableMp4:              pi.ProtocolEnableMp4,
		ProtocolEnableRtmp:             pi.ProtocolEnableRtmp,
		ProtocolEnableRtsp:             pi.ProtocolEnableRtsp,
		ProtocolEnableTs:               pi.ProtocolEnableTs,
		ProtocolFmp4Demand:             pi.ProtocolFmp4Demand,
		ProtocolHlsDemand:              pi.ProtocolHlsDemand,
		ProtocolHlsSavePath:            pi.ProtocolHlsSavePath,
		ProtocolModifyStamp:            pi.ProtocolModifyStamp,
		ProtocolMp4AsPlayer:            pi.ProtocolMp4AsPlayer,
		ProtocolMp4MaxSecond:           pi.ProtocolMp4MaxSecond,
		ProtocolMp4SavePath:            pi.ProtocolMp4SavePath,
		ProtocolRtmpDemand:             pi.ProtocolRtmpDemand,
		ProtocolRtspDemand:             pi.ProtocolRtspDemand,
		ProtocolTsDemand:               pi.ProtocolTsDemand,
		RecordAppName:                  pi.RecordAppName,
		RecordFastStart:                pi.RecordFastStart,
		RecordFileBufSize:              pi.RecordFileBufSize,
		RecordFileRepeat:               pi.RecordFileRepeat,
		RecordSampleMS:                 pi.RecordSampleMS,
		RtcExternIP:                    pi.RtcExternIP,
		RtcPort:                        pi.RtcPort,
		RtcPreferredCodecA:             pi.RtcPreferredCodecA,
		RtcPreferredCodecV:             pi.RtcPreferredCodecV,
		RtcRembBitRate:                 pi.RtcRembBitRate,
		RtcTcpPort:                     pi.RtcTcpPort,
		RtcTimeoutSec:                  pi.RtcTimeoutSec,
		RtmpHandshakeSecond:            pi.RtmpHandshakeSecond,
		RtmpKeepAliveSecond:            pi.RtmpKeepAliveSecond,
		RtmpPort:                       pi.RtmpPort,
		RtmpSslport:                    pi.RtmpSslport,
		RtpAudioMtuSize:                pi.RtpAudioMtuSize,
		RtpH264StapA:                   pi.RtpH264StapA,
		RtpLowLatency:                  pi.RtpLowLatency,
		RtpRtpMaxSize:                  pi.RtpRtpMaxSize,
		RtpVideoMtuSize:                pi.RtpVideoMtuSize,
		RtpProxyDumpDir:                pi.RtpProxyDumpDir,
		RtpProxyGopCache:               pi.RtpProxyGopCache,
		RtpProxyH264Pt:                 pi.RtpProxyH264Pt,
		RtpProxyH265Pt:                 pi.RtpProxyH265Pt,
		RtpProxyOpusPt:                 pi.RtpProxyOpusPt,
		RtpProxyPort:                   pi.RtpProxyPort,
		RtpProxyPortRange:              pi.RtpProxyPortRange,
		RtpProxyPsPt:                   pi.RtpProxyPsPt,
		RtpProxyTimeoutSec:             pi.RtpProxyTimeoutSec,
		RtspAuthBasic:                  pi.RtspAuthBasic,
		RtspDirectProxy:                pi.RtspDirectProxy,
		RtspHandshakeSecond:            pi.RtspHandshakeSecond,
		RtspKeepAliveSecond:            pi.RtspKeepAliveSecond,
		RtspLowLatency:                 pi.RtspLowLatency,
		RtspPort:                       pi.RtspPort,
		RtspRtpTransportType:           pi.RtspRtpTransportType,
		RtspSslport:                    pi.RtspSslport,
		ShellMaxReqSize:                pi.ShellMaxReqSize,
		ShellPort:                      pi.ShellPort,
		SrtLatencyMul:                  pi.SrtLatencyMul,
		SrtPktBufSize:                  pi.SrtPktBufSize,
		SrtPort:                        pi.SrtPort,
		SrtTimeoutSec:                  pi.SrtTimeoutSec,
	}
	return dpi
}

// func ToVidmgrConfigRpc(pi *vid.VidmgrConfig) *types.ServerConfig {
func ToVidmgrConfigApi(pi *relationDB.VidmgrConfig) *types.ServerConfig {
	dpi := &types.ServerConfig{
		GeneralMediaServerId:           pi.VidmgrID,
		ApiDebug:                       pi.ApiDebug,
		ApiDefaultSnap:                 pi.ApiDefaultSnap,
		ApiSecret:                      pi.ApiSecret,
		ApiSnapRoot:                    pi.ApiSnapRoot,
		ClusterOriginUrl:               pi.ClusterOriginUrl,
		ClusterRetryCount:              pi.ClusterRetryCount,
		ClusterTimeoutSec:              pi.ClusterTimeoutSec,
		FfmpegBin:                      pi.FfmpegBin,
		FfmpegCmd:                      pi.FfmpegCmd,
		FfmpegLog:                      pi.FfmpegLog,
		FfmpegRestartSec:               pi.FfmpegRestartSec,
		FfmpegSnap:                     pi.FfmpegSnap,
		GeneralCheckNvidiaDev:          pi.GeneralCheckNvidiaDev,
		GeneralEnableVhost:             pi.GeneralEnableVhost,
		GeneralEnableFfmpegLog:         pi.GeneralEnableFfmpegLog,
		GeneralFlowThreshold:           pi.GeneralFlowThreshold,
		GeneralMaxStreamWaitMS:         pi.GeneralMaxStreamWaitMS,
		GeneralMergeWriteMS:            pi.GeneralMergeWriteMS,
		GeneralResetWhenRePlay:         pi.GeneralResetWhenRePlay,
		GeneralStreamNoneReaderDelayMS: pi.GeneralStreamNoneReaderDelayMS,
		GeneralUnreadyFrameCache:       pi.GeneralUnreadyFrameCache,
		GeneralWaitAddTrackMs:          pi.GeneralWaitAddTrackMs,
		GeneralWaitTrackReadyMs:        pi.GeneralWaitTrackReadyMs,
		HlsBroadcastRecordTs:           pi.HlsBroadcastRecordTs,
		HlsDeleteDelaySec:              pi.HlsDeleteDelaySec,
		HlsFileBufSize:                 pi.HlsFileBufSize,
		HlsSegDur:                      pi.HlsSegDur,
		HlsSegKeep:                     pi.HlsSegKeep,
		HlsSegNum:                      pi.HlsSegNum,
		HlsSegRetain:                   pi.HlsSegRetain,
		HookAliveInterval:              pi.HookAliveInterval,
		HookEnable:                     pi.HookEnable,
		HookOnFlowReport:               pi.HookOnFlowReport,
		HookOnHttpAccess:               pi.HookOnHttpAccess,
		HookOnPlay:                     pi.HookOnPlay,
		HookOnPublish:                  pi.HookOnPublish,
		HookOnRecordMp4:                pi.HookOnRecordMp4,
		HookOnRecordTs:                 pi.HookOnRecordTs,
		HookOnRtpServerTimeout:         pi.HookOnRtpServerTimeout,
		HookOnRtspAuth:                 pi.HookOnRtspAuth,
		HookOnRtspRealm:                pi.HookOnRtspRealm,
		HookOnSendRtpStopped:           pi.HookOnSendRtpStopped,
		HookOnServerExited:             pi.HookOnServerExited,
		HookOnServerKeepalive:          pi.HookOnServerKeepalive,
		HookOnServerStarted:            pi.HookOnServerStarted,
		HookOnShellLogin:               pi.HookOnShellLogin,
		HookOnStreamChanged:            pi.HookOnStreamChanged,
		HookOnStreamNoneReader:         pi.HookOnStreamNoneReader,
		HookOnStreamNotFound:           pi.HookOnStreamNotFound,
		HookRetry:                      pi.HookRetry,
		HookRetryDelay:                 pi.HookRetryDelay,
		HookStreamChangedSchemas:       pi.HookStreamChangedSchemas,
		HookTimeoutSec:                 pi.HookTimeoutSec,
		HttpAllowCrossDomains:          pi.HttpAllowCrossDomains,
		HttpAllowIpRange:               pi.HttpAllowIpRange,
		HttpCharSet:                    pi.HttpCharSet,
		HttpDirMenu:                    pi.HttpDirMenu,
		HttpForbidCacheSuffix:          pi.HttpForbidCacheSuffix,
		HttpForwardedIpHeader:          pi.HttpForwardedIpHeader,
		HttpKeepAliveSecond:            pi.HttpKeepAliveSecond,
		HttpMaxReqSize:                 pi.HttpMaxReqSize,
		HttpNotFound:                   pi.HttpNotFound,
		HttpPort:                       pi.HttpPort,
		HttpRootPath:                   pi.HttpRootPath,
		HttpSendBufSize:                pi.HttpSendBufSize,
		HttpSslport:                    pi.HttpSslport,
		HttpVirtualPath:                pi.HttpVirtualPath,
		MulticastAddrMax:               pi.MulticastAddrMax,
		MulticastAddrMin:               pi.MulticastAddrMin,
		MulticastUdpTTL:                pi.MulticastUdpTTL,
		ProtocolAddMuteAudio:           pi.ProtocolAddMuteAudio,
		ProtocolAutoClose:              pi.ProtocolAutoClose,
		ProtocolContinuePushMs:         pi.ProtocolContinuePushMs,
		ProtocolEnableAudio:            pi.ProtocolEnableAudio,
		ProtocolEnableFmp4:             pi.ProtocolEnableFmp4,
		ProtocolEnableHls:              pi.ProtocolEnableHls,
		ProtocolEnableHlsFmp4:          pi.ProtocolEnableHlsFmp4,
		ProtocolEnableMp4:              pi.ProtocolEnableMp4,
		ProtocolEnableRtmp:             pi.ProtocolEnableRtmp,
		ProtocolEnableRtsp:             pi.ProtocolEnableRtsp,
		ProtocolEnableTs:               pi.ProtocolEnableTs,
		ProtocolFmp4Demand:             pi.ProtocolFmp4Demand,
		ProtocolHlsDemand:              pi.ProtocolHlsDemand,
		ProtocolHlsSavePath:            pi.ProtocolHlsSavePath,
		ProtocolModifyStamp:            pi.ProtocolModifyStamp,
		ProtocolMp4AsPlayer:            pi.ProtocolMp4AsPlayer,
		ProtocolMp4MaxSecond:           pi.ProtocolMp4MaxSecond,
		ProtocolMp4SavePath:            pi.ProtocolMp4SavePath,
		ProtocolRtmpDemand:             pi.ProtocolRtmpDemand,
		ProtocolRtspDemand:             pi.ProtocolRtspDemand,
		ProtocolTsDemand:               pi.ProtocolTsDemand,
		RecordAppName:                  pi.RecordAppName,
		RecordFastStart:                pi.RecordFastStart,
		RecordFileBufSize:              pi.RecordFileBufSize,
		RecordFileRepeat:               pi.RecordFileRepeat,
		RecordSampleMS:                 pi.RecordSampleMS,
		RtcExternIP:                    pi.RtcExternIP,
		RtcPort:                        pi.RtcPort,
		RtcPreferredCodecA:             pi.RtcPreferredCodecA,
		RtcPreferredCodecV:             pi.RtcPreferredCodecV,
		RtcRembBitRate:                 pi.RtcRembBitRate,
		RtcTcpPort:                     pi.RtcTcpPort,
		RtcTimeoutSec:                  pi.RtcTimeoutSec,
		RtmpHandshakeSecond:            pi.RtmpHandshakeSecond,
		RtmpKeepAliveSecond:            pi.RtmpKeepAliveSecond,
		RtmpPort:                       pi.RtmpPort,
		RtmpSslport:                    pi.RtmpSslport,
		RtpAudioMtuSize:                pi.RtpAudioMtuSize,
		RtpH264StapA:                   pi.RtpH264StapA,
		RtpLowLatency:                  pi.RtpLowLatency,
		RtpRtpMaxSize:                  pi.RtpRtpMaxSize,
		RtpVideoMtuSize:                pi.RtpVideoMtuSize,
		RtpProxyDumpDir:                pi.RtpProxyDumpDir,
		RtpProxyGopCache:               pi.RtpProxyGopCache,
		RtpProxyH264Pt:                 pi.RtpProxyH264Pt,
		RtpProxyH265Pt:                 pi.RtpProxyH265Pt,
		RtpProxyOpusPt:                 pi.RtpProxyOpusPt,
		RtpProxyPort:                   pi.RtpProxyPort,
		RtpProxyPortRange:              pi.RtpProxyPortRange,
		RtpProxyPsPt:                   pi.RtpProxyPsPt,
		RtpProxyTimeoutSec:             pi.RtpProxyTimeoutSec,
		RtspAuthBasic:                  pi.RtspAuthBasic,
		RtspDirectProxy:                pi.RtspDirectProxy,
		RtspHandshakeSecond:            pi.RtspHandshakeSecond,
		RtspKeepAliveSecond:            pi.RtspKeepAliveSecond,
		RtspLowLatency:                 pi.RtspLowLatency,
		RtspPort:                       pi.RtspPort,
		RtspRtpTransportType:           pi.RtspRtpTransportType,
		RtspSslport:                    pi.RtspSslport,
		ShellMaxReqSize:                pi.ShellMaxReqSize,
		ShellPort:                      pi.ShellPort,
		SrtLatencyMul:                  pi.SrtLatencyMul,
		SrtPktBufSize:                  pi.SrtPktBufSize,
		SrtPort:                        pi.SrtPort,
		SrtTimeoutSec:                  pi.SrtTimeoutSec,
	}
	return dpi
}
