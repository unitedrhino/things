package edge

import (
	"bytes"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic/things/device/edge"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"io"
	"net/http"
)

// 设备使用http协议用云端交互,需要在http头中带上mqtt的账号密码(basic auth)
func SendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeviceInteractEdgeSendReq
		reqBody, _ := io.ReadAll(r.Body)                //读取 reqBody
		r.Body = io.NopCloser(bytes.NewReader(reqBody)) //重建 reqBody
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}
		l := edge.NewSendLogic(r.Context(), svcCtx)
		resp, err := l.Send(r, reqBody, &req)
		result.Http(w, r, resp, err)
	}
}
