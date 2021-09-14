package storage

import "context"

type ItemService interface {
	GetItems(context.Context, string) ([]StorageItem, error)
}
