package gateway

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtTokenHandler interface {
	Generate(key string, expirationInMs int64) (string, error)
	Validate(tokenStr string) (string, error)
}

type jwtTokenHandler struct {
	secret []byte
}

var _ JwtTokenHandler = (*jwtTokenHandler)(nil)

func NewJwtTokenHandler(secret string) JwtTokenHandler {
	return &jwtTokenHandler{
		secret: []byte(secret),
	}
}

func (j *jwtTokenHandler) Generate(key string, expirationInMs int64) (string, error) {
	expiration := time.Now().Add(time.Duration(expirationInMs) * time.Millisecond)

	claims := jwt.MapClaims{"key": key, "exp": expiration.Unix()}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.secret)
}

func (j *jwtTokenHandler) Validate(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", err
	}

	key, ok := claims["key"].(string)
	if !ok {
		return "", err
	}

	return key, nil
}
