package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang-jwt/jwt"
)

func getCertName(path, version string) string {
	return fmt.Sprintf("%v/private_key%v.pem", path, version)
}

// read certificate as RSA key and certificate version
func readRSAPrivateKey(certVersion string) (*rsa.PrivateKey, error) {
	// data to get certificate
	path := os.Getenv("CERT_PATH")
	if path == "" {
		return nil, errors.New("wrong cert path")
	}

	// read certificate as byte array
	b, err := ioutil.ReadFile(getCertName(path, certVersion))
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
