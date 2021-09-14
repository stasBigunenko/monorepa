package storage

import "context"

type StorageItemService interface {
	GetItems(context.Context, string) ([]StorageItem, error)
}
