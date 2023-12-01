package index

import (
	"context"
	"encoding/json"
	"fmt"
	zlmediakitapi "github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/types"
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

func proxySetMediaServer(ctx context.Context, preUrl string, vidmgrID string, values []byte) (data []byte, err error) {
	pi, err := relationDB.NewVidmgrInfoRepo(ctx).FindOneByFilter(ctx, relationDB.VidmgrFilter{
		VidmgrIDs: []string{vidmgrID},
	})
	if err != nil {
		er := errors.Fmt(err).AddMsg("数据库查询失败")
		fmt.Print("%s rpc.VidmgrInfoRead  err=%+v", utils.FuncName(), er)
		return nil, er
	}
	if pi != nil {
		mediaSrv := zlmediakitapi.NewMeidaServer(utils.InetNtoA(pi.VidmgrIpV4), pi.VidmgrPort)
		tdata := make(map[string]interface{})
		err = json.Unmarshal(values, &tdata)
		tdata["secret"] = pi.VidmgrSecret
		values, err = json.Marshal(tdata)
		if err != nil {
			er := errors.Fmt(err).AddMsg("构造服务数据失败")
			fmt.Print("%s map string phares failed  err=%+v", utils.FuncName(), er)
			return nil, er
		}
		vidRecv, error := mediaSrv.PostMediaServerJson(preUrl, values)
		if error != nil {
			er := errors.Fmt(error).AddMsg("服务不在线")
			fmt.Print("%s rpc.PostMediaServer  err=%+v", utils.FuncName(), er)
			return nil, er
		}
		return vidRecv, nil
	}
	return nil, errors.MediaNotfoundError.AddMsg("服务未找到")
}

func SetDefaultConfig(hostip string, hostport int64, config *types.ServerConfig) {
	//hostport
	config.ApiDebug = "1"
	config.HookEnable = "1"
	config.HookOnFlowReport = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onFlowReport", hostip, hostport)
	config.HookOnHttpAccess = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onHttpAccess", hostip, hostport)
	config.HookOnPlay = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onPlay", hostip, hostport)
	config.HookOnPublish = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onPublish", hostip, hostport)
	config.HookOnRecordMp4 = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onRecordMp4", hostip, hostport)
	config.HookOnRecordTs = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onRecordTs", hostip, hostport)
	config.HookOnRtpServerTimeout = fmt.Sprintf("https://%s:%d/api/v1/things/vidmgr/hooks/onRtpServerTimeout", hostip, hostport)
	config.HookOnRtspAuth = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onRtspAuth", hostip, hostport)
	config.HookOnRtspRealm = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onRtspRealm", hostip, hostport)
	config.HookOnSendRtpStopped = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onSendRtpStopped", hostip, hostport)
	config.HookOnServerExited = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onServerExited", hostip, hostport)
	config.HookOnServerKeepalive = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onServerKeepalive", hostip, hostport)
	config.HookOnServerStarted = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onServerStarted", hostip, hostport)
	config.HookOnShellLogin = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onShellLogin", hostip, hostport)
	config.HookOnStreamChanged = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onStreamChanged", hostip, hostport)
	config.HookOnStreamNoneReader = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onStreamNoneReader", hostip, hostport)
	config.HookOnStreamNotFound = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onStreamNotFound", hostip, hostport)
}
