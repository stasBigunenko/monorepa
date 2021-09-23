package mocks

import (
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
)

type MockAccountsGrpcServer struct {
	MockCreateAccount   func(userID uuid.UUID) (uuid.UUID, error)
	MockGetAccount      func(id uuid.UUID) (model.Account, error)
	MockGetUserAccounts func(userID uuid.UUID) ([]model.Account, error)
	MockGetAllAccounts  func() ([]model.Account, error)
	MockUpdateAccount   func(account model.Account) error
	MockDeleteAccount   func(id uuid.UUID) error
}

func (m *MockAccountsGrpcServer) CreateAccount(userID uuid.UUID) (uuid.UUID, error) {
	return m.MockCreateAccount(userID)
}
func (m *MockAccountsGrpcServer) GetAccount(id uuid.UUID) (model.Account, error) {
	return m.MockGetAccount(id)
}
func (m *MockAccountsGrpcServer) GetUserAccounts(userID uuid.UUID) ([]model.Account, error) {
	return m.MockGetUserAccounts(userID)
}
func (m *MockAccountsGrpcServer) GetAllAccounts() ([]model.Account, error) {
	return m.MockGetAllAccounts()
}
func (m *MockAccountsGrpcServer) UpdateAccount(account model.Account) error {
	return m.MockUpdateAccount(account)
}
func (m *MockAccountsGrpcServer) DeleteAccount(id uuid.UUID) error {
	return m.MockDeleteAccount(id)
}
