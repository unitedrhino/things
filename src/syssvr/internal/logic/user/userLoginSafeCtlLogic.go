package userlogic

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strings"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginSafeCtlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type LoginSafeCtlInfo struct {
	prefix    string
	key       string
	timeout   int
	times     int
	forbidden int
}

func parseList(listStr []string, keyV0 string, keyV1 string) []*LoginSafeCtlInfo {
	var res []*LoginSafeCtlInfo
	for i, val := range listStr {
		var timeout int
		var times int
		var forbidden int
		tlist := strings.Split(val, ",")
		for _, tval := range tlist {
			t := strings.Split(tval, ":")
			switch t[0] {
			case "year":
				timeout = cast.ToInt(t[1]) * 365 * 24 * 3600
			case "month":
				timeout = cast.ToInt(t[1]) * 30 * 24 * 3600
			case "day":
				timeout = cast.ToInt(t[1]) * 24 * 3600
			case "times":
				times = cast.ToInt(t[1])
			case "forbidden":
				forbidden = cast.ToInt(t[1]) * 60
			}
		}

		info := &LoginSafeCtlInfo{
			prefix:    keyV0,
			key:       keyV0 + cast.ToString(i+1) + ":" + keyV1,
			timeout:   timeout,
			times:     times,
			forbidden: forbidden,
		}
		res = append(res, info)
	}

	return res
}

func parseWrongpassConf(counter conf.WrongPasswordCounter, userID string, ip string) []*LoginSafeCtlInfo {
	var res []*LoginSafeCtlInfo
	res = append(res, parseList(counter.Account, "login:wrongPassword:account:", userID)...)
	res = append(res, parseList(counter.Ip, "login:wrongPassword:ip:", ip)...)
	res = append(res, &LoginSafeCtlInfo{
		prefix:  "login:wrongPassword:captcha:",
		key:     "login:wrongPassword:captcha:" + userID,
		timeout: 24 * 3600,
		times:   counter.Captcha,
	})
	return res
}

func checkAccountForbidden(conn *redis.Redis, list []*LoginSafeCtlInfo) bool {
	for _, v := range list {
		if v.prefix == "login:wrongPassword:account:" {
			ret, err := conn.Get(v.key)
			if err != nil {
				continue
			}
			if cast.ToInt(ret) >= v.times {
				return true
			}
		}
	}
	return false
}

func checkIpForbidden(conn *redis.Redis, list []*LoginSafeCtlInfo) bool {

	for _, v := range list {
		if v.prefix == "login:wrongPassword:ip:" {
			ret, err := conn.Get(v.key)
			if err != nil {
				continue
			}
			if cast.ToInt(ret) >= v.times {
				return true
			}
		}
	}
	return false
}

func checkCaptchaTimes(conn *redis.Redis, list []*LoginSafeCtlInfo) (bool, error) {
	for _, v := range list {
		ret, err := conn.Get(v.key)
		if err != nil {
			if err == redis.Nil {
				err = conn.Setex(v.key, "1", v.timeout)
				if err != nil {
					return false, errors.Redis.AddMsgf("创建 redis key：%s 失败", v.key)
				}
			} else {
				return false, errors.Redis.AddMsgf("获取 redis key：%s 失败", v.key)
			}
		}
		if v.prefix != "login:wrongPassword:captcha:" {
			_, err = conn.Incr(v.key)
			if err != nil {
				return false, errors.Redis.AddMsgf("redis key：%s 自增失败", v.key)
			}
		} else {
			if cast.ToInt(ret) >= v.times {
				return true, nil
			}
		}
	}
	return false, nil
}

func NewUserLoginSafeCtlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginSafeCtlLogic {
	return &UserLoginSafeCtlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserLoginSafeCtlLogic) UserLoginSafeCtl(in *sys.UserLoginSafeCtlReq) (*sys.Response, error) {
	redis, err := redis.NewRedis(l.svcCtx.Config.CacheRedis[0].RedisConf)
	if err != nil {
		return nil, err
	}
	list := parseWrongpassConf(l.svcCtx.Config.WrongPasswordCounter, in.UserID, in.Ip)
	if !in.WrongPassword {
		if checkAccountForbidden(redis, list) {
			return nil, errors.AccountForbidden
		}
		if checkIpForbidden(redis, list) {
			return nil, errors.IpForbidden
		}
	} else {
		ret, err := checkCaptchaTimes(redis, list)
		if err != nil {
			return nil, err
		}
		if ret {
			return nil, errors.UseCaptcha
		}
	}

	return &sys.Response{}, nil
}
