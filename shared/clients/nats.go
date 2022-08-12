package clients

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/nats-io/nats.go"
	"sync"
)

var (
	natsClient   *nats.Conn
	natsInitOnce sync.Once
)

func NewNatsClient(conf conf.NatsConf) (nc *nats.Conn, err error) {
	natsInitOnce.Do(func() {
		connectOpts := nats.Options{
			Url:      conf.Url,
			User:     conf.User,
			Password: conf.Pass,
			Token:    conf.Token,
		}
		nc, err = connectOpts.Connect()
		if err != nil {
			return
		}
		natsClient = nc
	})
	return natsClient, err
}
