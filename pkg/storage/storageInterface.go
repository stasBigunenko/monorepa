package storage

import "context"

type StorageInterface interface {
	GetItems(context.Context, string) ([]Storage, error)
}
