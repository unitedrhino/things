package middleware

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gcharset"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/i-Things/things/shared/domain/userHeader"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/config"
	operLog "github.com/i-Things/things/src/syssvr/client/log"
	"github.com/i-Things/things/src/syssvr/client/user"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"strings"
)

type TeardownWareMiddleware struct {
	cfg    config.Config
	LogRpc operLog.Log
}

func NewTeardownWareMiddleware(cfg config.Config, LogRpc operLog.Log) *TeardownWareMiddleware {
	return &TeardownWareMiddleware{cfg: cfg, LogRpc: LogRpc}
}

func (m *TeardownWareMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.WithContext(r.Context()).Infof("%s.Lifecycle.Before", utils.FuncName())

		//记录 接口响应日志
		err := m.OperationLogRecord(r)
		if err != nil {
			logx.WithContext(r.Context()).Errorf("%s.OperationLogRecord responseInfo error=%s", utils.FuncName(), err)
		}

		next(w, r)

		logx.WithContext(r.Context()).Infof("%s.Lifecycle.After", utils.FuncName())
	}
}

// 接口操作日志记录
func (m *TeardownWareMiddleware) OperationLogRecord(r *http.Request) error {
	useCtx := userHeader.GetUserCtx(r.Context())
	if useCtx.IsOpen || useCtx.Uid == 0 {
		return nil
	}

	reqBody, _ := io.ReadAll(r.Body)                //读取 reqBody
	r.Body = io.NopCloser(bytes.NewReader(reqBody)) //重建 reqBody
	reqBodyStr := string(reqBody)

	respStatusCode := http.StatusOK
	respStatusMsg := ""
	respBodyStr := ""

	if r.Response != nil {
		respStatusCode = r.Response.StatusCode
		respStatusMsg = r.Response.Status
		respBody, _ := io.ReadAll(r.Response.Body)                //读取 respBody
		r.Response.Body = io.NopCloser(bytes.NewReader(respBody)) //重建 respBody
		respBodyStr = string(respBody)
	}

	uri := "https://"
	if !strings.Contains(r.Proto, "HTTPS") {
		uri = "http://"
	}

	ipAddr, err := utils.GetIP(r)
	if err != nil {
		logx.WithContext(r.Context()).Errorf("%s.GetIP is error : %s req:%v",
			utils.FuncName(), err.Error(), utils.Fmt(r))
		ipAddr = "0.0.0.0"
	}

	_, err = m.LogRpc.OperLogCreate(r.Context(), &user.OperLogCreateReq{
		Uid:          userHeader.GetUserCtx(r.Context()).Uid,
		Uri:          uri + r.Host + r.RequestURI,
		Route:        r.RequestURI,
		OperIpAddr:   ipAddr,
		OperLocation: m.GetCityByIp(ipAddr),
		Code:         int64(respStatusCode),
		Msg:          respStatusMsg,
		Req:          reqBodyStr,
		Resp:         respBodyStr,
	})
	if err != nil {
		logx.WithContext(r.Context()).Errorf("%s.OperationLogRecord is error : %s",
			utils.FuncName(), err.Error())
	}

	return nil
}

// 获取ip所属城市
func (m *TeardownWareMiddleware) GetCityByIp(ip string) string {
	if ip == "" {
		return ""
	}
	if ip == "[::1]" || ip == "127.0.0.1" {
		return "内网IP"
	}

	url := "http://whois.pconline.com.cn/ipJson.jsp?json=true&ip=" + ip
	bytes := g.Client().GetBytes(context.TODO(), url)
	src := string(bytes)
	srcCharset := "GBK"

	tmp, _ := gcharset.ToUTF8(srcCharset, src)
	json, err := gjson.DecodeToJson(tmp)
	if err != nil {
		return ""
	}

	if json.Get("code").Int() == 0 {
		city := fmt.Sprintf("%s %s", json.Get("pro").String(), json.Get("city").String())
		return city
	} else {
		return ""
	}
}
