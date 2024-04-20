package devicemsglogic

import (
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceLog"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func ToDataHubLogIndex(log *deviceLog.Hub) *dm.HubLogInfo {
	return &dm.HubLogInfo{
		Timestamp:  log.Timestamp.UnixMilli(),
		Action:     log.Action,
		RequestID:  log.RequestID,
		TraceID:    log.TraceID,
		Topic:      log.Topic,
		Content:    log.Content,
		ResultCode: log.ResultCode,
	}
}

func ToDataSendLogIndex(log *deviceLog.Send) *dm.SendLogInfo {
	return &dm.SendLogInfo{
		Timestamp:  log.Timestamp.UnixMilli(),
		Account:    log.Account,
		UserID:     log.UserID,
		ProductID:  log.ProductID,
		DeviceName: log.DeviceName,
		Action:     log.Action,
		DataID:     log.DataID,
		TraceID:    log.TraceID,
		Content:    log.Content,
		ResultCode: log.ResultCode,
	}
}

func ToDataStatusLogIndex(log *deviceLog.Status) *dm.StatusLogInfo {
	return &dm.StatusLogInfo{
		Timestamp:  log.Timestamp.UnixMilli(),
		Status:     log.Status,
		ProductID:  log.ProductID,
		DeviceName: log.DeviceName,
	}
}

// SDK调试日志
func ToDataSdkLogIndex(log *deviceLog.SDK) *dm.SdkLogInfo {
	return &dm.SdkLogInfo{
		Timestamp: log.Timestamp.UnixMilli(),
		Loglevel:  log.LogLevel,
		Content:   log.Content,
	}
}
