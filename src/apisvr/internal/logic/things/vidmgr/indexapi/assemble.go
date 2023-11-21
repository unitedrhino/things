package indexapi

import (
	"context"
	"encoding/json"
	"fmt"
	zlmediakitapi "github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/vidsvr/pb/vid"
	"net/url"
)

var (
	ADDFFMPEGSOURCE      = "addFFmpegSource"
	ADDSTREAMPROXY       = "addStreamProxy"
	CLOSESTREAM          = "close_stream"
	CLOSESTREAMS         = "close_streams"
	DELFFMPEGSOURCE      = "delFFmpegSource"
	DELSTREAMPROXY       = "delStreamProxy"
	GETALLSESSION        = "getAllSession"
	GETAPILIST           = "getApiList"
	GETMEDIALIST         = "getMediaList"
	GETSERVERCONFIG      = "getServerConfig"
	GETTHREADSLOAD       = "getThreadsLoad"
	GETWORKTHREADSLOAD   = "getWorkThreadsLoad"
	KICKSESSION          = "kick_session"
	KICKSESSIONS         = "kick_sessions"
	RESTARTSERVER        = "restartServer"
	ISRECORDING          = "isRecording"
	SETSERVERCONFIG      = "setServerConfig"
	ISMEDIAONLINE        = "isMediaOnline"
	GETMEDIAINFO         = "getMediaInfo"
	GETRTPINFO           = "getRtpInfo"
	GETMP4RECORDFILE     = "getMp4RecordFile"
	STARTRECORD          = "startRecord"
	STOPRECORD           = "stopRecord"
	GETRECORDSTATUS      = "getRecordStatus"
	STARTSENDRTPPASSIVE  = "startSendRtpPassive"
	GETSNAP              = "getSnap"
	OPENRTPSERVER        = "openRtpServer"
	CLOSERTPSERVER       = "closeRtpServer"
	LISTRTPSERVER        = "listRtpServer"
	STARTSENDRTP         = "startSendRtp"
	STOPSENDRTP          = "stopSendRtp"
	GETSTATISTIC         = "getStatistic"
	ADDSTREAMPUSHERPROXY = "addStreamPusherProxy"
	DELSTREAMPUSHERPROXY = "delStreamPusherProxy"
	VERSION              = "version"
	GETMEDIAPLAYERLIST   = "getMediaPlayerList"
)

func proxyMediaServer(ctx context.Context, svcCtx *svc.ServiceContext, preUrl string, vidmgrID string) (data []byte, err error) {
	vidResp, err := svcCtx.VidmgrM.VidmgrInfoRead(ctx, &vid.VidmgrInfoReadReq{
		VidmgrtID: vidmgrID,
	})
	if err != nil {
		er := errors.Fmt(err)
		fmt.Print("%s rpc.VidmgrInfoRead  err=%+v", utils.FuncName(), er)
		return nil, er
	}
	mediaSrv := zlmediakitapi.NewMeidaServer(vidResp)
	values := make(url.Values)
	values.Add("secret", vidResp.VidmgrSecret)
	vidRecv, error := mediaSrv.PostMediaServer(preUrl, values)
	if error != nil {
		er := errors.Fmt(error)
		fmt.Print("%s rpc.PostMediaServer  err=%+v", utils.FuncName(), er)
		return nil, er
	}
	return vidRecv, nil
}

func proxySetMediaServer(ctx context.Context, svcCtx *svc.ServiceContext, preUrl string, vidmgrID string, values []byte) (data []byte, err error) {
	vidResp, err := svcCtx.VidmgrM.VidmgrInfoRead(ctx, &vid.VidmgrInfoReadReq{
		VidmgrtID: vidmgrID,
	})
	if err != nil {
		er := errors.Fmt(err)
		fmt.Print("%s rpc.VidmgrInfoRead  err=%+v", utils.FuncName(), er)
		return nil, er
	}
	mediaSrv := zlmediakitapi.NewMeidaServer(vidResp)
	var tdata map[string]interface{}
	err = json.Unmarshal(values, &tdata)
	tdata["secret"] = vidResp.VidmgrSecret
	values, err = json.Marshal(tdata)
	if err != nil {
		er := errors.Fmt(err)
		fmt.Print("%s map string phares failed  err=%+v", utils.FuncName(), er)
		return nil, er
	}
	vidRecv, error := mediaSrv.PostMediaServerJson(preUrl, values)
	if error != nil {
		er := errors.Fmt(error)
		fmt.Print("%s rpc.PostMediaServer  err=%+v", utils.FuncName(), er)
		return nil, er
	}
	return vidRecv, nil
}

const HostIP = "192.168.10.33"

func SetDefaultConfig(svcCtx *svc.ServiceContext, config *types.ServerConfig) {
	config.ApiDebug = "1"
	config.HookEnable = "1"
	config.HookOnFlowReport = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onFlowReport", HostIP, svcCtx.Config.Port)
	config.HookOnHttpAccess = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onHttpAccess", HostIP, svcCtx.Config.Port)
	config.HookOnPlay = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onPlay", HostIP, svcCtx.Config.Port)
	config.HookOnPublish = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onPublish", HostIP, svcCtx.Config.Port)
	config.HookOnRecordMp4 = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onRecordMp4", HostIP, svcCtx.Config.Port)
	config.HookOnRecordTs = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onRecordTs", HostIP, svcCtx.Config.Port)
	config.HookOnRtpServerTimeout = fmt.Sprintf("https://%s:%d/api/v1/things/vidmgr/hooks/onRtpServerTimeout", HostIP, svcCtx.Config.Port)
	config.HookOnRtspAuth = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onRtspAuth", HostIP, svcCtx.Config.Port)
	config.HookOnRtspRealm = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onRtspRealm", HostIP, svcCtx.Config.Port)
	config.HookOnSendRtpStopped = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onSendRtpStopped", HostIP, svcCtx.Config.Port)
	config.HookOnServerExited = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onServerExited", HostIP, svcCtx.Config.Port)
	config.HookOnServerKeepalive = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onServerKeepalive", HostIP, svcCtx.Config.Port)
	config.HookOnServerStarted = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onServerStarted", HostIP, svcCtx.Config.Port)
	config.HookOnShellLogin = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onShellLogin", HostIP, svcCtx.Config.Port)
	config.HookOnStreamChanged = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onStreamChanged", HostIP, svcCtx.Config.Port)
	config.HookOnStreamNoneReader = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onStreamNoneReader", HostIP, svcCtx.Config.Port)
	config.HookOnStreamNotFound = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onStreamNotFound", HostIP, svcCtx.Config.Port)
}
