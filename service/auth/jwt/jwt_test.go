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

func TestCreateUserJWTToken(t *testing.T) {
	type request struct {
		userName string
	}

	type want struct {
		isError bool
	}

	type env struct {
		set   func()
		unset func()
	}

	testCases := []struct {
		name    string
		request request
		want    want
		env     env
	}{
		{
			name: "normal request",
			request: request{
				userName: "bob",
			},

			want: want{
				isError: false,
			},
			env: env{
				set: func() {
					os.Setenv("TOKEN_EXPIRE", "10")
					os.Setenv("CERT_VERSION", "1")
					os.Setenv("CERT_PATH", "../../../pkg/storage/certificates")
				},
				unset: func() {
					os.Unsetenv("TOKEN_EXPIRE")
					os.Unsetenv("CERT_VERSION")
					os.Unsetenv("CERT_PATH")
				},
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
			env: env{
				set: func() {
					os.Setenv("TOKEN_EXPIRE", "10")
					os.Setenv("CERT_VERSION", "1")
					os.Setenv("CERT_PATH", "./")
				},
				unset: func() {
					os.Unsetenv("TOKEN_EXPIRE")
					os.Unsetenv("CERT_VERSION")
					os.Unsetenv("CERT_PATH")
				},
			},
		},
		{
			name: "wrong time",
			request: request{
				userName: "bob",
			},
			want: want{
				isError: true,
			},
			env: env{
				set: func() {
					os.Setenv("TOKEN_EXPIRE", "")
					os.Setenv("CERT_VERSION", "1")
					os.Setenv("CERT_PATH", "./")
				},
				unset: func() {
					os.Unsetenv("TOKEN_EXPIRE")
					os.Unsetenv("CERT_VERSION")
					os.Unsetenv("CERT_PATH")
				},
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
			env: env{
				set: func() {
					os.Setenv("TOKEN_EXPIRE", "10")
					os.Setenv("CERT_VERSION", "0")
					os.Setenv("CERT_PATH", "../../../pkg/storage/certificates")
				},
				unset: func() {
					os.Unsetenv("TOKEN_EXPIRE")
					os.Unsetenv("CERT_VERSION")
					os.Unsetenv("CERT_PATH")
				},
			},
		},
	}

	for _, tc := range testCases {
		tc.env.set()
		defer tc.env.unset()

		token, err := CreateUserJWTToken(tc.request.userName)

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

	type env struct {
		set   func()
		unset func()
	}

	testCases := []struct {
		name    string
		request request
		want    want
		env     env
	}{
		{
			name: "normal request",
			request: request{
				certVersion: "1",
			},

			want: want{
				isError: false,
			},
			env: env{
				set: func() {
					os.Setenv("CERT_PATH", "../../../pkg/storage/certificates")
				},
				unset: func() {
					os.Unsetenv("CERT_PATH")
				},
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
			env: env{
				set: func() {
					os.Setenv("CERT_PATH", "./")
				},
				unset: func() {
					os.Unsetenv("CERT_PATH")
				},
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
			env: env{
				set: func() {
					os.Setenv("CERT_PATH", "../../../pkg/storage/certificates")
				},
				unset: func() {
					os.Unsetenv("CERT_PATH")
				},
			},
		},
	}

	for _, tc := range testCases {
		tc.env.set()
		defer tc.env.unset()

		cert, err := GetCertificateKey(tc.request.certVersion)

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
	//
	//JWT config
	os.Setenv("TOKEN_EXPIRE", "10") // minutes
	defer os.Unsetenv("TOKEN_EXPIRE")

	// certificates
	os.Setenv("CERT_VERSION", "1")
	defer os.Unsetenv("CERT_VERSION")
	os.Setenv("CERT_PATH", "../../../pkg/storage/certificates")
	defer os.Unsetenv("CERT_PATH")

	// create user
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
	pubInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)

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
