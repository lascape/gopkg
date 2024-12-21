package authx

import (
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-jwt/jwt/v5"
	"github.com/lascape/gopkg/envx"
	"github.com/pkg/errors"
	"time"
)

var jwtSecret = envx.ValueByEnv("PKG_AUTHX_JWT_SECRET", "your_secret_key").Bytes()

type JWT struct{}

func (j *JWT) GenerateToken(kv map[string]interface{}, expired time.Duration) (string, error) {
	claims := jwt.MapClaims(kv)
	claims["exp"] = float64(time.Now().Add(expired).Unix())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (j *JWT) ValidateToken(token string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
