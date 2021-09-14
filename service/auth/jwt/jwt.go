package jwt

import (
	"crypto/x509"
	"encoding/pem"

	"github.com/golang-jwt/jwt"
)

func newClaim(name, version string) *UserClaims {
	return &UserClaims{
		Name:       name,
		KeyVersion: version,
	}
}

// create new tocken for User
func CreateUserJWTToken(userName string, conf *Config) (string, error) {
	privateKey, err := readRSAPrivateKey(conf.certVersion, conf.pathCert)
	if err != nil {
		return "", err
	}

	claims := newClaim(userName, conf.certVersion)
	claims.addExpTime(conf.tokenExpireDuration)

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodRS512, claims).SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// get certifikate of some version
func GetCertificateKey(certVersion string, conf *Config) ([]byte, error) {
	privateKey, err := readRSAPrivateKey(certVersion, conf.pathCert)
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
}
