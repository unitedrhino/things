package script

import (
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic/things/protocol/script"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// 获取协议脚本列表
func IndexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProtocolScriptIndexReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := script.NewIndexLogic(r.Context(), svcCtx)
		resp, err := l.Index(&req)
		result.Http(w, r, resp, err)
	}
}
