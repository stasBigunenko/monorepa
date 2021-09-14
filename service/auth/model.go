package auth

import "monorepa/model"

type Session struct {
}

type Service interface {
	Login(model.User) (string, error)
	GetCert(string) ([]byte, error)
}
