package verify

import (
	"fmt"
	"gitee.com/godLei6/things/shared/errors"
	"github.com/mojocn/base64Captcha"
	"github.com/tal-tech/go-zero/core/syncx"
	"time"
)

import "github.com/tal-tech/go-zero/core/stores/cache"

//type Store interface {
//	// Set sets the digits for the captcha id.
//	Set(id string, value string)
//
//	// Get returns stored digits for the captcha id. Clear indicates
//	// whether the captcha must be deleted from the store.
//	Get(id string, clear bool) string
//
//	//Verify captcha's answer directly
//	Verify(id, answer string, clear bool) bool
//}

var (
	// can't use one SharedCalls per conn, because multiple conns may share the same cache key.
	exclusiveCalls = syncx.NewSharedCalls()
	stats          = cache.NewStat("verify")
	cachePrefix    = "verify"
)

type Captcha struct {
	cache    cache.Cache   //redis的结构体
	keyPre   string        //redis中key的前缀
	keyExp   time.Duration //key的过期时间
	height   int           //验证码图片的长度
	width    int           //验证码图片的宽度
	length   int           //验证码的个数
	maxSkew  float64
	dotCount int
}

func NewCaptcha(height int, width int, length int, c cache.ClusterConf, keyExp time.Duration, opts ...cache.Option) *Captcha {
	return &Captcha{
		cache:    cache.New(c, exclusiveCalls, stats, errors.NotFind, opts...),
		keyExp:   keyExp,
		keyPre:   "captcha",
		height:   height,
		width:    width,
		length:   length,
		maxSkew:  0.7,
		dotCount: 80,
	}
}

func (c *Captcha) Verify(id, answer string) bool {
	var ans string
	key := fmt.Sprintf("%s#%s#%s", cachePrefix, c.keyPre, id)
	err := c.cache.Get(key, &ans)
	if err != nil {
		return false
	}
	c.cache.Del(key)
	return ans == answer
}

func (c *Captcha) Get() (string, string, error) {
	driver := base64Captcha.NewDriverDigit(c.height, c.width, c.length, c.maxSkew, c.dotCount)
	id, content, answer := driver.GenerateIdQuestionAnswer()
	item, _ := driver.DrawCaptcha(content)
	b64s := item.EncodeB64string()
	key := fmt.Sprintf("%s#%s#%s", cachePrefix, c.keyPre, id)
	err := c.cache.SetWithExpire(key, answer, c.keyExp)
	if err != nil {
		return "", "", err
	}
	return id, b64s, err
}
