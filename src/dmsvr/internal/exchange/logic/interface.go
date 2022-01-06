package logic

import "github.com/go-things/things/src/dmsvr/internal/exchange/types"

type LogicHandle interface {
	Handle(msg *types.Elements) error
}
