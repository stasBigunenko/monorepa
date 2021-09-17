package model

import "github.com/golang-jwt/jwt"

type JWTUserClaims struct {
	Name       string `json:"name"`
	KeyVersion string `json:"keyVersion"`
	jwt.StandardClaims
}
