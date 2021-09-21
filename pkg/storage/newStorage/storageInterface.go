package newStorage

import (
	"context"
	"github.com/google/uuid"
)

type NewStore interface {
	Get(context.Context, uuid.UUID) ([]byte, error)
	GetUser(context.Context, uuid.UUID) ([][]byte, error)
	GetAll(context.Context) ([][]byte, error)
	Create(context.Context, []byte) ([]byte, uuid.UUID, error)
	Update(context.Context, uuid.UUID, []byte) ([]byte, error)
	Delete(context.Context, uuid.UUID) (bool, error)
}
