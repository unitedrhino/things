package nats

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"testing"
	"time"
)

var (
	natsConf conf.NatsConf
)

func TestGetApplicationMsg(t *testing.T) {
	nc, err := clients.NewNatsClient(natsConf)
	if err != nil {
		t.Error(err)
	}
	//订阅所有消息
	_, err = nc.Subscribe(">", events.NatsSubscription(func(ctx context.Context, msg []byte) error {
		ele, err := deviceMsg.GetDevPublish(ctx, msg)
		if err != nil {
			return err
		}
		fmt.Println(ele)
		return err
	}))
	if err != nil {
		t.Error(err)
	}
	for {
		time.Sleep(time.Hour)
	}
}
