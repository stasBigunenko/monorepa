package accountsDB

import (
	"context"
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
)

type AccountsService interface {
	GetAccount(context.Context, uuid.UUID) (model.Account, error)
	GetUserAccounts(context.Context, uuid.UUID) ([]model.Account, error)
	GetAllUsers(context.Context) ([]model.Account, error)
	CreateAccount(context.Context, uuid.UUID) (model.Account, error)
	UpdateAccount(context.Context, model.Account) (model.Account, error)
	DeleteAccount(context.Context, uuid.UUID) error
}
