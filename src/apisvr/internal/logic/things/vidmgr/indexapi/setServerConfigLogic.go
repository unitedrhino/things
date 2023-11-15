package indexapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetServerConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetServerConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetServerConfigLogic {
	return &SetServerConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func setDefaultConfig(svcCtx *svc.ServiceContext, config *types.IndexApiServerConfig) {
	config.ApiDebug = "1"
	config.HookEnable = "1"
	config.HookOnFlowReport = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onFlowReport", utils.GetHostIp(), svcCtx.Config.Port)
	config.HookOnHttpAccess = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onHttpAccess", utils.GetHostIp(), svcCtx.Config.Port)
	config.HookOnPlay = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onPlay", utils.GetHostIp(), svcCtx.Config.Port)
	config.HookOnPublish = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onPublish", utils.GetHostIp(), svcCtx.Config.Port)
	config.HookOnRecordMp4 = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onRecordMp4", utils.GetHostIp(), svcCtx.Config.Port)
	config.HookOnRecordTs = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onRecordTs", utils.GetHostIp(), svcCtx.Config.Port)
	config.HookOnRtpServerTimeout = fmt.Sprintf("https://%s:%d/api/v1/things/vidmgr/hooks/onRtpServerTimeout", utils.GetHostIp(), svcCtx.Config.Port)
	config.HookOnRtspAuth = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onRtspAuth", utils.GetHostIp(), svcCtx.Config.Port)
	config.HookOnRtspRealm = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onRtspRealm", utils.GetHostIp(), svcCtx.Config.Port)
	config.HookOnSendRtpStopped = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onSendRtpStopped", utils.GetHostIp(), svcCtx.Config.Port)
	config.HookOnServerExited = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onServerExited", utils.GetHostIp(), svcCtx.Config.Port)
	config.HookOnServerKeepalive = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onServerKeepalive", utils.GetHostIp(), svcCtx.Config.Port)
	config.HookOnServerStarted = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onServerStarted", utils.GetHostIp(), svcCtx.Config.Port)
	config.HookOnShellLogin = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onShellLogin", utils.GetHostIp(), svcCtx.Config.Port)
	config.HookOnStreamChanged = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onStreamChanged", utils.GetHostIp(), svcCtx.Config.Port)
	config.HookOnStreamNoneReader = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onStreamNoneReader", utils.GetHostIp(), svcCtx.Config.Port)
	config.HookOnStreamNotFound = fmt.Sprintf("http://%s:%d/api/v1/things/vidmgr/hooks/onStreamNotFound", utils.GetHostIp(), svcCtx.Config.Port)
}

func (l *SetServerConfigLogic) SetServerConfig(req *types.IndexApiSetServerConfigReq) (resp *types.IndexApiSetServerConfigResp, err error) {
	// todo: add your logic here and delete this line
	//serverConfig := types.IndexApiServerConfig{}
	fmt.Println("[airgens] ID  ", req.VidmgrID)
	fmt.Println("[airgens] Data", req.Data)
	//json.Marshal()
	////data, err := proxySetMediaServer(l.ctx, l.svcCtx, SETSERVERCONFIG, req.VidmgrID)
	dataRecv := new(types.IndexApiSetServerConfigResp)

	strmConfig := new(types.IndexApiServerConfig)
	err = json.Unmarshal([]byte(req.Data), strmConfig)
	if err != nil {
		fmt.Println("parse Json failed:", err)
		return dataRecv, err
	}
	strmConfig.GeneralMediaServerId = req.VidmgrID
	setDefaultConfig(l.svcCtx, strmConfig)
	fmt.Println("[_______IndexApiServerConfig________]IndexApiServerConfig struct:", strmConfig)
	//set default
	byte4, err := json.Marshal(strmConfig)
	fmt.Println("strmConfig TOjSON:", string(byte4))

	mdata, err := proxySetMediaServer(l.ctx, l.svcCtx, SETSERVERCONFIG, req.VidmgrID, byte4)
	fmt.Println("proxySetMediaServer:", string(mdata))

	//fmt.Println("IP:", getHostIp() /*l.svcCtx.Config.Host*/, " Port:", l.svcCtx.Config.Port)
	//json.Unmarshal(data, dataRecv)
	err = json.Unmarshal(mdata, dataRecv)
	if err != nil {
		fmt.Println("parse Json failed:", err)
		return dataRecv, err
	}
	return dataRecv, err
}
