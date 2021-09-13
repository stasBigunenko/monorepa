package auth

import "monorepa/model"

type Session struct {
}

type AuthService interface {
	Login(model.User) (string, error)
	GetCert(string) ([]byte, error)
}
