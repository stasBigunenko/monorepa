package mocks

import (
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
)

type MockUsersGrpcServer struct {
	MockCreateUser  func(name string) (uuid.UUID, error)
	MockGetUser     func(id uuid.UUID) (model.UserHTTP, error)
	MockGetAllUsers func() ([]model.UserHTTP, error)
	MockUpdateUser  func(user model.UserHTTP) error
	MockDeleteUser  func(id uuid.UUID) error
}

func (m *MockUsersGrpcServer) CreateUser(name string) (uuid.UUID, error) {
	return m.MockCreateUser(name)
}
func (m *MockUsersGrpcServer) GetUser(id uuid.UUID) (model.UserHTTP, error) {
	return m.MockGetUser(id)
}
func (m *MockUsersGrpcServer) GetAllUsers() ([]model.UserHTTP, error) {
	return m.MockGetAllUsers()
}
func (m *MockUsersGrpcServer) UpdateUser(user model.UserHTTP) error {
	return m.MockUpdateUser(user)
}
func (m *MockUsersGrpcServer) DeleteUser(id uuid.UUID) error {
	return m.MockDeleteUser(id)
}
