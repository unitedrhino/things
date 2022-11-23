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
	NatsClient struct {
		client *nats.Conn
	}
)

func newNatsClient(conf conf.NatsConf) (*NatsClient, error) {
	nc, err := clients.NewNatsClient(conf)
	if err != nil {
		return nil, err
	}
	return &NatsClient{client: nc}, nil
}

func (n *NatsClient) ProductSchemaUpdate(ctx context.Context, info *events.DataUpdateInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	err = n.client.Publish(topics.DmProductUpdateSchema, events.NewEventMsg(ctx, data))
	logx.WithContext(ctx).Infof("%s info:%v,err:%v", utils.FuncName(),
		utils.Fmt(info), err)
	return err
}

func (n *NatsClient) DeviceLogLevelUpdate(ctx context.Context, info *events.DataUpdateInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	err = n.client.Publish(topics.DmDeviceUpdateLogLevel, events.NewEventMsg(ctx, data))
	logx.WithContext(ctx).Infof("%s info:%v,err:%v", utils.FuncName(),
		info, err)
	return err
}

func (n *NatsClient) DeviceGatewayUpdate(ctx context.Context, info *events.GatewayUpdateInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	err = n.client.Publish(topics.DmDeviceUpdateGateway, events.NewEventMsg(ctx, data))
	logx.WithContext(ctx).Infof("%s info:%v,err:%v", utils.FuncName(),
		utils.Fmt(info), err)
	return err
}

func (n *NatsClient) DeviceRemoteConfigUpdate(ctx context.Context, info *events.DataUpdateInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	err = n.client.Publish(topics.DmDeviceUpdateRemoteConfig, events.NewEventMsg(ctx, data))
	logx.WithContext(ctx).Infof("%s info:%v,err:%v", utils.FuncName(),
		utils.Fmt(info), err)
	return err
}
