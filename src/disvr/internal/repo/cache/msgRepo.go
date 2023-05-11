package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/zeromicro/go-zero/core/stores/kv"
)

const (
	msgExpr = 10 * 60
)

func genDeviceMsgKey(msgType string, handle string, Type string, device devices.Core, clientToken string) string {
	return fmt.Sprintf("device:%s:%s:%s:%s:%s:%s",
		handle, Type, msgType, device.ProductID, device.DeviceName, clientToken)
}

func SetDeviceMsg(ctx context.Context, store kv.Store, msgType string, req *deviceMsg.PublishMsg, clientToken string) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}
	err = store.SetexCtx(ctx, genDeviceMsgKey(msgType, req.Handle, req.Type, devices.Core{
		ProductID:  req.ProductID,
		DeviceName: req.DeviceName,
	}, clientToken), string(payload), msgExpr)
	return err
}

func GetDeviceMsg[reqType any](ctx context.Context, store kv.Store, msgType string /*请求还是回复*/, handle string, Type string, device devices.Core, clientToken string) (*reqType, error) {
	val, err := store.GetCtx(ctx, genDeviceMsgKey(msgType, handle, Type, devices.Core{
		ProductID:  device.ProductID,
		DeviceName: device.DeviceName,
	}, clientToken))
	if val == "" || err != nil {
		return nil, err
	}
	var req reqType
	err = json.Unmarshal([]byte(val), &req)
	if err != nil {
		return nil, err
	}
	return &req, err
}
