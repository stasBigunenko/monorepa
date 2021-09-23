package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
)

type MockAccountsGrpcServer struct {
	MockCreateAccount   func(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	MockGetAccount      func(ctx context.Context, id uuid.UUID) (model.Account, error)
	MockGetUserAccounts func(ctx context.Context, userID uuid.UUID) ([]model.Account, error)
	MockGetAllAccounts  func(ctx context.Context) ([]model.Account, error)
	MockUpdateAccount   func(ctx context.Context, account model.Account) error
	MockDeleteAccount   func(ctx context.Context, id uuid.UUID) error
}

func (m *MockAccountsGrpcServer) CreateAccount(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	return m.MockCreateAccount(ctx, userID)
}
func (m *MockAccountsGrpcServer) GetAccount(ctx context.Context, id uuid.UUID) (model.Account, error) {
	return m.MockGetAccount(ctx, id)
}
func (m *MockAccountsGrpcServer) GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]model.Account, error) {
	return m.MockGetUserAccounts(ctx, userID)
}
func (m *MockAccountsGrpcServer) GetAllAccounts(ctx context.Context) ([]model.Account, error) {
	return m.MockGetAllAccounts(ctx)
}
func (m *MockAccountsGrpcServer) UpdateAccount(ctx context.Context, account model.Account) error {
	return m.MockUpdateAccount(ctx, account)
}
func (m *MockAccountsGrpcServer) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	return m.MockDeleteAccount(ctx, id)
}
