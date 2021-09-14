package jwt

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	pathCert            string
	certVersion         string
	tokenExpireDuration int
}

func NewJTWConfig() (*Config, error) {
	certVersion := os.Getenv("CERT_VERSION")
	if certVersion == "" {
		return nil, errors.New("wrong cert path")
	}

	expireTime := os.Getenv("TOKEN_EXPIRE")
	duration, err := strconv.Atoi(expireTime)
	if err != nil {
		return nil, err
	}

	path := os.Getenv("CERT_PATH")
	if path == "" {
		return nil, errors.New("wrong cert path")
	}

	token := &Config{
		pathCert:            path,
		certVersion:         certVersion,
		tokenExpireDuration: duration,
	}

	return token, nil
}
