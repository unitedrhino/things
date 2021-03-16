package logic

import "github.com/dgrijalva/jwt-go"

func  getJwtToken(secretKey string, iat, seconds, uid int64,tokenType string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["uid"] = uid
	claims["type"] = tokenType
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}