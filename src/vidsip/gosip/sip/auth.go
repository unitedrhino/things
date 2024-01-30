package sip

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

// Authorization currently only Digest and MD5
type Authorization struct {
	realm     string
	nonce     string
	algorithm string
	username  string
	password  string
	uri       string
	response  string
	method    string
	qop       string
	nc        string
	cnonce    string
	other     map[string]string
	Data      map[string]string
}

// AuthFromValue AuthFromValue
func AuthFromValue(value string) *Authorization {
	auth := &Authorization{
		algorithm: "MD5",
		other:     make(map[string]string),
		Data:      make(map[string]string),
	}

	re := regexp.MustCompile(`([\w]+)="?([^",]+)"?`)
	matches := re.FindAllStringSubmatch(value, -1)
	for _, match := range matches {
		switch match[1] {
		case "realm":
			auth.realm = match[2]
		case "algorithm":
			auth.algorithm = match[2]
		case "nonce":
			auth.nonce = match[2]
		case "username":
			auth.username = match[2]
		case "uri":
			auth.uri = match[2]
		case "response":
			auth.response = match[2]
		case "qop":
			for _, v := range strings.Split(match[2], ",") {
				v = strings.Trim(v, " ")
				if v == "auth" || v == "auth-int" {
					auth.qop = "auth"
					break
				}
			}
		case "nc":
			auth.nc = match[2]
		case "cnonce":
			auth.cnonce = match[2]
		default:
			auth.other[match[1]] = match[2]
		}
		auth.Data[match[1]] = match[2]
	}
	return auth
}

// Get Get
func (auth *Authorization) Get(key string) string {
	return auth.Data[key]
}

// SetUsername SetUsername
func (auth *Authorization) SetUsername(username string) *Authorization {
	auth.username = username

	return auth
}

// SetURI SetURI
func (auth *Authorization) SetURI(uri string) *Authorization {
	auth.uri = uri

	return auth
}

// SetMethod SetMethod
func (auth *Authorization) SetMethod(method string) *Authorization {
	auth.method = method

	return auth
}

// SetPassword SetPassword
func (auth *Authorization) SetPassword(password string) *Authorization {
	auth.password = password

	return auth
}

// CalcResponse CalcResponse
func (auth *Authorization) CalcResponse() string {
	auth.response = CalcResponse(
		auth.username,
		auth.realm,
		auth.password,
		auth.method,
		auth.uri,
		auth.nonce,
		auth.qop,
		auth.cnonce,
		auth.nc,
	)

	return auth.response
}

func (auth *Authorization) String() string {
	if auth == nil {
		return "<nil>"
	}

	str := fmt.Sprintf(
		`Digest realm="%s",algorithm=%s,nonce="%s",username="%s",uri="%s",response="%s"`,
		auth.realm,
		auth.algorithm,
		auth.nonce,
		auth.username,
		auth.uri,
		auth.response,
	)
	if auth.qop == "auth" {
		str += fmt.Sprintf(`,qop=%s,nc=%s,cnonce="%s"`, auth.qop, auth.nc, auth.cnonce)
	}

	return str
}

// CalcResponse Authorization response https://www.ietf.org/rfc/rfc2617.txt
func CalcResponse(username, realm, password, method, uri, nonce, qop, cnonce, nc string) string {
	calcA1 := func() string {
		encoder := md5.New()
		encoder.Write([]byte(username + ":" + realm + ":" + password))

		return hex.EncodeToString(encoder.Sum(nil))
	}
	calcA2 := func() string {
		encoder := md5.New()
		encoder.Write([]byte(method + ":" + uri))

		return hex.EncodeToString(encoder.Sum(nil))
	}

	encoder := md5.New()
	encoder.Write([]byte(calcA1() + ":" + nonce + ":"))
	if qop != "" {
		encoder.Write([]byte(nc + ":" + cnonce + ":" + qop + ":"))
	}
	encoder.Write([]byte(calcA2()))
	return hex.EncodeToString(encoder.Sum(nil))
}
