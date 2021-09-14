package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	Name       string `json:"name"`
	KeyVersion string `json:"keyVersion"`
	jwt.StandardClaims
}

func (cl *UserClaims) addExpTime(duration int) {
	cl.ExpiresAt = time.Now().Add(time.Duration(duration) * time.Minute).Unix()
}
