package pubInner

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/i-Things/share/clients"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/events/topics"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/sdk/service/protocol"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	NatsJsClient struct {
		client       *clients.NatsClient
		protocolCode string
	}
)

func newNatsClient(conf conf.EventConf, protocolCode string, nodeID int64) (PubInner, error) {
	nc, err := clients.NewNatsClient2(conf.Mode, conf.Nats.Consumer, conf.Nats, nodeID)
	if err != nil {
		return nil, err
	}
	return &NatsJsClient{client: nc, protocolCode: protocolCode}, nil
}

func (n *NatsJsClient) DevPubGateway(ctx context.Context, publishMsg *devices.DevPublish) error {
	publishMsg.ProtocolCode = n.protocolCode
	pubStr, _ := json.Marshal(publishMsg)
	return n.publish(ctx,
		fmt.Sprintf(topics.DeviceUpMsg, publishMsg.Handle, publishMsg.ProductID, publishMsg.DeviceName), pubStr)
}

func (n *NatsJsClient) DevPubMsg(ctx context.Context, publishMsg *devices.DevPublish) error {
	protocol.UpdateDeviceActivity(ctx, devices.Core{
		ProductID:  publishMsg.ProductID,
		DeviceName: publishMsg.DeviceName,
	})
	publishMsg.ProtocolCode = n.protocolCode
	pubStr, _ := json.Marshal(publishMsg)
	err := n.publish(ctx,
		fmt.Sprintf(topics.DeviceUpMsg, publishMsg.Handle, publishMsg.ProductID, publishMsg.DeviceName), pubStr)
	if err != nil {
		logx.Errorf("%s.publish  err:%v", utils.FuncName(), err)
		return err
	}
	return err
}

func (n *NatsJsClient) DevPubThing(ctx context.Context, publishMsg *devices.DevPublish) error {
	publishMsg.ProtocolCode = n.protocolCode
	pubStr, _ := json.Marshal(publishMsg)
	err := n.publish(ctx,
		fmt.Sprintf(topics.DeviceUpThing, publishMsg.ProductID, publishMsg.DeviceName), pubStr)
	if err != nil {
		logx.Errorf("%s.publish  err:%v", utils.FuncName(), err)
		return err
	}
	return err
}

func (n *NatsJsClient) DevPubOta(ctx context.Context, publishMsg *devices.DevPublish) error {
	publishMsg.ProtocolCode = n.protocolCode
	pubStr, _ := json.Marshal(publishMsg)
	err := n.publish(ctx,
		fmt.Sprintf(topics.DeviceUpOta, publishMsg.ProductID, publishMsg.DeviceName), pubStr)
	if err != nil {
		logx.Errorf("%s.publish  err:%v", utils.FuncName(), err)
		return err
	}
	return err
}

func (n *NatsJsClient) DevPubConfig(ctx context.Context, publishMsg *devices.DevPublish) error {
	publishMsg.ProtocolCode = n.protocolCode
	pubStr, _ := json.Marshal(publishMsg)
	err := n.publish(ctx,
		fmt.Sprintf(topics.DeviceUpConfig, publishMsg.ProductID, publishMsg.DeviceName), pubStr)
	if err != nil {
		logx.Errorf("%s.publish  err:%v", utils.FuncName(), err)
		return err
	}
	return err
}

func (n *NatsJsClient) DevPubShadow(ctx context.Context, publishMsg *devices.DevPublish) error {
	publishMsg.ProtocolCode = n.protocolCode
	pubStr, _ := json.Marshal(publishMsg)
	err := n.publish(ctx,
		fmt.Sprintf(topics.DeviceUpShadow, publishMsg.ProductID, publishMsg.DeviceName), pubStr)
	if err != nil {
		logx.Errorf("%s.publish  err:%v", utils.FuncName(), err)
		return err
	}
	return err
}

func (n *NatsJsClient) DevPubSDKLog(ctx context.Context, publishMsg *devices.DevPublish) error {
	publishMsg.ProtocolCode = n.protocolCode
	pubStr, _ := json.Marshal(publishMsg)
	err := n.publish(ctx,
		fmt.Sprintf(topics.DeviceUpSDKLog, publishMsg.ProductID, publishMsg.DeviceName), pubStr)
	if err != nil {
		logx.Errorf("%s.publish  err:%v", utils.FuncName(), err)
		return err
	}
	return err
}

func (n *NatsJsClient) PubConn(ctx context.Context, conn ConnType, info *devices.DevConn) error {
	str, _ := json.Marshal(info)
	var err error
	switch conn {
	case Connect:
		err = n.publish(ctx, topics.DeviceUpStatusConnected, str)
	case DisConnect:
		err = n.publish(ctx, topics.DeviceUpStatusDisconnected, str)
	default:
		panic("not support conn type")
	}
	if err != nil {
		logx.Errorf("%s.publish  err:%v", utils.FuncName(), err)
		return err
	}
	return err
}

func (n *NatsJsClient) publish(ctx context.Context, topic string, payload []byte) error {
	err := n.client.Publish(ctx, topic, payload)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s nats publish failure err:%v topic:%v", err, topic)
	}
	return err
}
