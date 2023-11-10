package result

import (
	"bytes"
	"encoding/json"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"io/ioutil"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Http http返回
func Http(w http.ResponseWriter, r *http.Request, resp any, err error) {
	var code int
	var msg string
	if err == nil {
		//成功返回
		re := Success(resp)
		httpx.WriteJson(w, http.StatusOK, re)
		code = 200
		msg = "success"

	} else {
		//错误返回
		er := errors.Fmt(err)
		logx.WithContext(r.Context()).Errorf("【http handle err】router:%v err: %v ",
			r.URL.Path, utils.Fmt(er))
		httpx.WriteJson(w, http.StatusBadRequest, Error(er.Code, er.Msg))
		code = int(er.Code)
		msg = er.Msg
	}

	//将接口的应答结果写入r.Response，为操作日志记录接口提供应答信息
	bs, _ := json.Marshal(resp)
	var temp http.Response
	temp.StatusCode = code
	temp.Status = msg
	temp.Body = ioutil.NopCloser(bytes.NewReader(bs))
	r.Response = &temp

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

// indexapi http返回
func IndexApiHttp(w http.ResponseWriter, r *http.Request, resp any, err error) {
	var code int
	var msg string
	if err == nil {
		//成功返回
		re := Success(resp)
		httpx.WriteJson(w, http.StatusOK, re)
		code = 200
		msg = "success"

	} else {
		//错误返回
		er := errors.Fmt(err)
		logx.WithContext(r.Context()).Errorf("【http handle err】router:%v err: %v ",
			r.URL.Path, utils.Fmt(er))
		httpx.WriteJson(w, http.StatusBadRequest, Error(er.Code, er.Msg))
		code = int(er.Code)
		msg = er.Msg
	}

	//将接口的应答结果写入r.Response，为操作日志记录接口提供应答信息
	bs, _ := json.Marshal(resp)
	var temp http.Response
	temp.StatusCode = code
	temp.Status = msg
	temp.Body = ioutil.NopCloser(bytes.NewReader(bs))
	r.Response = &temp
}

// hook http返回
func HooksApiHttp(w http.ResponseWriter, r *http.Request, resp any, err error) {
	var code int
	var msg string
	if err == nil {
		//成功返回
		//re := Success(resp)
		httpx.WriteJson(w, http.StatusOK, resp)
		code = 200
		msg = "success"
	} else {
		//错误返回
		er := errors.Fmt(err)
		logx.WithContext(r.Context()).Errorf("【http handle err】router:%v err: %v ",
			r.URL.Path, utils.Fmt(er))
		httpx.WriteJson(w, http.StatusBadRequest, Error(er.Code, er.Msg))
		code = int(er.Code)
		msg = er.Msg
	}

	//将接口的应答结果写入r.Response，为操作日志记录接口提供应答信息
	bs, _ := json.Marshal(resp)
	var temp http.Response
	temp.StatusCode = code
	temp.Status = msg
	temp.Body = ioutil.NopCloser(bytes.NewReader(bs))
	r.Response = &temp
}
