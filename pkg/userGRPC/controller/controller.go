package usergrpccontroller

import (
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
	pb "github.com/stasBigunenko/monorepa/pkg/userGRPC/proto"
)

type UserGRPCСontroller struct {
	client pb.UserGRPCServiceClient
}

func New(cli pb.UserGRPCServiceClient) UserGRPCСontroller {
	return UserGRPCСontroller{
		client: cli,
	}
}

func (s UserGRPCСontroller) CreateUser(name string) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (s UserGRPCСontroller) GetUser(id uuid.UUID) (model.UserHTTP, error) {
	return model.UserHTTP{}, nil
}

func (s UserGRPCСontroller) GetAllUsers() ([]model.UserHTTP, error) {
	return nil, nil
}

func (s UserGRPCСontroller) UpdateUser(model.UserHTTP) error {
	return nil
}

func (s UserGRPCСontroller) DeleteUser(id uuid.UUID) error {
	return nil
}
