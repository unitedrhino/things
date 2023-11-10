package deviceinteractlogic

import (
	"github.com/i-Things/things/src/disvr/internal/domain/serverDo"
	"github.com/i-Things/things/src/disvr/pb/di"
)

func ToSendOptionDo(in *di.SendOption) *serverDo.SendOption {
	if in == nil {
		return nil
	}
	return &serverDo.SendOption{
		TimeoutToFail:  in.TimeoutToFail,
		RequestTimeout: in.RequestTimeout,
		RetryInterval:  in.RetryInterval,
	}
}
