package logic

import "gitee.com/godLei6/things/src/dmsvr/device/msgquque/types"

type LogicHandle interface {
	Handle(msg *types.Elements) error
}
