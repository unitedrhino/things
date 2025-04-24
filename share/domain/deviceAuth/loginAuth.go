// Package device 设备登录权限认证及对应clientID,userName的处理
package deviceAuth

import (
	"encoding/base64"
	"fmt"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"
	"strings"
	"time"
)

const (
	ProductIDLen = 11
	HmacSha256   = "hmacsha256"
	HmacSha1     = "hmacsha1"
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
		ClientID       string
		ProductID      string //产品id
		DeviceName     string //设备名称
		SdkAppID       int64  //appid 直接填 12010126
		ConnID         string //随机6字节字符串 帮助查bug
		Expiry         int64  //过期时间 unix时间戳
		IsNeedRegister bool
	}
)

// 生成 MQTT 的 username 部分, 格式为 ${clientid};${sdkappid};${connid};${expiry}
func GetLoginDevice(userName string) (*LoginDevice, error) {
	keys := strings.Split(userName, ";")
	if len(keys) != 4 {
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

func GenClientID(productID string, deviceName string) string {
	return fmt.Sprintf("%s&%s", productID, deviceName)
}

func GetClientIDInfo(ClientID string) (*LoginDevice, error) {
	clientIDs := strings.Split(ClientID, "&")
	switch len(clientIDs) {
	case 2:
		lg := &LoginDevice{
			ProductID:  clientIDs[0],
			DeviceName: clientIDs[1],
			ClientID:   ClientID,
		}
		return lg, nil
	case 3:
		if clientIDs[0] != "register" {
			return nil, errors.Parameter.AddDetail("clientID not right")
		}
		lg := &LoginDevice{
			ProductID:      clientIDs[1],
			DeviceName:     clientIDs[2],
			IsNeedRegister: true,
			ClientID:       ClientID,
		}
		return lg, nil
	default:
		// 兼容老的clientID
		if len(ClientID) < ProductIDLen {
			return nil, errors.Parameter.AddDetail("clientID length not enough")
		}
		lg := &LoginDevice{
			ProductID:  ClientID[0:ProductIDLen],
			DeviceName: ClientID[ProductIDLen:],
			ClientID:   ClientID,
		}
		return lg, nil
	}
}

func GetStrProductID(id int64) string {
	return fmt.Sprintf("%X", id)
}

func GetInt64ProductID(id string) int64 {
	return utils.AnyToDecimal(id, 62)
}

func NewPwdInfo(signature string, signMethod string) (*PwdInfo, error) {
	var HmacHandle func(data string, secret []byte) string
	switch signMethod {
	case HmacSha256:
		HmacHandle = utils.HmacSha256
	case HmacSha1:
		HmacHandle = utils.HmacSha1
	default:
		return nil, errors.Parameter.AddDetail("password not support encrypt method:" + signMethod)
	}
	return &PwdInfo{
		token:      signature,
		hmac:       signMethod,
		HmacHandle: HmacHandle,
	}, nil
}

func NewPwdInfoWithPwd(password string) (*PwdInfo, error) {
	keys := strings.Split(password, ";")
	if len(keys) != 2 {
		return nil, errors.Parameter.AddDetail("password not right")
	}
	return NewPwdInfo(keys[0], keys[1])
}

/*
比较设备秘钥是否正确
@param secret 设备秘钥
*/
func (p *PwdInfo) CmpPwd(signature, secret string) error {
	pwd, _ := base64.StdEncoding.DecodeString(secret)
	passwrod := p.HmacHandle(signature, pwd)
	if passwrod != p.token {
		return errors.Password
	}
	return nil
}

func GenSecretDeviceInfo(hmacType string, productID string, deviceName string, deviceSecret string) (
	clientID, userName, password string) {
	var (
		connID = utils.Random(5, 1)
		expiry = time.Now().AddDate(0, 0, 10).Unix()
		token  string
		pwd, _ = base64.StdEncoding.DecodeString(deviceSecret)
	)
	clientID = productID + "&" + deviceName
	userName = fmt.Sprintf("%s;12010126;%s;%d", clientID, connID, expiry)
	if hmacType == HmacSha1 {
		token = utils.HmacSha1(userName, pwd)
		password = token + ";hmacsha1"
	} else {
		token = utils.HmacSha256(userName, pwd)
		password = token + ";hmacsha256"
	}
	return
}
