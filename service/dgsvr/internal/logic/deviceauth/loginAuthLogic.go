package deviceauthlogic

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/deviceAuth"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/dgsvr/internal/svc"
	"gitee.com/i-Things/things/service/dgsvr/pb/dg"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginAuthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
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

func init() {
	clientCertificate, _ = tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
	x509Cert, _ = x509.ParseCertificate(clientCertificate.Certificate[0])
}
func NewLoginAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginAuthLogic {
	return &LoginAuthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 设备登录认证
func (l *LoginAuthLogic) LoginAuth(in *dg.LoginAuthReq) (*dg.Response, error) {
	l.Infof("%s req=%+v", utils.FuncName(), in)
	if len(in.Certificate) > 0 {
		if bytes.Equal(in.Certificate, x509Cert.Signature) {
			l.Error("it is same")
		}
		l.Errorf("cert len=%d signature len=%d",
			len(x509Cert.Raw), len(x509Cert.Signature))
	}
	inLg, err := deviceAuth.GetClientIDInfo(in.ClientID)
	if err != nil {
		return nil, err
	}
	if inLg.IsNeedRegister { //如果是需要注册的直接放过
		return &dg.Response{}, nil
	}
	//生成 MQTT 的 username 部分, 格式为 ${clientid};${sdkappid};${connid};${expiry}
	lg, err := deviceAuth.GetLoginDevice(in.Username)
	if err != nil {
		return nil, err
	}

	if lg.ProductID != inLg.ProductID || lg.DeviceName != inLg.DeviceName {
		return nil, errors.Parameter.AddDetail("userName'clientID not equal real client id")
	}
	//if lg.Expiry < time.Now().Unix() {
	//	return nil, errors.SignatureExpired
	//}

	di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
		ProductID:  lg.ProductID,
		DeviceName: lg.DeviceName,
	})
	if err != nil {
		return nil, err
	}
	pwd, err := deviceAuth.NewPwdInfoWithPwd(in.Password)
	if err != nil {
		return nil, err
	}
	err = pwd.CmpPwd(in.Username, di.Secret)
	if err != nil {
		return nil, err
	}
	return &dg.Response{}, nil
}
