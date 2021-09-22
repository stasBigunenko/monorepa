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
	storage newStorage.NewStore
}

func NewAccService(s newStorage.NewStore) *AccService {
	return &AccService{
		storage: s,
	}
}

func (a *AccService) Get(c context.Context, id uuid.UUID) (model.Account, error) {

	res, err := a.storage.Get(c, id)
	if err != nil {
		return model.Account{}, err
	}

	if res == nil {
		return model.Account{}, errors.New("couldn't get account")
	}

	var acc model.Account

	err = json.Unmarshal(res, &acc)
	if err != nil {
		return model.Account{}, errors.New("couldn't unmarshal data")
	}

	return acc, nil

}

func (a *AccService) GetUser(c context.Context, userID uuid.UUID) ([]model.Account, error) {
	res, err := a.storage.GetUser(c, userID)
	if err != nil {
		return []model.Account{}, err
	}

	if res == nil {
		return []model.Account{}, errors.New("storage problem")
	}

	var acs []model.Account

	for _, val := range res {
		var acc model.Account
		err := json.Unmarshal(val, &acc)
		if err != nil {
			return []model.Account{}, errors.New("couldn't unmarshal data")
		}
		acs = append(acs, acc)
	}

	return acs, nil
}

func (a *AccService) GetAll(c context.Context) ([]model.Account, error) {

	res, err := a.storage.GetAll(c)
	if err != nil {
		return nil, err
	}

	acs := []model.Account{}

	for _, val := range res {
		var acc model.Account
		err := json.Unmarshal(val, &acc)
		if err != nil {
			return nil, errors.New("couldn't unmarshal data")
		}

		acs = append(acs, acc)
	}
	return acs, nil
}

func (a *AccService) Create(c context.Context, userID uuid.UUID) (model.Account, error) {

	acc := model.Account{UserID: userID, Balance: 0}

	jUsrID, err := json.Marshal(acc)
	if err != nil {
		return model.Account{}, errors.New("couldn't marshal data")
	}

	res, id, err := a.storage.Create(c, jUsrID)
	if err != nil {
		return model.Account{}, err
	}

	err = json.Unmarshal(res, &acc)
	if err != nil {
		return model.Account{}, errors.New("couldn't unmarshal data")
	}

	acc.ID = id

	return acc, nil
}

func (a *AccService) Update(c context.Context, account model.Account) (model.Account, error) {
	id := account.ID

	j, err := json.Marshal(account)
	if err != nil {
		return model.Account{}, errors.New("couldn't marshal data")
	}

	res, err := a.storage.Update(c, id, j)
	if err != nil {
		return model.Account{}, err
	}

	var acc model.Account

	err = json.Unmarshal(res, &acc)
	if err != nil {
		return model.Account{}, errors.New("couldn't unmarshal data")
	}

	return acc, nil
}

func (a *AccService) Delete(c context.Context, id uuid.UUID) error {
	b, err := a.storage.Delete(c, id)
	if err != nil || !b {
		return err
	}

	if !b {
		return errors.New("not found")
	}

	return nil
}