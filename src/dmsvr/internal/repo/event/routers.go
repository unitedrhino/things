package event

import (
	deviceSend "github.com/go-things/things/src/dmsvr/internal/domain/deviceSend"
)

func (k *Kafka) AddRouters() {
	k.AddRouter(Router{
		Topic:   "onConnect",
		Handler: deviceSend.NewConnectLogic,
	})
	k.AddRouter(Router{
		Topic:   "onPublish",
		Handler: deviceSend.NewPublishLogic,
	})
	k.AddRouter(Router{
		Topic:   "onDisconnect",
		Handler: deviceSend.NewDisConnectLogic,
	})
}
