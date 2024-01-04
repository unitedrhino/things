package websocket

import (
	"encoding/json"
	"net/http"
)

type response struct {
	req  *WsReq
	resp WsResp
}

func (r *response) Header() http.Header {
	ret := make(http.Header)
	for k, v := range r.resp.Handler {
		ret[k] = []string{v}
	}
	return ret
}

func (r *response) Write(buf []byte) (int, error) {
	var body map[string]any
	err := json.Unmarshal(buf, &body)
	if err != nil {
		r.resp.Body = body
		return len(buf), nil
	}
	r.resp.Body = body
	return len(buf), nil
}

func (r *response) WriteHeader(statusCode int) {
	r.resp.StatusCode = statusCode
}
