package msg

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/result"
	"github.com/i-Things/things/src/apisvr/internal/logic/things/device/msg"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func SchemaLogIndexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeviceMsgSchemaLogIndexReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.AddMsgf("参数校验失败:%v", err.Error()))
			return
		}

		l := msg.NewSchemaLogIndexLogic(r.Context(), svcCtx)
		resp, err := l.SchemaLogIndex(&req)
		result.Http(w, r, resp, err)
	}
}
