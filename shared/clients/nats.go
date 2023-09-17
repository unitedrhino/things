package clients

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/utils"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

func NewNatsClient(conf conf.NatsConf) (nc *nats.Conn, err error) {
	connectOpts := nats.GetDefaultOptions()
	connectOpts.Url = conf.Url
	connectOpts.User = conf.User
	connectOpts.Password = conf.Pass
	connectOpts.Token = conf.Token
	connectOpts.DisconnectedErrCB = func(conn *nats.Conn, err error) {
		logx.Errorf("nats.DisconnectedErrCB  err:%v", err)
	}
	connectOpts.AsyncErrorCB = func(conn *nats.Conn, subscription *nats.Subscription, err error) {
		logx.Errorf("nats.AsyncErrorCB subscription:%v err:%v", utils.Fmt(subscription), err)
	}

	nc, err = connectOpts.Connect()
	if err != nil {
		return
	}
	return nc, err
}
