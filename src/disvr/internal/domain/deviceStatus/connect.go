// Package device 设备发送来的消息解析
package deviceStatus

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/devices"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

//连接和断连消息信息
type ConnectMsg struct {
	ClientID  string
	Username  string
	Timestamp time.Time
	Address   string
	Action    string
	Reason    string
}

const (
	ConnectStatus    = "connected"
	DisConnectStatus = "disconnected"
)

func GetDevConnMsg(ctx context.Context, data []byte) (*ConnectMsg, error) {
	logInfo := devices.DevConn{}
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
