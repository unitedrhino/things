package websocket

import (
	"github.com/zeromicro/go-zero/rest"
	"time"
)

type WsType string

const (
	Sub        WsType = "up.sub"          //订阅
	SubRet     WsType = "down.subRet"     //订阅回复
	Pub        WsType = "down.pub"        //发布
	Control    WsType = "up.control"      //控制
	ControlRet WsType = "down.controlRet" //控制回复
	UnSub      WsType = "up.unSub"        //取消订阅
	UnSubRet   WsType = "down.unSubRet"   //取消订阅回复
)

type (
	WsBody struct {
		Handler map[string]string `json:"handler,omitempty"`
		Type    WsType            `json:"type,omitempty"` //req 请求类型
		Path    string            `json:"path,omitempty"` //url路径或发布及订阅的主题
		Body    any               `json:"body,omitempty"` //消息体
	}
	WsReq struct {
		// Method specifies the HTTP method (GET, POST, PUT, etc.).
		// For client requests, an empty string means GET.
		//
		// Go's HTTP client does not support sending a request with
		// the CONNECT method. See the documentation on Transport for
		// details.
		Method string `json:"method"`
		WsBody
	}

	WsResp struct {
		StatusCode int `json:"statusCode"` // http状态码 200
		WsBody
	}

	Ping struct {
		Ping int64 `json:"ping"` //时间戳
	}

	Pong struct {
		Pong int64 `json:"pong"` //时间戳
	}
)

type (
	jwtSetting struct {
		enabled    bool
		secret     string
		prevSecret string
	}
	RouteOption      func(r *featuredRoutes)
	signatureSetting struct {
		rest.SignatureConf
		enabled bool
	}

	featuredRoutes struct {
		timeout   time.Duration
		priority  bool
		jwt       jwtSetting
		signature signatureSetting
		routes    []rest.Route
		maxBytes  int64
	}
)
