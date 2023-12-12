package hooks

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/i-Things/things/src/vidsvr/pb/vid"
)

func ToVidmgrConfigRpc(pi *types.ServerConfig) *vid.VidmgrConfig {
	dpi := &vid.VidmgrConfig{
		GeneralMediaServerId:           pi.GeneralMediaServerId,
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

func ToVidmgrConfigApi(pi *vid.VidmgrConfig) *types.ServerConfig {
	dpi := &types.ServerConfig{
		GeneralMediaServerId:           pi.GeneralMediaServerId,
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
func ToVidmgrTrackRpc(in *types.StreamTrack) *relationDB.StreamTrack {
	pi := &relationDB.StreamTrack{
		Channels:    in.Channels,
		CodecId:     in.CodecId,
		CodecIdName: in.CodecIdName,
		CodecType:   in.CodecType,
		Ready:       in.Ready,
		SampleBit:   in.SampleBit,
		SampleRate:  in.SampleRate,
		Fps:         in.Fps,
		Height:      in.Height,
		Width:       in.Width,
	}
	return pi
}

func TovidmgrTracksRpc(pi []types.StreamTrack) []*relationDB.StreamTrack {
	if len(pi) > 0 {
		info := make([]*relationDB.StreamTrack, 0, len(pi))
		for _, v := range pi {
			info = append(info, ToVidmgrTrackRpc(&v))
		}
		return info
	}
	return nil

}

func ToVidmgrStreamRpc(pi *types.HooksApiStreamChangedRep) *relationDB.VidmgrStream {
	dpi := &relationDB.VidmgrStream{
		VidmgrID: pi.MediaServerId,

		Stream: pi.Stream,
		Vhost:  pi.Vhost,
		App:    pi.App,

		Identifier: pi.OriginSock.Identifier,
		LocalIP:    utils.InetAtoN(pi.OriginSock.LocalIp),
		LocalPort:  pi.OriginSock.LocalPort,
		PeerIP:     utils.InetAtoN(pi.OriginSock.PeerIp),
		PeerPort:   pi.OriginSock.PeerPort,

		OriginType: pi.OriginType,
		OriginUrl:  pi.OriginUrl,
		OriginStr:  pi.OriginTypeStr,

		ReaderCount:      pi.ReaderCount,
		TotalReaderCount: pi.TotalReaderCount,
		Tracks:           TovidmgrTracksRpc(pi.Tracks),
	}
	return dpi
}

func ToVidmgrTrackApi(in *vid.StreamTrack) *types.StreamTrack {
	pi := &types.StreamTrack{
		Channels:    in.Channels,
		CodecId:     in.CodecId,
		CodecIdName: in.CodecIdName,
		CodecType:   in.CodecType,
		Ready:       in.Ready,
		SampleBit:   in.SampleBit,
		SampleRate:  in.SampleRate,
		Fps:         in.Fps,
		Height:      in.Height,
		Width:       in.Width,
	}
	return pi
}

func TovidmgrTracksApi(pi []vid.StreamTrack) []*types.StreamTrack {
	if len(pi) > 0 {
		info := make([]*types.StreamTrack, 0, len(pi))
		for _, v := range pi {
			info = append(info, ToVidmgrTrackApi(&v))
		}
		return info
	}
	return nil
}

/*
 *Protocol为视频协议
 *当前协议支持类型有 rtmp/rtsp/ts/fmp4/hls/hls.fmp4/
 *分别用一个bit位来表示一个协议
 *对应关系:
 *          bit位        	5	          4   	  3   	  2      1      0
 *                          hls.fmp4      hls     fmp4    ts     rtsp   rtmp
 */
//                          1            1        0      1      0      1
const (
	RTMP = 1 << iota
	RTSP
	TS
	FMP4
	HLS
	HLS_FMP4
)

func GetProtocol(schema string) uint32 {
	var val uint32
	switch schema {
	case "rtmp":
		val = RTMP
	case "rtsp":
		val = RTSP
	case "ts":
		val = TS
	case "fmp4":
		val = FMP4
	case "hls":
		val = HLS
	case "hls.fmp4":
		val = HLS_FMP4
	default:
		val = 0
	}
	return val
}

const (
	UNKNOWN = iota
	RTMP_PUSH
	RTSP_PUSH
	RTP_PUSH
	PULL
	FFMPEG_PULL
	MP4_VOD
	DEVICE_CHN
	RTC_PUSH
)

//产生源类型，包括 unknown = 0,rtmp_push=1,rtsp_push=2,rtp_push=3,pull=4,ffmpeg_pull=5,mp4_vod=6,device_chn=7,rtc_push=8
