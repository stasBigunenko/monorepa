package user

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"

	"github.com/stasBigunenko/monorepa/model"
	"github.com/stasBigunenko/monorepa/pkg/storage/newStorage"
)

type LoggingService interface {
	WriteLog(ctx context.Context, message string)
}

type UsrService struct {
	storage        newStorage.NewStore
	loggingService LoggingService
}

func NewUsrService(s newStorage.NewStore, loggingService LoggingService) *UsrService {
	return &UsrService{
		storage:        s,
		loggingService: loggingService,
	}
}

func (u *UsrService) Get(c context.Context, id uuid.UUID) (model.UserHTTP, error) {
	u.loggingService.WriteLog(c, "User service: Command Get received...")

	res, err := u.storage.Get(c, id)
	if err != nil {
		return model.UserHTTP{}, err
	}

	if res == nil {
		return model.UserHTTP{}, errors.New("not found")
	}

	var m model.UserHTTP

	err = json.Unmarshal(res, &m)
	if err != nil {
		return model.UserHTTP{}, errors.New("unmarshal problem")
	}

	return m, nil
}

func (u *UsrService) GetAll(c context.Context) ([]model.UserHTTP, error) {
	u.loggingService.WriteLog(c, "User service: Command GetAll received...")

	res, err := u.storage.GetAll(c)
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

func (u *UsrService) Create(c context.Context, name string) (model.UserHTTP, error) {
	u.loggingService.WriteLog(c, "User service: Command Create received...")

	m := model.UserHTTP{
		Name: name,
	}

	bt, err := json.Marshal(m)
	if err != nil {
		return model.UserHTTP{}, errors.New("marshal problem")
	}

	res, id, err := u.storage.Create(c, bt)
	if err != nil {
		return model.UserHTTP{}, err
	}

	err = json.Unmarshal(res, &m)
	if err != nil {
		return model.UserHTTP{}, errors.New("unmarshal problem")
	}

	m.ID = id

	return m, nil
}

func (u *UsrService) Update(c context.Context, user model.UserHTTP) (model.UserHTTP, error) {
	u.loggingService.WriteLog(c, "User service: Command Update received...")

	id := user.ID

	bt, err := json.Marshal(user)
	if err != nil {
		return model.UserHTTP{}, errors.New("marshal problem")
	}

	res, err := u.storage.Update(c, id, bt)
	if err != nil {
		return model.UserHTTP{}, err
	}

	if res == nil && err == nil {
		return model.UserHTTP{}, errors.New("not found")
	}

	var m model.UserHTTP

	err = json.Unmarshal(res, &m)
	if err != nil {
		return model.UserHTTP{}, errors.New("unmarshal problem")
	}

	return m, nil
}

func (u *UsrService) Delete(c context.Context, id uuid.UUID) error {
	u.loggingService.WriteLog(c, "User service: Command Delete received...")

	b, err := u.storage.Delete(c, id)
	if err != nil {
		return err
	}

	if !b {
		return errors.New("not found")
	}

	return nil
}
