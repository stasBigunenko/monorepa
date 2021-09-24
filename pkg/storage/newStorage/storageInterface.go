package newStorage

import (
	"context"
	"github.com/google/uuid"
)

type NewStore interface {
	Get(context.Context, uuid.UUID) (interface{}, error)
	GetUserAccounts(context.Context, uuid.UUID) (interface{}, error)
	GetAll(context.Context) (interface{}, error)
	Create(context.Context, interface{}) (interface{}, error)
	Update(context.Context, interface{}) (interface{}, error)
	Delete(context.Context, uuid.UUID) error
}
