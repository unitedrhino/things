package indexapi

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CloseRtpServerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloseRtpServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloseRtpServerLogic {
	return &CloseRtpServerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloseRtpServerLogic) CloseRtpServer(req *types.IndexApiReq) (resp *types.IndexApiCloseRtpServerResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, CLOSERTPSERVER, req.VidmgrID)
	dataRecv := new(types.IndexApiCloseRtpServerResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func SetDefaultConfig(svcCtx *svc.ServiceContext, config *types.ServerConfig) {
	HostIP := svcCtx.Config.OssConf.CustomHost
	HPort := svcCtx.Config.Port
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
