package users

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/i-Things/things/shared/errors"
)

// 创建一个token
func CreateToken(secretKey string, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// 解析 token
func ParseToken(claim jwt.Claims, tokenString string, secretKey string) *errors.CodeError {
	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (i any, e error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return errors.TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return errors.TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return errors.TokenNotValidYet
			} else {
				return errors.TokenInvalid
			}
		}
	}
	if token != nil {
		if token.Valid {
			return nil
		}
		return errors.TokenInvalid

	} else {
		return errors.TokenInvalid
	}
}
