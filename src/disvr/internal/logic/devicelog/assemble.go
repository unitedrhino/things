package deviceloglogic

import (
	"database/sql"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/pb/di"
)

func GetNullTime(time sql.NullTime) int64 {
	if time.Valid == false {
		return 0
	}
	return time.Time.Unix()
}

func ToDataHubLogIndex(log *deviceMsg.HubLog) *di.DataHubLogIndex {
	return &di.DataHubLogIndex{
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
func ToDataSdkLogIndex(log *deviceMsg.SDKLog) *di.DataSdkLogIndex {
	return &di.DataSdkLogIndex{
		Timestamp: log.Timestamp.UnixMilli(),
		Loglevel:  log.LogLevel,
		Content:   log.Content,
	}
}
