package account

import (
	"context"
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
)

type AccountInterface interface {
	Get(context.Context, uuid.UUID) (model.Account, error)
	GetUser(context.Context, uuid.UUID) ([]model.Account, error)
	GetAll(context.Context) ([]model.Account, error)
	Create(context.Context, uuid.UUID) (model.Account, error)
	Update(context.Context, model.Account) (model.Account, error)
	Delete(context.Context, uuid.UUID) error
}
