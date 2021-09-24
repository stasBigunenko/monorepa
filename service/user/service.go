package user

import (
	"context"
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

	val, err := u.storage.Get(c, id)
	if err != nil {
		return model.UserHTTP{}, err
	}

	res, ok := val.(model.UserHTTP)
	if !ok {
		return model.UserHTTP{}, err
	}

	return res, nil
}

func (u *UsrService) GetAll(c context.Context) ([]model.UserHTTP, error) {
	u.loggingService.WriteLog(c, "User service: Command GetAll received...")

	val, err := u.storage.GetAll(c)
	if err != nil {
		return nil, errors.New("not found")
	}

	res, ok := val.([]model.UserHTTP)
	if !ok {
		return []model.UserHTTP{}, err
	}

	return res, nil
}

func (u *UsrService) Create(c context.Context, name string) (model.UserHTTP, error) {

	u.loggingService.WriteLog(c, "User service: Command Create received...")

	m := model.UserHTTP{ID: uuid.Nil, Name: name}

	val, err := u.storage.Create(c, m)
	if err != nil {
		return model.UserHTTP{}, err
	}

	res, ok := val.(model.UserHTTP)
	if !ok {
		return model.UserHTTP{}, errors.New("invalid data")
	}

	return res, nil
}

func (u *UsrService) Update(c context.Context, user model.UserHTTP) (model.UserHTTP, error) {
	u.loggingService.WriteLog(c, "User service: Command Update received...")

	val, err := u.storage.Update(c, user)
	if err != nil {
		return model.UserHTTP{}, err
	}

	res, ok := val.(model.UserHTTP)
	if !ok {
		return model.UserHTTP{}, errors.New("invalid data")
	}

	return res, nil
}

func (u *UsrService) Delete(c context.Context, id uuid.UUID) error {
	u.loggingService.WriteLog(c, "User service: Command Delete received...")

	err := u.storage.Delete(c, id)
	if err != nil {
		return err
	}

	return nil
}
