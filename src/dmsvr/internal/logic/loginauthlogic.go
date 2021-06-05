package logic

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/dmsvr/internal/svc"
	"gitee.com/godLei6/things/src/dmsvr/model"
	"github.com/spf13/cast"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
)

type LoginAuthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	di *model.DeviceInfo
}
var clientCert string = `-----BEGIN CERTIFICATE-----
MIIC3zCCAcegAwIBAgIBAjANBgkqhkiG9w0BAQsFADATMREwDwYDVQQDEwhNeVRl
c3RDQTAeFw0xNjEyMjYwMzA4MjNaFw0xNzEyMjYwMzA4MjNaMCIxDzANBgNVBAMM
BmNsb3VkMzEPMA0GA1UECgwGY2xpZW50MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEA40Q+bMUjxCOdDcdC2jZaX8HuNCdm6Mu1rgj8ZfyTJIzsKtv00LYd
xfdhlNFj1uq8wi/zK/cB95wBpG1Loo/WicqSP2G/A7aPnzIBPj3zzP7HdyM5EaHW
zDWLzK+f0+MmAsrp7UW/zBR5O+ScnmIWm2H7KJY36dJPKllzzw+R6a4eJ6vthBcm
nueIYrhdXnunaWzkWQqAWlSZCzD8/MfTkgAPYW7OoS6aAQugTBzhHRo1meOVIT7u
y+hmZE4kE8V98Iy1rGPV5Uz/1vSEJziJGvQkyVr3gcAv5DwLWnyX1vOXKIKOgRAz
VYB4tbfEP5tQ0jimG5ErftF/sGYFNTRslQIDAQABoy8wLTAJBgNVHRMEAjAAMAsG
A1UdDwQEAwIHgDATBgNVHSUEDDAKBggrBgEFBQcDAjANBgkqhkiG9w0BAQsFAAOC
AQEAMtVzZdj1y/TLxP7KZcqkd/Z/vdW6moo12tahDHR4vPq0NdGaHADRfZHbCBmb
JEI9Qz3CKSdKZRZ4/A3ui/ZltbvkCao9ilmhQXXDT3Yz5hxk5ZBC9+Zs1IZmrsis
Qg/cdLUx4+ei/eR0OgWyC5D9AKNzshQExKBGGojedb98VcuS5ccJKrq0kVzZ/BZQ
k1EswNC9ifKBcqPIO1rTD9T3PB7dv9ZRpxwslmgYWWsqQu9x/dnOPEHJ1yXr7KJh
47NP7OrtG21la8EcAtjA3GXDjiZ+tXIR1RMbx1iuJAQPBeddWJCPbyofdNRL59BQ
caLVQ77r0hpvnNpkafa5QptyfQ==
-----END CERTIFICATE-----`

var clientKey string = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA40Q+bMUjxCOdDcdC2jZaX8HuNCdm6Mu1rgj8ZfyTJIzsKtv0
0LYdxfdhlNFj1uq8wi/zK/cB95wBpG1Loo/WicqSP2G/A7aPnzIBPj3zzP7HdyM5
EaHWzDWLzK+f0+MmAsrp7UW/zBR5O+ScnmIWm2H7KJY36dJPKllzzw+R6a4eJ6vt
hBcmnueIYrhdXnunaWzkWQqAWlSZCzD8/MfTkgAPYW7OoS6aAQugTBzhHRo1meOV
IT7uy+hmZE4kE8V98Iy1rGPV5Uz/1vSEJziJGvQkyVr3gcAv5DwLWnyX1vOXKIKO
gRAzVYB4tbfEP5tQ0jimG5ErftF/sGYFNTRslQIDAQABAoIBAQCT4CzKM4AxOIcR
lw0t1V36nsIy10yDv0EI67nnVnAbwUJOJO7n+wfmby/kWFahWf3WUMLmYYO7LJx4
89DaBsOuxstgSGa0sM5E5JGggUkoosMBBz8z9N1B5LmBRuk1QsDR4lxR0ieZT90O
lpM+D07sbdWxtATPtNNkF+5d1aC4riPaNenwPXdb88bamcqCcARExwNxVhUogu88
frBeIfBdvNTZTmsqiqWrmAm4l1QnoQ1kCd3br4vbOlI4aQZCAPhECMaBSM7soNax
6XHUAA35vB3njNgvQYb6X2HvfktenwKXxDKDm7T8E6Ckof0kySncu2tpcIU/aHi4
QxS2TenhAoGBAPeU48RIbKYt158xmBbiY6EzMRHI1mq+iItiYGwjt43td4l1nEX+
UVGNnRJDffPPWIwNabPnOw9ZClwyEWgkJNJ/OS542B5QtFA5don5uAiX5OZCtQ6/
jyedC2HLq+e4No00pBkko3sVKbUHD98qRd45PFwhC34HJGjzxj3C16ZdAoGBAOr+
hN/2JSDOW+0dpbwVJUT1u3Ir9nOFZ3N5LDgvkE7dosKaHY9+AtUjMvhJ8Vea3jbJ
3VZrmacVtOPrrbsVWeacibqdDvRkPbQeg8vJjymLAFuvpv9WI6rih7PoiiG0HfSR
8aS14QTId31+6d9vgH/oWQuqNcTnqG3xWkK8HAuZAoGBAKE0INm9DoFld+//qre7
0IM1gc3Cp1n5lY6sD3xaBTo0VJD8MzSf0vL28j7iEzCc4VrPoPOyq5HiuAwvzYWx
gwhMLj9ED/QtODrEL5rHLjzqKfCDnsBrmhqA9thGdTf7igmHLRHx+UA7F1z3rC3y
qGt5eQPDwGfe3qY3k+zC4QdBAoGBANvaF3J5FS9mITbr39zhY6bqx93/J2nYy3qL
SUWfqkE+tkGecj2HRRsm/U6xzyuI5pEXtw5dSLm7YytBmZ5IUX2hwnFm81DOX7Qe
QGvuPRQ+yaz93x1P97quiQtWabUykDv6NrtEtisFalVs4V17Mht4w6ZYLknz+e4y
OaHp38sxAoGAF2ZBRadUjrfYN+BKJxekdvLEzGlRICvBRdB6vDfJPNULp1cVIzbF
rNhpjJJb53OvSJwI6OwRt5ehfIg1sRjoSYXhE6yJyEBQRIRdPLbxSAaQB20P9ZlL
blA0kLm6HiGNSu1CTAst23i2WueGQgOHHdBQoLUU5xEBNFYB2S7OB74=
-----END RSA PRIVATE KEY-----`
var clientCertificate tls.Certificate
var x509Cert *x509.Certificate
func NewLoginAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginAuthLogic {
	clientCertificate, _ = tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
	x509Cert, _ = x509.ParseCertificate(clientCertificate.Certificate[0])
	return &LoginAuthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}


/*
username 字段的格式为：
${productId}${deviceName};${sdkappid};${connid};${expiry}
注意：${} 表示变量，并非特定的拼接符号。

*/
type LoginDevice struct {
	ClientID 	string	//clientID
	ProductID 	int64	//产品id
	DeviceName 	string	//设备名称
	SdkAppID   	int64	//appid 直接填 12010126
	ConnID		string	//随机6字节字符串 帮助查bug
	Expiry 		int64	//过期时间 unix时间戳
}

/*
password 字段格式为：
${token};hmac 签名方法
其中 hmac 签名方法字段填写第三步用到的摘要算法，可选的值有 hmacsha256 和 hmacsha1。
*/
type PwdInfo struct {
	token string	//userName通过加密方法后的token
	hmac  string    //签名的加密方法,共有两种:"hmacsha256","hmacsha1"
	HmacHandle  func(data string, secret []byte) string
}
const (
	Hmacsha256 = "hmacsha256"
	Hmacsha1   = "hmacsha1"
)



func (l *LoginAuthLogic)GetPwdInfo(password string) (*PwdInfo,error){
	keys :=strings.Split(password,";")
	if len(keys) != 2{
		return nil,errors.Parameter.AddDetail("password not right")
	}
	var HmacHandle  func(data string, secret []byte) string
	switch keys[1] {
	case Hmacsha256:
		HmacHandle = utils.HmacSha256
	case Hmacsha1:
		HmacHandle = utils.HmacSha1
	default:
		return nil,errors.Parameter.AddDetail("password not suppot encrypt method:"+keys[1])
	}

	return &PwdInfo{
		token: keys[0],
		hmac: keys[1],
		HmacHandle:HmacHandle,
	}, nil
}

func (l *LoginAuthLogic)GetLoginDevice(userName string) (*LoginDevice,error){
	keys :=strings.Split(userName,";")
	if len(keys) != 4 || len(keys[0]) < 11{
		return nil,errors.Parameter.AddDetail("userName not right")
	}
	ProductID := dm.GetInt64ProductID(keys[0][0:11])
	if ProductID < 0 {
		return nil,errors.Parameter.AddDetail("product id not right")
	}
	DeviceName := keys[0][11:]
	lg:= &LoginDevice{
		ClientID	: keys[0],
		ProductID	: ProductID,
		DeviceName	: DeviceName,
		SdkAppID	: cast.ToInt64(keys[1]),
		ConnID		: keys[2],
		Expiry		: cast.ToInt64(keys[3]),
	}
	l.Slowf("LoginDevice=%+v",lg)
	return lg,nil
}

func (l *LoginAuthLogic)CmpPwd(in *dm.LoginAuthReq) error{
	if l.di == nil {
		panic("neet select  device info db first")
	}
	pwdInfo, err:= l.GetPwdInfo(in.Password)
	if err != nil {
		return err
	}
	pwd,_ := base64.StdEncoding.DecodeString(l.di.Secret)
	passwrod := pwdInfo.HmacHandle(in.Username,pwd)
	if passwrod != pwdInfo.token{
		return errors.Password
	}
	return nil
}

func (l *LoginAuthLogic)UpdateLoginTime(){
	if l.di == nil {
		panic("neet select  device info db first")
	}
	now := sql.NullTime{
		Valid: true,
		Time: time.Now(),
	}
	if l.di.FirstLogin.Valid == false {
		l.di.FirstLogin = now
	}
	l.di.UpdatedTime = now
	l.di.LastLogin = now
	l.svcCtx.DeviceInfo.Update(*l.di)
}


func (l *LoginAuthLogic) LoginAuth(in *dm.LoginAuthReq) (*dm.Response, error) {
	l.Infof("LoginAuth|req=%+v",in)
	if len(in.Certificate) > 0 {
		if bytes.Equal(in.Certificate,x509Cert.Signature){
			l.Error("it is same")
		}
		l.Errorf("cert len=%d|signature len=%d",
			len(x509Cert.Raw),len(x509Cert.Signature))
	}
	//生成 MQTT 的 username 部分, 格式为 ${clientid};${sdkappid};${connid};${expiry}
	lg, err :=l.GetLoginDevice(in.Username)
	if err != nil {
		return nil, err
	}
	if lg.ClientID != in.Clientid {
		return nil, errors.Parameter.AddDetail("userName'clientID not equal real client id")
	}
	if lg.Expiry < time.Now().Unix(){
		return nil, errors.SignatureExpired
	}
	l.di,err = l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(lg.ProductID,lg.DeviceName)
	if err!= nil {
		if err == model.ErrNotFound{
			return nil, errors.Password
		}else {
			l.Errorf("LoginAuth|FindOneByProductIDDeviceName failure|err=%+v",err)
			return nil, errors.Database
		}
	}
	err = l.CmpPwd(in)
	if err != nil {
		return nil, err
	}
	l.UpdateLoginTime()
	return &dm.Response{}, nil
}
