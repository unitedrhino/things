package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"gitee.com/godLei6/things/shared/errors"
)

// Custom claims structure
type CustomClaims struct {
	Uid      int64
	jwt.StandardClaims
}



func  GetJwtToken(secretKey string, iat, seconds, uid int64) (string, error) {
	claims := CustomClaims{
		Uid: uid,
		StandardClaims:jwt.StandardClaims{
			ExpiresAt:iat+seconds,
			IssuedAt:iat,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// 创建一个token
func CreateToken(secretKey string,claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// 解析 token
func ParseToken(tokenString string,secretKey string) (*CustomClaims, *errors.CodeError) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
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
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, errors.TokenInvalid

	} else {
		return nil, errors.TokenInvalid
	}

}

// 更新token
func RefreshToken(tokenString string,secretKey string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return CreateToken(secretKey,*claims)
	}
	return "", errors.TokenInvalid
}