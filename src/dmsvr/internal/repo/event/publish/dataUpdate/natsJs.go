package dataUpdate

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/i-Things/things/shared/utils"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	NatsJsClient struct {
		client nats.JetStreamContext
	}
)

func newNatsJsClient(conf conf.NatsConf) (*NatsJsClient, error) {
	js, err := clients.NewNatsJetStreamClient(conf)
	if err != nil {
		return nil, err
	}
	return &NatsJsClient{client: js}, nil
}

func (n *NatsJsClient) ProductSchemaUpdate(ctx context.Context, info *events.DeviceUpdateInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	_, err = n.client.Publish(topics.DmProductSchemaUpdate, events.NewEventMsg(ctx, data))
	logx.WithContext(ctx).Infof("%s info:%v,err:%v", utils.FuncName(),
		utils.Fmt(info), err)
	return err
}

func (n *NatsJsClient) ProductCustomUpdate(ctx context.Context, info *events.DeviceUpdateInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	_, err = n.client.Publish(topics.DmProductCustomUpdate, events.NewEventMsg(ctx, data))
	logx.WithContext(ctx).Infof("%s info:%v,err:%v", utils.FuncName(),
		utils.Fmt(info), err)
	return err
}

func (n *NatsJsClient) DeviceLogLevelUpdate(ctx context.Context, info *events.DeviceUpdateInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	_, err = n.client.Publish(topics.DmDeviceLogLevelUpdate, events.NewEventMsg(ctx, data))
	logx.WithContext(ctx).Infof("%s info:%v,err:%v", utils.FuncName(),
		info, err)
	return err
}

func (n *NatsJsClient) DeviceGatewayUpdate(ctx context.Context, info *events.GatewayUpdateInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	_, err = n.client.Publish(topics.DmDeviceGatewayUpdate, events.NewEventMsg(ctx, data))
	logx.WithContext(ctx).Infof("%s info:%v,err:%v", utils.FuncName(),
		utils.Fmt(info), err)
	return err
}

func (n *NatsJsClient) DeviceRemoteConfigUpdate(ctx context.Context, info *events.DeviceUpdateInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	_, err = n.client.Publish(topics.DmDeviceRemoteConfigUpdate, events.NewEventMsg(ctx, data))
	logx.WithContext(ctx).Infof("%s info:%v,err:%v", utils.FuncName(),
		utils.Fmt(info), err)
	return err
}
