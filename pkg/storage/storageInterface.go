package storage

type StorageInterface interface {
	GetItems(username string) ([]Storage, error)
}
