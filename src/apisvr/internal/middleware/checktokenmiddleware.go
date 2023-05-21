package middleware

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gcharset"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/userHeader"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/config"
	auth "github.com/i-Things/things/src/syssvr/client/auth"
	operLog "github.com/i-Things/things/src/syssvr/client/log"
	user "github.com/i-Things/things/src/syssvr/client/user"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"strings"
)

type CheckTokenMiddleware struct {
	UserRpc user.User
	LogRpc  operLog.Log
	AuthRpc auth.Auth
	c       config.Config
}

func NewCheckTokenMiddleware(c config.Config, UserRpc user.User, AuthRpc auth.Auth, LogRpc operLog.Log) *CheckTokenMiddleware {
	return &CheckTokenMiddleware{UserRpc: UserRpc, c: c, AuthRpc: AuthRpc, LogRpc: LogRpc}
}

func (m *CheckTokenMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userCtx *userHeader.UserCtx

		isOpen, userCtx, err := m.OpenAuth(w, r)
		if isOpen { //如果是开放请求
			if err != nil {
				logx.WithContext(r.Context()).Errorf("%s.OpenAuth error=%s", utils.FuncName(), err)
				http.Error(w, "开放请求失败："+err.Error(), http.StatusUnauthorized)
				return
			}
		} else { //如果是用户请求
			//校验 Jwt Token
			userCtx, err = m.UserAuth(w, r)
			if err != nil {
				logx.WithContext(r.Context()).Errorf("%s.UserAuth error=%s", utils.FuncName(), err)
				http.Error(w, "用户请求失败："+err.Error(), http.StatusUnauthorized)
				return
			}

			//校验 Casbin Rule
			_, err = m.AuthRpc.AuthApiCheck(r.Context(), &user.CheckAuthReq{
				RoleID: userCtx.Role,
				Path:   r.URL.Path,
				Method: utils.MethodToNum(r.Method),
			})
			if err != nil {
				logx.WithContext(r.Context()).Errorf("%s.AuthApiCheck error=%s", utils.FuncName(), err)
				http.Error(w, "接口权限不足："+err.Error(), http.StatusUnauthorized)
				return
			}
		}

		//注入 用户信息 到 ctx
		ctx2 := userHeader.SetUserCtx(r.Context(), userCtx)
		r2 := r.WithContext(ctx2)

		//记录 接口响应日志
		err = m.OperationLogRecord("requestInfo", r2)
		if err != nil {
			logx.WithContext(r2.Context()).Errorf("%s.OperationLogRecord requestInfo error=%s", utils.FuncName(), err)
		}

		next(w, r2)

		//记录 接口响应日志
		err = m.OperationLogRecord("responseInfo", r2)
		if err != nil {
			logx.WithContext(r2.Context()).Errorf("%s.OperationLogRecord responseInfo error=%s", utils.FuncName(), err)
		}
	}
}

// 如果有开放认证的字段才进行认证
func (m *CheckTokenMiddleware) OpenAuth(w http.ResponseWriter, r *http.Request) (bool, *userHeader.UserCtx, error) {
	var isOpen bool
	userName, password, ok := r.BasicAuth()
	if !ok {
		return isOpen, nil, nil
	} else {
		isOpen = true
	}

	strIP, _ := utils.GetIP(r)
	if !m.c.OpenAuth.Auth(userName, password, strIP) {
		return isOpen, nil, errors.Permissions.AddMsg("开放认证没通过")
	}

	return isOpen, &userHeader.UserCtx{
		IsOpen:    isOpen,
		Uid:       0,
		Role:      0,
		IsAllData: true,
		IP:        strIP,
		Os:        r.Header.Get("User-Agent"),
	}, nil
}

func (m *CheckTokenMiddleware) UserAuth(w http.ResponseWriter, r *http.Request) (*userHeader.UserCtx, error) {
	strIP, _ := utils.GetIP(r)

	strToken := r.Header.Get(userHeader.UserToken)
	if strToken == "" {
		logx.WithContext(r.Context()).Errorf("%s.CheckToken ip=%s not find token",
			utils.FuncName(), strIP)
		return nil, errors.NotLogin
	}

	resp, err := m.UserRpc.UserCheckToken(r.Context(), &user.UserCheckTokenReq{
		Ip:    strIP,
		Token: strToken,
	})
	if err != nil {
		er := errors.Fmt(err)
		logx.WithContext(r.Context()).Errorf("%s.CheckToken ip=%s token=%s return=%s",
			utils.FuncName(), strIP, strToken, err)
		return nil, er
	}

	if resp.Token != "" {
		w.Header().Set("Access-Control-Expose-Headers", userHeader.UserSetToken)
		w.Header().Set(userHeader.UserSetToken, resp.Token)
	}
	logx.WithContext(r.Context()).Infof("%s.CheckToken ip:%v in.token=%s checkResp:%v",
		utils.FuncName(), strIP, strToken, utils.Fmt(resp))

	return &userHeader.UserCtx{
		IsOpen:    false,
		Uid:       resp.Uid,
		Role:      resp.Role,
		IsAllData: resp.IsAllData == def.True,
		IP:        strIP,
		Os:        r.Header.Get("User-Agent"),
	}, nil
}

// 获取ip所属城市
func (m *CheckTokenMiddleware) GetCityByIp(ip string) string {
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

// 操作日志记录
func (m *CheckTokenMiddleware) OperationLogRecord(logTitle string, r *http.Request) error {
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
		Msg:          fmt.Sprintf("logTitle:%s; statusMsg:%s", logTitle, respStatusMsg),
		Req:          reqBodyStr,
		Resp:         respBodyStr,
	})
	if err != nil {
		logx.WithContext(r.Context()).Errorf("%s.OperationLogRecord is error : %s",
			utils.FuncName(), err.Error())
	}

	return nil
}
