package jwt

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func newClaim(name, version string) *UserClaims {
	return &UserClaims{
		Name:       name,
		KeyVersion: version,
	}
}

func addExpTime(claims *UserClaims) error {
	expireTime := os.Getenv("TOKEN_EXPIRE")
	duration, err := strconv.Atoi(expireTime)
	if err != nil {
		return err
	}

	claims.ExpiresAt = time.Now().Add(time.Duration(duration) * time.Minute).Unix()

	return nil
}

// create new tocken for User
func CreateUserJWTToken(userName string) (string, error) {
	certVersion := os.Getenv("CERT_VERSION")
	if certVersion == "" {
		return "", errors.New("wrong cert path")
	}

	privateKey, err := readRSAPrivateKey(certVersion)
	if err != nil {
		return "", err
	}

	claims := newClaim(userName, certVersion)
	err = addExpTime(claims)
	if err != nil {
		return "", err
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodRS512, claims).SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// get certifikate of some version
func GetCertificateKey(certVersion string) ([]byte, error) {
	privateKey, err := readRSAPrivateKey(certVersion)
	if err != nil {
		return nil, err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	return pem.EncodeToMemory(publicKeyBlock), nil
	// return publicKeyBytes, nil

}
