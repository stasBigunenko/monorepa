package jwt

import "github.com/golang-jwt/jwt"

type UserClaims struct {
	Name       string `json:"name"`
	KeyVersion string `json:"keyVersion"`
	jwt.StandardClaims
}
