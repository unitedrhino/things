package result

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Http http返回
func Http(w http.ResponseWriter, r *http.Request, resp any, err error) {
	if err == nil {
		//成功返回
		r := Success(resp)
		httpx.WriteJson(w, http.StatusOK, r)
	} else {
		//错误返回
		er := errors.Fmt(err)
		logx.WithContext(r.Context()).Errorf("【http handle err】router:%v err: %v ",
			r.URL.Path, utils.Fmt(er))
		httpx.WriteJson(w, http.StatusBadRequest, Error(er.Code, er.Msg))
	}
}

// HttpWithoutWrap http返回，无包装
func HttpWithoutWrap(w http.ResponseWriter, r *http.Request, resp any, err error) {
	if err == nil {
		//成功返回
		httpx.WriteJson(w, http.StatusOK, resp)
	} else {
		//错误返回
		er := errors.Fmt(err)
		logx.WithContext(r.Context()).Errorf("【http handle err】router:%v err: %v ",
			r.URL.Path, utils.Fmt(er))
		httpx.WriteJson(w, http.StatusBadRequest, Error(er.Code, er.Msg))
	}
}
