package devicemsglogic

import (
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/pb/di"
)

func ToDataHubLogIndex(log *deviceMsg.HubLog) *di.HubLogIndex {
	return &di.HubLogIndex{
		Timestamp:  log.Timestamp.UnixMilli(),
		Action:     log.Action,
		RequestID:  log.RequestID,
		TranceID:   log.TranceID,
		Topic:      log.Topic,
		Content:    log.Content,
		ResultType: log.ResultType,
	}
}

//SDK调试日志
func ToDataSdkLogIndex(log *deviceMsg.SDKLog) *di.SdkLogIndex {
	return &di.SdkLogIndex{
		Timestamp: log.Timestamp.UnixMilli(),
		Loglevel:  log.LogLevel,
		Content:   log.Content,
	}
}
