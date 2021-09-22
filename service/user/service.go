package user

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
	"github.com/stasBigunenko/monorepa/pkg/storage/newStorage"
)

type UsrService struct {
	storage newStorage.NewStore
}

func NewUsrService(s newStorage.NewStore) *UsrService {
	return &UsrService{
		storage: s,
	}
}

func (u *UsrService) Get(_ context.Context, id uuid.UUID) (model.UserHTTP, error) {

	res, err := u.storage.Get(context.Background(), id)
	if (res == nil && err == nil) || err != nil {
		return model.UserHTTP{}, errors.New("not found")
	}

	var m model.UserHTTP

	err = json.Unmarshal(res, &m)
	if err != nil {
		return model.UserHTTP{}, errors.New("unmarshal problem")
	}

	return m, nil
}

func (u *UsrService) GetAll(_ context.Context) ([]model.UserHTTP, error) {

	res, err := u.storage.GetAll(context.Background())
	if err != nil {
		return nil, errors.New("not found")
	}

	ac := []model.UserHTTP{}

	var m model.UserHTTP

	for _, val := range res {
		err := json.Unmarshal(val, &m)
		if err != nil {
			return nil, errors.New("unmarshal problem")
		}

		ac = append(ac, m)
	}
	return ac, nil
}

func (u *UsrService) Create(_ context.Context, name string) (model.UserHTTP, error) {

	if len(name) <= 2 {
		return model.UserHTTP{}, errors.New("invalid username")
	}

	m := model.UserHTTP{
		Name: name,
	}

	bt, err := json.Marshal(m)
	if err != nil {
		return model.UserHTTP{}, errors.New("marshal problem")
	}

	res, id, err := u.storage.Create(context.Background(), bt)
	if err != nil {
		return model.UserHTTP{}, errors.New("storage problem")
	}

	err = json.Unmarshal(res, &m)
	if err != nil {
		return model.UserHTTP{}, errors.New("unmarshal problem")
	}

	m.ID = id

	return m, nil
}

func (u *UsrService) Update(_ context.Context, user model.UserHTTP) (model.UserHTTP, error) {
	id := user.ID

	bt, err := json.Marshal(user)
	if err != nil {
		return model.UserHTTP{}, errors.New("marshal problem")
	}

	res, err := u.storage.Update(context.Background(), id, bt)
	if err != nil || (res == nil && err == nil) {
		return model.UserHTTP{}, errors.New("not found")
	}

	var m model.UserHTTP

	err = json.Unmarshal(res, &m)
	if err != nil {
		return model.UserHTTP{}, errors.New("unmarshal problem")
	}

	return m, nil
}

func (u *UsrService) Delete(_ context.Context, id uuid.UUID) error {
	b, err := u.storage.Delete(context.Background(), id)
	if err != nil || !b {
		return errors.New("not found")
	}

	return nil
}
