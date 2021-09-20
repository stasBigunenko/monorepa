package accountsDB

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
	"sync"
)

type AccountStorage struct {
	MapAccount map[uuid.UUID]model.Account
	mu         sync.Mutex
}

func NewAccountStorage() *AccountStorage {
	as := AccountStorage{}
	as.MapAccount = make(map[uuid.UUID]model.Account)
	return &as
}

func (as *AccountStorage) GetAccount(_ context.Context, id uuid.UUID) (model.Account, error) {
	as.mu.Lock()
	val, ok := as.MapAccount[id]
	as.mu.Unlock()
	if !ok {
		return model.Account{}, fmt.Errorf("account with id %d not found", id)
	}

	return val, nil
}

func (as *AccountStorage) GetUserAccounts(_ context.Context, user_id uuid.UUID) ([]model.Account, error) {
	as.mu.Lock()
	accounts := []model.Account{}

	for _, account := range as.MapAccount {
		if user_id == account.UserID {
			accounts = append(accounts, account)
		}
	}

	as.mu.Unlock()
	return accounts, nil
}

func (as *AccountStorage) GetAllUsers(_ context.Context) ([]model.Account, error) {
	as.mu.Lock()

	var accounts []model.Account

	for _, account := range as.MapAccount {
		accounts = append(accounts, account)
	}

	as.mu.Unlock()

	return accounts, nil
}

func (as *AccountStorage) CreateAccount(_ context.Context, userID uuid.UUID) (model.Account, error) {
	as.mu.Lock()

	id := uuid.New()

	as.mu.Unlock()

	account := model.Account{
		ID:      id,
		UserID:  userID,
		Balance: 0,
	}

	as.MapAccount[id] = account

	return account, nil
}

func (as *AccountStorage) UpdateAccount(_ context.Context, acc model.Account) (model.Account, error) {
	as.mu.Lock()
	if _, ok := as.MapAccount[acc.ID]; !ok {
		return model.Account{}, fmt.Errorf("account with id %d not found", acc.ID)
	}

	as.MapAccount[acc.ID] = acc
	as.mu.Unlock()
	return acc, nil
}

func (as *AccountStorage) DeleteAccount(_ context.Context, id uuid.UUID) error {

	if _, ok := as.MapAccount[id]; !ok {
		return fmt.Errorf("account with id %d not found", id)
	}

	delete(as.MapAccount, id)

	return nil
}
