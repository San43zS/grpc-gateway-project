package services

import (
	"github.com/golang-jwt/jwt"
	"time"
)

const secret = "secret"

func CreateToken(usrEmail string) (string, error) {
	claims := jwt.MapClaims{}
	claims["email"] = usrEmail
	claims["exp"] = time.Now()

	tokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenJwt.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, err
}
