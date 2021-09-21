package account

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
	"github.com/stasBigunenko/monorepa/pkg/storage/newStorage"
)

type AccService struct {
	Acc     model.Account
	storage newStorage.NewStore
}

func NewAccService(s newStorage.NewStore) *AccService {
	return &AccService{
		storage: s,
	}
}

func (a *AccService) Get(_ context.Context, id uuid.UUID) (model.Account, error) {

	res, err := a.storage.Get(context.Background(), id)
	if err != nil && (err == nil || err != nil) {
		return model.Account{}, errors.New("couldn't get account")
	}

	err = json.Unmarshal(res, &a.Acc)
	if err != nil {
		return model.Account{}, errors.New("couldn't unmarshal data")
	}

	return a.Acc, nil

}

func (a *AccService) GetUser(_ context.Context, userID uuid.UUID) ([]model.Account, error) {
	res, err := a.storage.GetUser(context.Background(), userID)
	if err != nil || (res == nil && err == nil) {
		return []model.Account{}, errors.New("storage problem")
	}

	var ac []model.Account

	for _, val := range res {
		var a model.Account
		err := json.Unmarshal(val, &a)
		if err != nil {
			return []model.Account{}, errors.New("couldn't unmarshal data")
		}
		ac = append(ac, a)
	}

	return ac, nil
}

func (a *AccService) GetAll(_ context.Context) ([]model.Account, error) {

	res, err := a.storage.GetAll(context.Background())
	if err != nil {
		return nil, errors.New("storage problem")
	}

	ac := []model.Account{}

	for _, val := range res {
		err := json.Unmarshal(val, &a.Acc)
		if err != nil {
			return nil, errors.New("couldn't unmarshal data")
		}

		ac = append(ac, a.Acc)
	}
	return ac, nil
}

func (a *AccService) Create(_ context.Context, b uuid.UUID) (model.Account, error) {

	br, err := json.Marshal(b)
	if err != nil {
		return model.Account{}, errors.New("couldn't unmarshal data")
	}

	_, id, err := a.storage.Create(context.Background(), br)
	if err != nil {
		return model.Account{}, errors.New("storage problem")
	}

	a.Acc.ID = id
	a.Acc.UserID = b
	a.Acc.Balance = 0

	return a.Acc, nil
}

func (a *AccService) Update(_ context.Context, account model.Account) (model.Account, error) {
	id := account.ID

	bt, err := json.Marshal(account)
	if err != nil {
		return model.Account{}, errors.New("couldn't marshal data")
	}

	res, err := a.storage.Update(context.Background(), id, bt)
	if err != nil {
		return model.Account{}, errors.New("not found")
	}

	err = json.Unmarshal(res, &a.Acc)
	if err != nil {
		return model.Account{}, errors.New("couldn't unmarshal data")
	}

	return a.Acc, nil
}

func (a *AccService) Delete(_ context.Context, id uuid.UUID) error {
	b, err := a.storage.Delete(context.Background(), id)
	if err != nil || !b {
		return errors.New("not found")
	}
	return nil
}
