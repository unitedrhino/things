package deviceSend

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/ddsvr/ddExport"
	"github.com/zeromicro/go-zero/core/logx"
)

//Elements 设备发送的所有属性
type Elements struct {
	ProductID  string
	DeviceName string
	ClientID   string
	Username   string
	Address    string
	Topic      string
	Payload    []byte
	Timestamp  int64
	Action     string
	Reason     string
}

func GetDevPublish(ctx context.Context, data []byte) (*Elements, error) {
	pubInfo := ddExport.DevPublish{}
	err := json.Unmarshal(data, &pubInfo)
	if err != nil {
		logx.WithContext(ctx).Error("GetDevPublish", string(data), err)
		return nil, err
	}
	ele := Elements{
		Topic:      pubInfo.Topic,
		Payload:    pubInfo.Payload,
		Timestamp:  pubInfo.Timestamp,
		ProductID:  pubInfo.ProductID,
		DeviceName: pubInfo.DeviceName,
	}
	return &ele, nil
}
