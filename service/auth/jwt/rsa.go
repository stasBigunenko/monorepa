package jwt

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/golang-jwt/jwt"
)

func getCertName(path, version string) string {
	certName := fmt.Sprintf("private_key%v.pem", version)

	return filepath.Join(path, certName)
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
