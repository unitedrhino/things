package devices

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/i-Things/things/shared/errors"
	"time"
)

// Custom claims structure
type OssJwtToken struct {
	Bucket string //oss的token
	Dir    string //对象路径
	jwt.StandardClaims
}

func GetJwtToken(secretKey string, iat, seconds int64, bucket string, dir string) (string, error) {
	claims := OssJwtToken{
		Bucket: bucket,
		Dir:    dir,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: iat + seconds,
			IssuedAt:  iat,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// 创建一个token
func CreateToken(secretKey string, claims OssJwtToken) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// 解析 token
func ParseToken(tokenString string, secretKey string) (*OssJwtToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &OssJwtToken{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, errors.TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.TokenNotValidYet
			} else {
				return nil, errors.TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*OssJwtToken); ok && token.Valid {
			return claims, nil
		}
		return nil, errors.TokenInvalid

	} else {
		return nil, errors.TokenInvalid
	}

}

// 更新token
func RefreshToken(tokenString string, secretKey string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &OssJwtToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*OssJwtToken); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return CreateToken(secretKey, *claims)
	}
	return "", errors.TokenInvalid
}
