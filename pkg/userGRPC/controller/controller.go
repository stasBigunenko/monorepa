package usergrpccontroller

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	customerrors "github.com/stasBigunenko/monorepa/customErrors"
	"github.com/stasBigunenko/monorepa/model"
	pb "github.com/stasBigunenko/monorepa/pkg/userGRPC/proto"
)

type UserGRPCСontroller struct {
	client pb.UserGRPCServiceClient
}

func New(cli pb.UserGRPCServiceClient) *UserGRPCСontroller {
	return &UserGRPCСontroller{
		client: cli,
	}
}

func (s UserGRPCСontroller) formatError(err error, message string) error {
	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("%s, failed to parse status or not a grpc error type: %w", message, err)
	}

	switch st.Code() {
	case codes.NotFound:
		return fmt.Errorf("%s: %w", message, customerrors.NotFound)
	case codes.AlreadyExists:
		return fmt.Errorf("%s: %w", message, customerrors.AlreadyExists)
	case codes.DeadlineExceeded:
		return fmt.Errorf("%s: %w", message, customerrors.DeadlineExceeded)
	}

	return fmt.Errorf("%s: %s", message, err.Error())
}

func (s UserGRPCСontroller) CreateUser(name string) (uuid.UUID, error) {
	resp, err := s.client.Create(context.Background(), &pb.Name{
		Name: name,
	})

	if err != nil {
		return uuid.Nil, s.formatError(err, "failed to create user")
	}

	userID, err := uuid.Parse(resp.Id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse user ID: %s, %w", err.Error(), customerrors.ParseError)
	}

	return userID, nil
}

func (s UserGRPCСontroller) GetUser(id uuid.UUID) (model.UserHTTP, error) {
	resp, err := s.client.Get(context.Background(), &pb.Id{
		Id: id.String(),
	})

	if err != nil {
		return model.UserHTTP{}, s.formatError(err, "failed to get user")
	}

	userID, err := uuid.Parse(resp.Id)
	if err != nil {
		return model.UserHTTP{}, fmt.Errorf("failed to parse user ID: %s, %w", err.Error(), customerrors.ParseError)
	}

	return model.UserHTTP{
		ID:   userID,
		Name: resp.Name,
	}, nil
}

func (s UserGRPCСontroller) GetAllUsers() ([]model.UserHTTP, error) {
	resp, err := s.client.GetAllUsers(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, s.formatError(err, "failed to get all users")
	}

	users := []model.UserHTTP{}
	for _, user := range resp.AllUsers {
		id, err := uuid.Parse(user.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to parse user ID: %s, %w", err.Error(), customerrors.ParseError)
		}

		users = append(users, model.UserHTTP{
			ID:   id,
			Name: user.Name,
		})
	}

	return users, nil
}

func (s UserGRPCСontroller) UpdateUser(user model.UserHTTP) error {
	_, err := s.client.Update(context.Background(), &pb.User{
		Id:   user.ID.String(),
		Name: user.Name,
	})

	if err != nil {
		return s.formatError(err, "failed to update user")
	}

	return nil
}

func (s UserGRPCСontroller) DeleteUser(id uuid.UUID) error {
	_, err := s.client.Delete(context.Background(), &pb.Id{
		Id: id.String(),
	})

	if err != nil {
		return s.formatError(err, "failed to delete user")
	}

	return nil
}
