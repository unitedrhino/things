package info

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic/things/device/info"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// canBindErrorResponse 是 can-bind 业务错误时附带项目或共享线索的 HTTP 响应。
type canBindErrorResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// CanBindHandler 处理设备可绑定检查，并在已有访问关系时返回前端线索。
func CanBindHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeviceInfoCanBindReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := info.NewCanBindLogic(r.Context(), svcCtx)
		resp, err := l.CanBind(&req)
		if err != nil && resp != nil {
			writeCanBindErrorWithData(w, r, err, resp)
			return
		}
		result.Http(w, r, resp, err)
	}
}

// writeCanBindErrorWithData 保持错误码响应形态，并给 App 返回项目或共享线索。
func writeCanBindErrorWithData(w http.ResponseWriter, r *http.Request, err error, data any) {
	er := errors.Fmt(err)
	msg := er.GetI18nMsg(ctxs.GetUserCtxNoNil(r.Context()).AcceptLanguage)

	logx.WithContext(r.Context()).Errorf("【http handle err】router:%v err: %#v ",
		r.URL.Path, err)
	httpx.WriteJson(w, http.StatusOK, canBindErrorResponse{Code: er.Code, Msg: msg, Data: data})

	ret := ctxs.GetResp(r)
	if ret == nil {
		return
	}
	var temp http.Response
	temp.StatusCode = int(er.Code)
	temp.Status = msg
	bs, _ := json.Marshal(data)
	temp.Body = io.NopCloser(bytes.NewReader(bs))
	*ret = temp
}
