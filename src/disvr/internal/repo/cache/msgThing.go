package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgThing"
	"github.com/zeromicro/go-zero/core/stores/kv"
)

type (
	MsgThingRepo struct {
		store kv.Store
	}
)

const (
	msgThingExpr = 10 * 60
)

func NewMsgThingRepo(store kv.Store) *MsgThingRepo {
	return &MsgThingRepo{
		store: store,
	}
}

func (m MsgThingRepo) genKey(msgType string /*请求还是回复*/, msgThingType string /*属性还是行为*/, device devices.Core, clientToken string) string {
	return fmt.Sprintf("msgThing_%s_%s_%s_%s_%s", msgType, msgThingType, device.ProductID, device.DeviceName, clientToken)
}

func (m MsgThingRepo) SetReq(ctx context.Context, msgThingType string, device devices.Core, req *msgThing.Req) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}
	err = m.store.SetexCtx(ctx, m.genKey(deviceMsg.ReqMsg, msgThingType, device, req.ClientToken), string(payload), msgThingExpr)
	return err
}
func (m MsgThingRepo) GetReq(ctx context.Context, msgThingType string, device devices.Core, clientToken string) (*msgThing.Req, error) {
	val, err := m.store.GetCtx(ctx, m.genKey(deviceMsg.ReqMsg, msgThingType, device, clientToken))
	if val == "" || err != nil {
		return nil, err
	}
	var req msgThing.Req
	err = json.Unmarshal([]byte(val), &req)
	if err != nil {
		return nil, err
	}
	return &req, err
}

func (m MsgThingRepo) SetResp(ctx context.Context, msgThingType string, device devices.Core, resp *msgThing.Resp) error {
	payload, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	err = m.store.SetexCtx(ctx, m.genKey(deviceMsg.RespMsg, msgThingType, device, resp.ClientToken), string(payload), msgThingExpr)
	return err
}
func (m MsgThingRepo) GetResp(ctx context.Context, msgThingType string, device devices.Core, clientToken string) (*msgThing.Resp, error) {
	val, err := m.store.GetCtx(ctx, m.genKey(deviceMsg.RespMsg, msgThingType, device, clientToken))
	if val == "" || err != nil {
		return nil, err
	}
	var resp msgThing.Resp
	err = json.Unmarshal([]byte(val), &resp)
	if err != nil {
		return nil, err
	}
	return &resp, err
}
