// Package device 设备登录权限认证及对应clientID,userName的处理
package device

import (
	"encoding/base64"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/spf13/cast"
	"strings"
)

const (
	PRODUCTID_LEN = 11
	Hmacsha256    = "hmacsha256"
	Hmacsha1      = "hmacsha1"
)

type (
	/*
		username 字段的格式为：
		${productId}${deviceName};${sdkappid};${connid};${expiry}
		注意：${} 表示变量，并非特定的拼接符号。

		password 字段格式为：
		${token};hmac 签名方法
		其中 hmac 签名方法字段填写第三步用到的摘要算法，可选的值有 hmacsha256 和 hmacsha1。
	*/
	PwdInfo struct {
		token      string //userName通过加密方法后的token
		hmac       string //签名的加密方法,共有两种:"hmacsha256","hmacsha1"
		HmacHandle func(data string, secret []byte) string
	}
	LoginDevice struct {
		ProductID  string //产品id
		DeviceName string //设备名称
		SdkAppID   int64  //appid 直接填 12010126
		ConnID     string //随机6字节字符串 帮助查bug
		Expiry     int64  //过期时间 unix时间戳
	}
)

func GetLoginDevice(userName string) (*LoginDevice, error) {
	keys := strings.Split(userName, ";")
	if len(keys) != 4 || len(keys[0]) < 11 {
		return nil, errors.Parameter.AddDetail("userName not right")
	}
	lg, err := GetClientIDInfo(keys[0])
	if err != nil {
		return nil, err
	}
	lg.SdkAppID = cast.ToInt64(keys[1])
	lg.ConnID = keys[2]
	lg.Expiry = cast.ToInt64(keys[3])
	return lg, nil
}

func GetClientIDInfo(ClientID string) (*LoginDevice, error) {
	if len(ClientID) < PRODUCTID_LEN {
		return nil, errors.Parameter.AddDetail("clientID length not enough")
	}
	lg := &LoginDevice{
		ProductID:  ClientID[0:PRODUCTID_LEN],
		DeviceName: ClientID[PRODUCTID_LEN:],
	}
	return lg, nil
}

//字符串类型的产品id有11个字节,不够的需要在前面补0
func GetStrProductID(id int64) string {
	str := utils.DecimalToAny(id, 62)
	return utils.ToLen(str, 11)
}

func GetInt64ProductID(id string) int64 {
	return utils.AnyToDecimal(id, 62)
}

func NewPwdInfo(password string) (*PwdInfo, error) {
	keys := strings.Split(password, ";")
	if len(keys) != 2 {
		return nil, errors.Parameter.AddDetail("password not right")
	}
	var HmacHandle func(data string, secret []byte) string
	switch keys[1] {
	case Hmacsha256:
		HmacHandle = utils.HmacSha256
	case Hmacsha1:
		HmacHandle = utils.HmacSha1
	default:
		return nil, errors.Parameter.AddDetail("password not suppot encrypt method:" + keys[1])
	}

	return &PwdInfo{
		token:      keys[0],
		hmac:       keys[1],
		HmacHandle: HmacHandle,
	}, nil
}

/*
比较设备秘钥是否正确
@param secret 设备秘钥
*/
func (p *PwdInfo) CmpPwd(userName, secret string) error {
	pwd, _ := base64.StdEncoding.DecodeString(secret)
	passwrod := p.HmacHandle(userName, pwd)
	if passwrod != p.token {
		return errors.Password
	}
	return nil
}
