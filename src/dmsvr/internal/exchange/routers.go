package exchange

import "github.com/go-things/things/src/dmsvr/internal/exchange/logic"

func (k *Kafka) AddRouters() {
	k.AddRouter(Router{
		Topic:   "onConnect",
		Handler: logic.NewConnectLogic,
	})
	k.AddRouter(Router{
		Topic:   "onPublish",
		Handler: logic.NewPublishLogic,
	})
	k.AddRouter(Router{
		Topic:   "onDisconnect",
		Handler: logic.NewDisConnectLogic,
	})
}
