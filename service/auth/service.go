package auth

import (
	"errors"
	"monorepa/model"
	"monorepa/service/auth/jwt"
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
		return "", errors.New("wrong password")
	}

	token, err := jwt.CreateUserJWTToken(user.Name)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Get certificate for user
func (s *Session) GetCert(keyCertVersion string) ([]byte, error) {
	res, err := jwt.GetCertificateKey(keyCertVersion)
	if err != nil {
		return nil, err
	}

	return res, nil
}
