package indexapi

import (
	"context"
	"fmt"
	zlmediakitapi "github.com/i-Things/things/shared/api/zlmediakit"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/svc"
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
	//fmt.Println("---------vidResp-----", vidResp)
	mediaSrv := zlmediakitapi.NewMeidaServer(vidResp)
	//fmt.Println("______mediaSrv______", mediaSrv.PreUrl)
	//var values url.Values
	values := make(url.Values)
	vidRecv, error := mediaSrv.PostMediaServer(preUrl, values)
	if error != nil {
		er := errors.Fmt(error)
		fmt.Print("%s rpc.PostMediaServer  err=%+v", utils.FuncName(), er)
		return nil, er
	}
	return vidRecv, nil
}
