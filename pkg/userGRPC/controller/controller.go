package usergrpccontroller

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	customerrors "github.com/stasBigunenko/monorepa/customErrors"
	"github.com/stasBigunenko/monorepa/model"
	pb "github.com/stasBigunenko/monorepa/pkg/userGRPC/proto"
)

type LoggingService interface {
	WriteLog(ctx context.Context, message string)
}

type UserGRPCСontroller struct {
	client         pb.UserGRPCServiceClient
	loggingService LoggingService
}

func New(cli pb.UserGRPCServiceClient, loggingService LoggingService) *UserGRPCСontroller {
	return &UserGRPCСontroller{
		client:         cli,
		loggingService: loggingService,
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

func (s UserGRPCСontroller) CreateUser(ctx context.Context, name string) (uuid.UUID, error) {
	s.loggingService.WriteLog(ctx, "GRPC Client: Command CreateUser received...")

	contextID, ok := ctx.Value(model.ContextKeyRequestID).(string)
	if !ok {
		log.Info("failed to convert context value and get context id")
	}

	c := metadata.AppendToOutgoingContext(ctx, "requestid", contextID)

	resp, err := s.client.Create(c, &pb.Name{
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

func (s UserGRPCСontroller) GetUser(ctx context.Context, id uuid.UUID) (model.UserHTTP, error) {
	s.loggingService.WriteLog(ctx, "GRPC Client: Command GetUser received...")

	contextID, ok := ctx.Value(model.ContextKeyRequestID).(string)
	if !ok {
		log.Info("failed to convert context value and get context id")
	}

	c := metadata.AppendToOutgoingContext(ctx, "requestid", contextID)

	resp, err := s.client.Get(c, &pb.Id{
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

func (s UserGRPCСontroller) GetAllUsers(ctx context.Context) ([]model.UserHTTP, error) {
	s.loggingService.WriteLog(ctx, "GRPC Client: Command GetAllUsers received...")

	contextID, ok := ctx.Value(model.ContextKeyRequestID).(string)
	if !ok {
		log.Info("failed to convert context value and get context id")
	}

	c := metadata.AppendToOutgoingContext(ctx, "requestid", contextID)

	resp, err := s.client.GetAllUsers(c, &emptypb.Empty{})
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

func (s UserGRPCСontroller) UpdateUser(ctx context.Context, user model.UserHTTP) error {
	s.loggingService.WriteLog(ctx, "GRPC Client: Command UpdateUser received...")

	contextID, ok := ctx.Value(model.ContextKeyRequestID).(string)
	if !ok {
		log.Info("failed to convert context value and get context id")
	}

	c := metadata.AppendToOutgoingContext(ctx, "requestid", contextID)

	_, err := s.client.Update(c, &pb.User{
		Id:   user.ID.String(),
		Name: user.Name,
	})

	if err != nil {
		return s.formatError(err, "failed to update user")
	}

	return nil
}

func (s UserGRPCСontroller) DeleteUser(ctx context.Context, id uuid.UUID) error {
	s.loggingService.WriteLog(ctx, "GRPC Client: Command DeleteUser received...")

	contextID, ok := ctx.Value(model.ContextKeyRequestID).(string)
	if !ok {
		log.Info("failed to convert context value and get context id")
	}

	c := metadata.AppendToOutgoingContext(ctx, "requestid", contextID)

	_, err := s.client.Delete(c, &pb.Id{
		Id: id.String(),
	})

	if err != nil {
		return s.formatError(err, "failed to delete user")
	}

	return nil
}
