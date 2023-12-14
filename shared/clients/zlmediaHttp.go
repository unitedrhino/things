package clients

//部分门接口直接http访问
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"io"
	"net/http"
)

var PreUrl string

const (
	VIDMGRTIMEOUT = 30
)

// 外部用
type SvcZlmedia struct {
	IP     string `json:"ip"`
	Port   int64  `json:"port"`
	Secret string `json:"Secret"`
}

const (
	MEDIA_DOCKER = 1 << iota
	MEDIA_OTHER
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

// 内部用 ZLMediakit初始化数据结构
type mediaConfig struct {
	ipv4   string
	port   int64
	preUrl string
}

func NewMeidaServer(ip string, port int64) *mediaConfig {
	return &mediaConfig{
		ipv4:   ip,
		port:   port,
		preUrl: fmt.Sprintf("http://%s:%d/index/api/", ip, port),
	}
}

func (f *mediaConfig) PostMediaServerJson(strurl string, values []byte) (data []byte, err error) {
	//fmt.Println("[---------------]PostMediaServer -", f.preUrl+strurl, "param:", values)
	request, error := http.NewRequest("POST", f.preUrl+strurl, bytes.NewBuffer(values))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		fmt.Println(err)
		return nil, err
	}
	defer response.Body.Close()
	//fmt.Println("response Status:", response.Status)
	//fmt.Println("response Headers:", response.Header)
	body, err := io.ReadAll(response.Body)
	//fmt.Println("response Body:", string(body))
	return body, err
}

func ProxyMediaServer(cmd string, mgr *SvcZlmedia, values []byte) (data []byte, err error) {
	mediaSrv := NewMeidaServer(mgr.IP, mgr.Port)
	tdata := make(map[string]interface{})
	err = json.Unmarshal(values, &tdata)
	tdata["secret"] = mgr.Secret
	values, err = json.Marshal(tdata)
	if err != nil {
		er := errors.Fmt(err).AddMsg("构造服务数据失败")
		fmt.Print("%s map string phares failed  err=%+v", utils.FuncName(), er)
		return nil, err
	}
	//fmt.Println("[****test***] ", string(values))
	vidRecv, error := mediaSrv.PostMediaServerJson(cmd, values)
	if error != nil {
		er := errors.Fmt(error).AddMsg("服务不在线")
		fmt.Print("%s rpc.PostMediaServer  err=%+v", utils.FuncName(), er)
		return nil, error
	}
	return vidRecv, nil
}
