package cache

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/stores/kv"
)

type Captcha struct {
	store kv.Store
}

func NewCaptcha(store kv.Store) *Captcha {
	return &Captcha{
		store: store,
	}
}
func (c *Captcha) GenKey(Type, codeID string) string {
	return "captcha:" + Type + ":" + codeID
}
func (c *Captcha) Verify(ctx context.Context, Type, codeID, code string) string {
	key := c.GenKey(Type, codeID)
	val, err := c.store.GetCtx(ctx, key)
	if err != nil || val == "" {
		return ""
	}
	//如果验证码存在，则删除验证码
	c.store.DelCtx(ctx, key)
	body := map[string]string{}
	json.Unmarshal([]byte(val), &body)
	if body["code"] == code {
		if body["account"] == "" {
			return " "
		}
		return body["account"]
	}
	return ""
}

func (c *Captcha) Store(ctx context.Context, Type, codeID, code string, account string, expire int64) error {
	body := map[string]interface{}{
		"code":    code,
		"account": account,
	}
	bodytStr, _ := json.Marshal(body)
	return c.store.SetexCtx(ctx, c.GenKey(Type, codeID), string(bodytStr), int(expire))
}
