package logic

import "gitee.com/godLei6/things/src/dmsvr/internal/exchange/types"

type LogicHandle interface {
	Handle(msg *types.Elements) error
}
