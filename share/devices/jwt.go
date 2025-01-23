package devices

import (
	"gitee.com/unitedrhino/share/errors"
	//"github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// OssJwtToken Custom claims structure
type OssJwtToken struct {
	Bucket string //oss的token
	Dir    string //对象路径
	jwt.RegisteredClaims
}

func GetJwtToken(secretKey string, t time.Time, seconds int64, bucket string, dir string) (string, error) {
	IssuedAt := jwt.NewNumericDate(t)
	claims := OssJwtToken{
		Bucket: bucket,
		Dir:    dir,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(t.Add(time.Duration(seconds) * time.Second)),
			IssuedAt:  IssuedAt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// CreateToken 创建一个token
func CreateToken(secretKey string, claims OssJwtToken) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// 解析 token
func ParseToken(tokenString string, secretKey string) (*OssJwtToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &OssJwtToken{}, func(token *jwt.Token) (i any, e error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, errors.TokenExpired.WithMsg("登录失效,请退出重新登录")
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, errors.TokenMalformed.WithMsg("登录失效,请退出重新登录")
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, errors.TokenNotValidYet.WithMsg("登录失效,请退出重新登录")
		default:
			return nil, errors.TokenInvalid.WithMsg("登录失效,请退出重新登录")
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*OssJwtToken); ok && token.Valid {
			return claims, nil
		}
		return nil, errors.TokenInvalid.WithMsg("登录失效,请退出重新登录")

	} else {
		return nil, errors.TokenInvalid.WithMsg("登录失效,请退出重新登录")
	}

}

// 更新token
func RefreshToken(tokenString string, secretKey string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &OssJwtToken{}, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*OssJwtToken); ok && token.Valid {
		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
		return CreateToken(secretKey, *claims)
	}
	return "", errors.TokenInvalid
}
