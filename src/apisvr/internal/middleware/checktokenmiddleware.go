package middleware

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gcharset"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/config"
	"github.com/i-Things/things/src/apisvr/internal/domain/userHeader"
	operLog "github.com/i-Things/things/src/syssvr/client/log"
	user "github.com/i-Things/things/src/syssvr/client/user"
	"github.com/zeromicro/go-zero/core/logx"
	"io/ioutil"
	"net/http"
	"strings"
)

type CheckTokenMiddleware struct {
	UserRpc user.User
	LogRpc  operLog.Log
	c       config.Config
}

func NewCheckTokenMiddleware(c config.Config, UserRpc user.User) *CheckTokenMiddleware {
	return &CheckTokenMiddleware{UserRpc: UserRpc, c: c}
}

func (m *CheckTokenMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//err, isOpen := m.OpenAuth(w, r)
		//if isOpen { //如果是开放请求
		//	if err == nil {
		//		next(w, r)
		//		//if r.Response != nil {
		//		//	m.OperationLogRecord(r)
		//		//}
		//	} else {
		//		http.Error(w, err.Error(), http.StatusUnauthorized)
		//	}
		//	return
		//}
		userCtx, err := m.UserAuth(w, r)
		if err == nil {
			next(w, r.WithContext(userHeader.SetUserCtx(r.Context(), userCtx)))
			if r.Response != nil {
				m.OperationLogRecord(r)
			}
			return
		}

		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
}

// 如果有开放认证的字段才进行认证
func (m *CheckTokenMiddleware) OpenAuth(w http.ResponseWriter, r *http.Request) (error, bool) {
	userName, password, ok := r.BasicAuth()
	if !ok {
		return nil, false
	}
	strIP, _ := utils.GetIP(r)
	if !m.c.OpenAuth.Auth(userName, password, strIP) {
		return errors.Permissions.AddMsg("开放认证没通过"), true
	}
	return nil, true
}

func (m *CheckTokenMiddleware) UserAuth(w http.ResponseWriter, r *http.Request) (*userHeader.UserCtx, error) {
	strIP, _ := utils.GetIP(r)
	strToken := r.Header.Get(userHeader.UserToken)
	if strToken == "" {
		logx.WithContext(r.Context()).Errorf("%s.CheckToken ip=%s not find token",
			utils.FuncName(), strIP)
		return nil, errors.NotLogin
	}
	resp, err := m.UserRpc.CheckToken(r.Context(), &user.CheckTokenReq{
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
		Uid:  resp.Uid,
		IP:   strIP,
		Role: resp.Role,
		Os:   r.Header.Get("User-Agent"),
	}, nil
}

//获取ip所属城市
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

//操作日志记录
func (m *CheckTokenMiddleware) OperationLogRecord(r *http.Request) error {
	re, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	res, err := ioutil.ReadAll(r.Response.Body)
	if err != nil {
		return err
	}

	var p []byte
	r.Body.Read(p)
	fmt.Println(p, p)
	ipAddr := r.Host[0:strings.Index(r.Host, ":")]
	_, err = m.LogRpc.OperLogCreate(r.Context(), &user.OperLogCreateReq{
		//Uid:          userHeader.GetUserCtx(r.Context()).Uid,
		Uri:          r.RequestURI,
		OperIpAddr:   ipAddr,
		OperLocation: m.GetCityByIp(ipAddr),
		Req:          string(re),
		Resp:         string(res),
		Code:         int64(r.Response.StatusCode),
		Msg:          r.Response.Status,
	})
	if err != nil {
		return err
	}
	return nil
}
