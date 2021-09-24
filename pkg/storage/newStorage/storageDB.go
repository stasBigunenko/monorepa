package newStorage

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"

	"github.com/stasBigunenko/monorepa/model"
)

type LoggingService interface {
	WriteLog(ctx context.Context, message string)
}

type StorageDB struct {
	Data           map[uuid.UUID]interface{}
	mu             sync.Mutex
	loggingService LoggingService
}

func NewDB(loggingService LoggingService) *StorageDB {
	sdb := StorageDB{}
	sdb.Data = make(map[uuid.UUID]interface{})
	sdb.loggingService = loggingService
	return &sdb
}

func (sdb *StorageDB) Get(c context.Context, id uuid.UUID) (interface{}, error) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	sdb.loggingService.WriteLog(c, "Storage: Command Get received...")

	val, ok := sdb.Data[id]
	if !ok {
		return nil, errors.New("not found")
	}

	return val, nil
}

func (sdb *StorageDB) GetUserAccounts(c context.Context, userID uuid.UUID) (interface{}, error) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	sdb.loggingService.WriteLog(c, "Storage: Command GetUserAccounts received...")

	var res []model.Account

	for _, account := range sdb.Data {
		acc, ok := account.(model.Account)
		if !ok {
			return nil, errors.New("invalid userID")
		}
		if userID == acc.UserID {
			res = append(res, acc)
		}
	}
	return res, nil
}

func (sdb *StorageDB) GetAll(c context.Context) (interface{}, error) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	sdb.loggingService.WriteLog(c, "Storage: Command GetAll received...")

	var users []model.UserHTTP
	var accounts []model.Account

	for _, val := range sdb.Data {
		res, ok := val.(model.UserHTTP)
		if ok {
			users = append(users, res)
			continue
		}

		res2, ok := val.(model.Account)
		if ok {
			accounts = append(accounts, res2)
		}
	}

	if users != nil { //nolint
		return users, nil
	} else if accounts != nil {
		return accounts, nil
	} else {
		return nil, nil
	}
}

func (sdb *StorageDB) Create(c context.Context, i interface{}) (interface{}, error) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	sdb.loggingService.WriteLog(c, "Storage: Command Create received...")

	id := uuid.New()

	res, ok := i.(model.UserHTTP)
	if ok {
		res.ID = id
		sdb.Data[id] = res
		return res, nil
	}

	res2, ok := i.(model.Account)
	if ok {

		res2.ID = id
		sdb.Data[id] = res2
		return res2, nil
	}

	return nil, errors.New("invalid data")
}

func (sdb *StorageDB) Update(c context.Context, i interface{}) (interface{}, error) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	sdb.loggingService.WriteLog(c, "Storage: Command Update received...")

	res, ok := i.(model.UserHTTP)
	if ok {
		if _, ok := sdb.Data[res.ID]; !ok { //nolint
			return nil, errors.New("not found in DB")
		}
		sdb.Data[res.ID] = res
		return res, nil
	}

	res2, ok := i.(model.Account)
	if ok {
		if _, ok := sdb.Data[res2.ID]; !ok {
			return nil, errors.New("not found in DB")
		}
		var val model.Account
		if res2.UserID == uuid.Nil {
			val, _ = sdb.Data[res2.ID].(model.Account)
			res2.UserID = val.UserID
		}
		sdb.Data[res2.ID] = res2
		return res2, nil
	}

	return nil, errors.New("not found")
}

func (sdb *StorageDB) Delete(c context.Context, id uuid.UUID) error {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	sdb.loggingService.WriteLog(c, "Storage: Command Delete received...")

	if _, ok := sdb.Data[id]; !ok {
		return errors.New("not found")
	}

	delete(sdb.Data, id)

	return nil
}
