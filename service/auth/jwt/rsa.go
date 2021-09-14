package jwt

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"

	"github.com/golang-jwt/jwt"
)

func getCertName(path, version string) string {
	return fmt.Sprintf("%v/private_key%v.pem", path, version)
}

// read certificate as RSA key and certificate version
func readRSAPrivateKey(certVersion string, pathCert string) (*rsa.PrivateKey, error) {
	// read certificate as byte array
	b, err := ioutil.ReadFile(getCertName(pathCert, certVersion))
	if err != nil {
		return nil, err
	}

	// convert to private key
	rsaPrivKey, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		return nil, err
	}

	return rsaPrivKey, err
}
