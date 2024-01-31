package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg"
	"github.com/zeromicro/go-zero/core/stores/kv"
)

const (
	msgExpr = 10 * 60
)

func genDeviceMsgKey(msgType string, handle string, Type string, device devices.Core, MsgToken string) string {
	return fmt.Sprintf("deviceMsg:%s:%s:%s:%s:%s:%s",
		handle, Type, msgType, device.ProductID, device.DeviceName, MsgToken)
}

func SetDeviceMsg(ctx context.Context, store kv.Store, msgType string, req *deviceMsg.PublishMsg, MsgToken string) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}
	err = store.SetexCtx(ctx, genDeviceMsgKey(msgType, req.Handle, req.Type, devices.Core{
		ProductID:  req.ProductID,
		DeviceName: req.DeviceName,
	}, MsgToken), string(payload), msgExpr)
	return err
}

func GetDeviceMsg[reqType any](ctx context.Context, store kv.Store, msgType string /*请求还是回复*/, handle string, Type string, device devices.Core, MsgToken string) (*reqType, error) {
	val, err := store.GetCtx(ctx, genDeviceMsgKey(msgType, handle, Type, devices.Core{
		ProductID:  device.ProductID,
		DeviceName: device.DeviceName,
	}, MsgToken))
	if val == "" || err != nil {
		return nil, err
	}
	var req deviceMsg.PublishMsg
	var ret reqType
	err = utils.Unmarshal([]byte(val), &req)
	if err != nil {
		return nil, err
	}
	err = utils.Unmarshal(req.Payload, &ret)
	return &ret, err
}
