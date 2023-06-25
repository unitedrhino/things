package users

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/i-Things/things/shared/errors"
	"time"
)

// Custom claims structure
type LoginClaims struct {
	UserID    int64 `json:",string"`
	Role      int64
	IsAllData int64
	jwt.StandardClaims
}

func GetLoginJwtToken(secretKey string, iat, seconds, userID, role, isAllData int64) (string, error) {
	claims := LoginClaims{
		UserID:    userID,
		Role:      role,
		IsAllData: isAllData,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: iat + seconds,
			IssuedAt:  iat,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// 更新token
func RefreshLoginToken(tokenString string, secretKey string, AccessExpire int64) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &LoginClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*LoginClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = AccessExpire
		return CreateToken(secretKey, *claims)
	}
	return "", errors.TokenInvalid
}
