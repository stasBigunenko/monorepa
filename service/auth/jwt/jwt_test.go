package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"monorepa/model"
	"os"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func init() {
	//JWT config
	os.Setenv("TOKEN_EXPIRE", "10") // minutes

	// certificates
	os.Setenv("CERT_VERSION", "1")
	os.Setenv("CERT_PATH", "../../../pkg/storage/certificates")
}

func TestGetCertificateKey(t *testing.T) {
	user := model.User{
		Name: "bob",
	}

	tokenString, err := CreateUserJWTToken(user.Name)
	assert.Nil(t, err)

	pbKeyBytes, err := GetCertificateKey(os.Getenv("CERT_VERSION"))
	assert.Nil(t, err)

	block, _ := pem.Decode(pbKeyBytes)
	assert.NotNil(t, block, "block")

	// pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	assert.Nil(t, err)
	pub := pubInterface.(*rsa.PublicKey)

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return pub, nil
	})
	assert.Nil(t, err)

	if claims, ok := token.Claims.(*UserClaims); ok {
		assert.Equal(t, user.Name, claims.Name)
		fmt.Println(claims.ExpiresAt)
	}
}
