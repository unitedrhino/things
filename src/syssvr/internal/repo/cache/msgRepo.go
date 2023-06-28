package cache

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/kv"
)

func CheckAccountOrIpForbidden(ctx context.Context, store kv.Store, list []*conf.LoginSafeCtlInfo) (int32, bool) {
	for _, v := range list {
		if v.Prefix != "login:wrongPassword:captcha:" {
			ret, err := store.GetCtx(ctx, v.Key)
			if err != nil {
				continue
			}
			if cast.ToInt(ret) >= v.Times {
				return int32(v.Forbidden), true
			}
		}
	}
	return 0, false
}

func CheckCaptchaTimes(ctx context.Context, store kv.Store, list []*conf.LoginSafeCtlInfo) (bool, error) {
	for _, v := range list {
		ret, err := store.GetCtx(ctx, v.Key)
		if ret == "" {
			err = store.SetexCtx(ctx, v.Key, "1", v.Timeout)
			if err != nil {
				return false, errors.Database.AddMsgf("创建 redis key：%s 失败", v.Key)
			}
			continue
		}

		_, err = store.Incr(v.Key)
		if err != nil {
			return false, errors.Database.AddMsgf("redis key：%s 自增失败", v.Key)
		}
		if v.Prefix != "login:wrongPassword:captcha:" {
			if cast.ToInt(ret)+1 >= v.Times {
				err = store.SetexCtx(ctx, v.Key, cast.ToString(cast.ToInt(ret)+1), v.Forbidden)
				if err != nil {
					return false, errors.Database.AddMsgf("重置 key：%s 时间失败", v.Key)
				}
			}
		} else {
			if cast.ToInt(ret)+1 >= v.Times {
				return true, nil
			}
		}
	}
	return false, nil
}

func ClearWrongpassKeys(ctx context.Context, store kv.Store, list []*conf.LoginSafeCtlInfo) {
	for _, v := range list {
		if v.Prefix != "login:wrongPassword:ip:" {
			store.DelCtx(ctx, v.Key)
		}
	}
}
