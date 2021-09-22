package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
)

type User interface {
	Get(context.Context, uuid.UUID) (model.UserHTTP, error)
	GetAll(context.Context) ([]model.UserHTTP, error)
	Create(context.Context, string) (model.UserHTTP, error)
	Update(context.Context, model.UserHTTP) (model.UserHTTP, error)
	Delete(context.Context, uuid.UUID) error
}
