package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"monorepa/model"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserJWTToken(t *testing.T) {
	type request struct {
		userName string
	}

	type want struct {
		isError bool
	}

	testCases := []struct {
		name      string
		request   request
		want      want
		envConfig Config
	}{
		{
			name:    "normal request",
			request: request{userName: "bob"},
			want:    want{isError: false},
			envConfig: Config{
				pathCert:            "../../../pkg/storage/certificates",
				certVersion:         "1",
				tokenExpireDuration: 10,
			},
		},
		{
			name: "wrong path",
			request: request{
				userName: "bob",
			},
			want: want{
				isError: true,
			},
			envConfig: Config{
				pathCert:            "./",
				certVersion:         "1",
				tokenExpireDuration: 10,
			},
		},
		{
			name: "wrong cert version",
			request: request{
				userName: "bob",
			},
			want: want{
				isError: true,
			},
			envConfig: Config{
				pathCert:            "../../../pkg/storage/certificates",
				certVersion:         "0",
				tokenExpireDuration: 10,
			},
		},
	}

	for _, tc := range testCases {
		token, err := CreateUserJWTToken(tc.request.userName, &tc.envConfig)

		if tc.want.isError {
			assert.Equal(t, token, "", tc.name)
			assert.Error(t, err, tc.name)
		} else {
			assert.NotEqual(t, token, "", tc.name)
			assert.Nil(t, err, tc.name)
		}
	}
}

func TestGetCertificateKey(t *testing.T) {
	type request struct {
		certVersion string
	}

	type want struct {
		isError bool
	}

	testCases := []struct {
		name      string
		request   request
		want      want
		envConfig Config
	}{
		{
			name: "normal request",
			request: request{
				certVersion: "1",
			},

			want: want{
				isError: false,
			},
			envConfig: Config{
				pathCert: "../../../pkg/storage/certificates",
			},
		},
		{
			name: "wrong path",
			request: request{
				certVersion: "1",
			},
			want: want{
				isError: true,
			},
			envConfig: Config{
				pathCert: "./",
			},
		},
		{
			name: "wrong cert version",
			request: request{
				certVersion: "0",
			},
			want: want{
				isError: true,
			},
			envConfig: Config{
				pathCert: "./",
			},
		},
	}

	for _, tc := range testCases {

		cert, err := GetCertificateKey(tc.request.certVersion, &tc.envConfig)

		if tc.want.isError {
			assert.Nil(t, cert, tc.name)
			assert.Error(t, err, tc.name)
		} else {
			assert.NotNil(t, cert, tc.name)
			assert.Nil(t, err, tc.name)
		}
	}
}

// general test of function work
func TesCertParsing(t *testing.T) {

	// config env
	// JWT config
	conf := Config{
		certVersion:         "1",
		tokenExpireDuration: 10, // minutes
		pathCert:            "../../../pkg/storage/certificates",
	}

	// create user
	user := model.User{
		Name: "bob",
	}

	tokenString, err := CreateUserJWTToken(user.Name, &conf)
	assert.Nil(t, err)

	pbKeyBytes, err := GetCertificateKey(conf.certVersion, &conf)
	assert.Nil(t, err)

	block, _ := pem.Decode(pbKeyBytes)
	assert.NotNil(t, block, "block")

	pubInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)

	pub, ok := pubInterface.(*rsa.PublicKey)
	assert.Equal(t, true, ok)

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return pub, nil
	})
	assert.Nil(t, err)

	if claims, ok := token.Claims.(*UserClaims); ok {
		assert.Equal(t, user.Name, claims.Name)
		fmt.Println(claims.ExpiresAt)
	}
}
