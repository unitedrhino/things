package device

import (
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic/things/protocol/script/device"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// 获取协议脚本列表设备
func IndexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProtocolScriptDeviceIndexReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := device.NewIndexLogic(r.Context(), svcCtx)
		resp, err := l.Index(&req)
		result.Http(w, r, resp, err)
	}
}
