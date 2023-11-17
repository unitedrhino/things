package deviceinteractlogic

import (
	"github.com/i-Things/things/src/dmsvr/internal/domain/serverDo"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func ToSendOptionDo(in *dm.SendOption) *serverDo.SendOption {
	if in == nil {
		return nil
	}
	return &serverDo.SendOption{
		TimeoutToFail:  in.TimeoutToFail,
		RequestTimeout: in.RequestTimeout,
		RetryInterval:  in.RetryInterval,
	}
}
