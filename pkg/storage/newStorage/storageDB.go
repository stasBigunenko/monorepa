package newStorage

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/google/uuid"

	"github.com/stasBigunenko/monorepa/model"
)

type LoggingService interface {
	WriteLog(ctx context.Context, message string)
}

type StorageDB struct {
	Data           map[uuid.UUID][]byte
	mu             sync.Mutex
	loggingService LoggingService
}

func NewDB(loggingService LoggingService) *StorageDB {
	sdb := StorageDB{}
	sdb.Data = make(map[uuid.UUID][]byte)
	sdb.loggingService = loggingService
	return &sdb
}

func (sdb *StorageDB) Get(c context.Context, id uuid.UUID) ([]byte, error) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	sdb.loggingService.WriteLog(c, "Storage: Command Get received...")

	val, ok := sdb.Data[id]

	if !ok {
		return nil, nil
	}

	return val, nil
}

func (sdb *StorageDB) GetUser(c context.Context, userID uuid.UUID) ([][]byte, error) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	sdb.loggingService.WriteLog(c, "Storage: Command GetUser received...")

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

func (sdb *StorageDB) GetAll(c context.Context) ([][]byte, error) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	sdb.loggingService.WriteLog(c, "Storage: Command GetAll received...")

	var res [][]byte

	for _, val := range sdb.Data {
		res = append(res, val)
	}

	return res, nil
}

func (sdb *StorageDB) Create(c context.Context, b []byte) ([]byte, uuid.UUID, error) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	sdb.loggingService.WriteLog(c, "Storage: Command Create received...")

	id := uuid.New()

	sdb.Data[id] = b

	return b, id, nil
}

func (sdb *StorageDB) Update(c context.Context, id uuid.UUID, b []byte) ([]byte, error) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	sdb.loggingService.WriteLog(c, "Storage: Command Update received...")

	if _, ok := sdb.Data[id]; !ok {
		return nil, nil
	}

	sdb.Data[id] = b

	return b, nil
}

func (sdb *StorageDB) Delete(c context.Context, id uuid.UUID) (bool, error) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	sdb.loggingService.WriteLog(c, "Storage: Command Delete received...")

	if _, ok := sdb.Data[id]; !ok {
		return false, nil
	}

	delete(sdb.Data, id)

	return true, nil
}
