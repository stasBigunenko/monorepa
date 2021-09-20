package usersDB

import (
	"context"
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
)

type UsersService interface {
	Get(context.Context, uuid.UUID) (model.UserHTTP, error)
	GetAllUsers(context.Context) ([]model.UserHTTP, error)
	Create(context.Context, string) (model.UserHTTP, error)
	Update(context.Context, model.UserHTTP) (model.UserHTTP, error)
	Delete(context.Context, uuid.UUID) error
}
