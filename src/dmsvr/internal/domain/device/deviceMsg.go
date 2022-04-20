// Package device 设备发送来的消息解析
package device

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/ddsvr/ddExport"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type PublishMsg struct {
	Topic      string
	Payload    []byte
	Timestamp  time.Time
	ProductID  string
	DeviceName string
}

//连接和断连消息信息
type ConnectMsg struct {
	ClientID  string
	Username  string
	Timestamp time.Time
	Address   string
	Action    string
	Reason    string
}

func GetDevConnMsg(ctx context.Context, data []byte) (*ConnectMsg, error) {
	logInfo := ddExport.DevConn{}
	err := json.Unmarshal(data, &logInfo)
	if err != nil {
		logx.WithContext(ctx).Error("getDevConnMsg", string(data), err)
		return nil, err
	}
	ele := ConnectMsg{
		ClientID:  logInfo.ClientID,
		Username:  logInfo.UserName,
		Timestamp: time.UnixMilli(logInfo.Timestamp),
		Address:   logInfo.Address,
		Action:    logInfo.Action,
		Reason:    logInfo.Reason,
	}
	return &ele, nil
}

func GetDevPublish(ctx context.Context, data []byte) (*PublishMsg, error) {
	pubInfo := ddExport.DevPublish{}
	err := json.Unmarshal(data, &pubInfo)
	if err != nil {
		logx.WithContext(ctx).Error("GetDevPublish", string(data), err)
		return nil, err
	}
	ele := PublishMsg{
		Topic:      pubInfo.Topic,
		Payload:    pubInfo.Payload,
		Timestamp:  time.UnixMilli(pubInfo.Timestamp),
		ProductID:  pubInfo.ProductID,
		DeviceName: pubInfo.DeviceName,
	}
	return &ele, nil
}
