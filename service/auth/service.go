package auth

import (
	er "github.com/stasBigunenko/monorepa/errors"
	"github.com/stasBigunenko/monorepa/model"
	"github.com/stasBigunenko/monorepa/service/auth/jwt"
)

// init item services
func New() *Session {
	return &Session{}
}

// verify user for login
func userVerify(password string) bool {
	return password != ""
}

// Create JWT tocken for user
func (s *Session) Login(user model.User) (string, error) {
	if ok := userVerify(user.Password); !ok {
		return "", er.WrongPassword
	}

	conf, err := jwt.NewJTWConfig()
	if err != nil {
		return "", err
	}

	token, err := jwt.CreateUserJWTToken(user.Name, conf)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Get certificate for user
func (s *Session) GetCert(keyCertVersion string) ([]byte, error) {
	conf, err := jwt.NewJTWConfig()
	if err != nil {
		return nil, err
	}

	res, err := jwt.GetCertificateKey(keyCertVersion, conf)
	if err != nil {
		return nil, err
	}

	return res, nil
}
