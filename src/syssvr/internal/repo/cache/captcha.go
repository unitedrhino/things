package cache

import (
	"context"
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
func (c *Captcha) Verify(ctx context.Context, Type, codeID, code string) bool {
	key := c.GenKey(Type, codeID)
	val, err := c.store.GetCtx(ctx, key)
	if err != nil {
		return false
	}
	if val != "" { //如果验证码存在，则删除验证码
		c.store.DelCtx(ctx, key)
	}
	if val == code {
		return true
	}
	return false
}

func (c *Captcha) Store(ctx context.Context, Type, codeID, code string, expire int64) error {
	return c.store.SetexCtx(ctx, c.GenKey(Type, codeID), code, int(expire))
}
