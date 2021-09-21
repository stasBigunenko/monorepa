package newStorage

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
	"sync"
)

type StorageDB struct {
	Data map[uuid.UUID][]byte
	mu   sync.Mutex
}

func NewDB() *StorageDB {
	sdb := StorageDB{}
	sdb.Data = make(map[uuid.UUID][]byte)
	return &sdb
}

func (sdb *StorageDB) Get(_ context.Context, id uuid.UUID) ([]byte, error) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	val, ok := sdb.Data[id]

	if !ok {
		return nil, nil
	}

	return val, nil
}

func (sdb *StorageDB) GetUser(_ context.Context, userID uuid.UUID) ([][]byte, error) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	var res [][]byte

	for _, account := range sdb.Data {
		acc := model.Account{}
		err := json.Unmarshal(account, &acc)
		if err != nil {
			return nil, nil
		}
		if userID == acc.UserID {
			r, err := json.Marshal(acc)
			if err != nil {
				return nil, nil
			}
			res = append(res, r)
		}
	}
	return res, nil
}

func (sdb *StorageDB) GetAll(_ context.Context) ([][]byte, error) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	var res [][]byte

	for _, val := range sdb.Data {
		res = append(res, val)
	}

	return res, nil
}

func (sdb *StorageDB) Create(_ context.Context, b []byte) ([]byte, uuid.UUID, error) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	id := uuid.New()

	sdb.Data[id] = b

	return b, id, nil
}

func (sdb *StorageDB) Update(_ context.Context, id uuid.UUID, b []byte) ([]byte, error) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	if _, ok := sdb.Data[id]; !ok {
		return nil, nil
	}

	sdb.Data[id] = b

	return b, nil
}

func (sdb *StorageDB) Delete(_ context.Context, id uuid.UUID) (bool, error) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	if _, ok := sdb.Data[id]; !ok {
		return false, nil
	}

	delete(sdb.Data, id)

	return true, nil
}
