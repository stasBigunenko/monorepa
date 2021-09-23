package httphandler

import (
	"context"

	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
)

type AccountGrpcService interface {
	CreateAccount(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	GetAccount(ctx context.Context, id uuid.UUID) (model.Account, error)
	GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]model.Account, error)
	GetAllAccounts(ctx context.Context) ([]model.Account, error)
	UpdateAccount(ctx context.Context, account model.Account) error
	DeleteAccount(ctx context.Context, id uuid.UUID) error
}

type UserGrpcService interface {
	CreateUser(ctx context.Context, name string) (uuid.UUID, error)
	GetUser(ctx context.Context, id uuid.UUID) (model.UserHTTP, error)
	GetAllUsers(ctx context.Context) ([]model.UserHTTP, error)
	UpdateUser(ctx context.Context, user model.UserHTTP) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type TokenService interface {
	ParseToken(tokenPart string) (string, error)
}

type LoggingService interface {
	WriteLog(ctx context.Context, message string)
}
