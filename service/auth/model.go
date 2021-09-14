package auth

import "github.com/stasBigunenko/monorepa/model"

type Session struct {
}

type Service interface {
	Login(model.User) (string, error)
	GetCert(string) ([]byte, error)
}
