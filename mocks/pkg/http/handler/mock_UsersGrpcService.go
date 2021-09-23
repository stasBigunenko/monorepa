package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
)

type MockUsersGrpcServer struct {
	MockCreateUser  func(ctx context.Context, name string) (uuid.UUID, error)
	MockGetUser     func(ctx context.Context, id uuid.UUID) (model.UserHTTP, error)
	MockGetAllUsers func(ctx context.Context) ([]model.UserHTTP, error)
	MockUpdateUser  func(ctx context.Context, user model.UserHTTP) error
	MockDeleteUser  func(ctx context.Context, id uuid.UUID) error
}

func (m *MockUsersGrpcServer) CreateUser(ctx context.Context, name string) (uuid.UUID, error) {
	return m.MockCreateUser(ctx, name)
}
func (m *MockUsersGrpcServer) GetUser(ctx context.Context, id uuid.UUID) (model.UserHTTP, error) {
	return m.MockGetUser(ctx, id)
}
func (m *MockUsersGrpcServer) GetAllUsers(ctx context.Context) ([]model.UserHTTP, error) {
	return m.MockGetAllUsers(ctx)
}
func (m *MockUsersGrpcServer) UpdateUser(ctx context.Context, user model.UserHTTP) error {
	return m.MockUpdateUser(ctx, user)
}
func (m *MockUsersGrpcServer) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return m.MockDeleteUser(ctx, id)
}
