package zlmedia

/*
func IndexApi(mgr *types.SvcZlmedia, cmd sring, indata []byte) (resp []byte, err error) {
	data, err := ProxyMediaServer(ADDFFMPEGSOURCE, mgr, indata)
	return dataRecv, err
}

// bytetmp := make([]byte, 0)
func AddFFmpegSource(mgr *types.SvcZlmedia, indata []byte) (resp *types.IndexApiAddFFmpegSourceResp, err error) {
	data, err := ProxyMediaServer(ADDFFMPEGSOURCE, mgr, indata)
	dataRecv := new(types.IndexApiAddFFmpegSourceResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func AddStreamProxy(mgr *types.SvcZlmedia, indata []byte) (resp *types.IndexApiAddStreamProxyResp, err error) {
	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ADDSTREAMPROXY, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiAddStreamProxyResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func AddStreamPusherProxy(mgr *types.SvcZlmedia, indata []byte) (resp *types.IndexApiAddStreamPusherProxyResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, ADDSTREAMPUSHERPROXY, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiAddStreamPusherProxyResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func CloseRtpServer(mgr *types.SvcZlmedia, indata []byte) (resp *types.IndexApiCloseRtpServerResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, CLOSERTPSERVER, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiCloseRtpServerResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func CloseStream(mgr *types.SvcZlmedia, indata []byte) (resp *types.IndexApiCloseStreamResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, CLOSESTREAM, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiCloseStreamResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func CloseStreams(mgr *types.SvcZlmedia, indata []byte) (resp *types.IndexApiCloseStreamsResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, CLOSESTREAMS, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiCloseStreamsResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func DelFFmpegSource(mgr *types.SvcZlmedia, indata []byte) (resp *types.IndexApiDelFFmpegSourceResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, DELFFMPEGSOURCE, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiDelFFmpegSourceResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func DelStreamProxy(mgr *types.SvcZlmedia, indata []byte) (resp *types.IndexApiDelStreamProxyResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, DELSTREAMPROXY, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiDelStreamProxyResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func DelStreamPusherProxy(mgr *types.SvcZlmedia, indata []byte) (resp *types.IndexApiDelStreamProxyResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, DELSTREAMPUSHERPROXY, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiDelStreamProxyResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
func GetAllSession(mgr *types.SvcZlmedia, indata []byte) (resp *types.IndexApiAllSessionResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, GETALLSESSION, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiAllSessionResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func GetApiList(mgr *types.SvcZlmedia, indata []byte) (resp *types.IndexApiListResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, GETAPILIST, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiListResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func GetMediaInfo(ctx context.Context, req *types.IndexApiReq) (resp *types.IndexApiMediaInfoResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, GETMEDIAINFO, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiMediaInfoResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err

}

func GetMediaList(ctx context.Context, req *types.IndexApiReq) (resp *types.IndexApiMediaListResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, GETMEDIALIST, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiMediaListResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func GetMediaPlayerList(ctx context.Context, req *types.IndexApiReq) (resp *types.IndexApiMediaPlayerListResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, GETMEDIAPLAYERLIST, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiMediaPlayerListResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func GetMp4RecordFile(ctx context.Context, req *types.IndexApiReq) (resp *types.IndexApiMp4RecordFileResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, GETMP4RECORDFILE, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiMp4RecordFileResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func GetRtpInfo(ctx context.Context, req *types.IndexApiReq) (resp *types.IndexApiRtpInfoResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, GETRTPINFO, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiRtpInfoResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func GetServerConfig(ctx context.Context, req *types.IndexApiReq) (resp *types.IndexApiServerConfigResp, err error) {

	bytetmp := make([]byte, 1024)
	data, err := ProxyMediaServer(ctx, GETSERVERCONFIG, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiServerConfigResp)
	json.Unmarshal(data, dataRecv)
	fmt.Println("GetServerConfig:", dataRecv)
	return dataRecv, err
}

func GetSnap(ctx context.Context, req *types.IndexApiReq) (resp *types.IndexApiSnapResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, GETSNAP, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiSnapResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func GetStatistic(ctx context.Context, req *types.IndexApiReq) (resp *types.IndexApiStatisticResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, GETSTATISTIC, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiStatisticResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func GetThreadsLoad(ctx context.Context, req *types.IndexApiReq) (resp *types.IndexApiThreadLoadResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, GETTHREADSLOAD, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiThreadLoadResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func GetWorkThreadsLoad(ctx context.Context, req *types.IndexApiReq) (resp *types.IndexApiWorkThreadLoadResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, GETWORKTHREADSLOAD, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiWorkThreadLoadResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func IsMediaOnline(ctx context.Context, req *types.IndexApiReq) (resp *types.IndexApiIsMediaOnlineResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, ISMEDIAONLINE, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiIsMediaOnlineResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func IsRecording(ctx context.Context, req *types.IndexApiReq) (resp *types.IndexApiIsRecordingResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, ISRECORDING, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiIsRecordingResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func KickSession(ctx context.Context, req *types.IndexApiReq) (resp *types.IndexApiKickSessionResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, KICKSESSION, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiKickSessionResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func KickSessions(ctx context.Context, req *types.IndexApiReq) (resp *types.IndexApiKickSessionsResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, KICKSESSIONS, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiKickSessionsResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func ListRtpServer(ctx context.Context, req *types.IndexApiReq) (resp *types.IndexApiListRtpServerResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, LISTRTPSERVER, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiListRtpServerResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err

}

func OpenRtpServer(ctx context.Context, req *types.IndexApiReq) (resp *types.IndexApiOpenRtpServerResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, OPENRTPSERVER, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiOpenRtpServerResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func RestartServer(ctx context.Context, req *types.IndexApiReq) (resp *types.IndexApiRestartServerResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(ctx, RESTARTSERVER, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiRestartServerResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func SetServerConfig(ctx context.Context, req *types.IndexApiSetServerConfigReq) (resp *types.IndexApiSetServerConfigResp, err error) {

	dataRecv := new(types.IndexApiSetServerConfigResp)
	strmConfig := new(types.ServerConfig)
	err = json.Unmarshal([]byte(req.Data), strmConfig)
	if err != nil {
		fmt.Println("parse Json failed:", err)
		return dataRecv, err
	}
	strmConfig.GeneralMediaServerId = req.VidmgrID
	SetDefaultConfig(l.svcCtx.Config.Mediakit.Host, l.svcCtx.Config.Mediakit.Port, strmConfig)
	//set default
	byte4, err := json.Marshal(strmConfig)
	mdata, err := ProxyMediaServer(l.ctx, SETSERVERCONFIG, req.VidmgrID, byte4)
	err = json.Unmarshal(mdata, dataRecv)
	if err != nil {
		fmt.Println("parse Json failed:", err)
		return dataRecv, err
	}
	return dataRecv, err
}

func StartRecord(req *types.IndexApiReq) (resp *types.IndexApiStartRecordResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(l.ctx, STARTRECORD, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiStartRecordResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func StartSendRtp(req *types.IndexApiReq) (resp *types.IndexApiStartSendRtpResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(l.ctx, STARTSENDRTP, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiStartSendRtpResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err

}

func StartSendRtpPassive(req *types.IndexApiReq) (resp *types.IndexApiStartSendRtpPassiveResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(l.ctx, STARTSENDRTPPASSIVE, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiStartSendRtpPassiveResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func StopRecord(req *types.IndexApiReq) (resp *types.IndexApiStopRecordResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(l.ctx, STOPRECORD, req.VidmgrID, bytetmp)

	dataRecv := new(types.IndexApiStopRecordResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func StopSendRtp(req *types.IndexApiReq) (resp *types.IndexApiStopSendRtpResp, err error) {

	bytetmp := make([]byte, 0)
	data, err := ProxyMediaServer(l.ctx, STOPSENDRTP, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiStopSendRtpResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}

func Version(ctx context.Context) (resp *types.IndexApiVersionResp, err error) {

	bytetmp := make([]byte, 0)
	fmt.Println("***Version ****")
	data, err := ProxyMediaServer(ctx, VERSION, req.VidmgrID, bytetmp)
	if err != nil {
		fmt.Println("***proxyMediaServer Error ****")
		er := errors.Fmt(err)
		fmt.Print("%s proxyMediaServer  err=%+v", utils.FuncName(), er)
		return nil, er
	}
	dataRecv := new(types.IndexApiVersionResp)
	fmt.Println(string(data))
	fmt.Println(dataRecv)
	json.Unmarshal(data, dataRecv)
	return dataRecv, nil
}
*/
