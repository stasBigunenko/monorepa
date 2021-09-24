package account

import (
	"context"
	"github.com/google/uuid"

	"github.com/stasBigunenko/monorepa/model"
	"github.com/stasBigunenko/monorepa/pkg/storage/newStorage"
)

type LoggingService interface {
	WriteLog(ctx context.Context, message string)
}

type AccService struct {
	storage        newStorage.NewStore
	loggingService LoggingService
}

func NewAccService(s newStorage.NewStore, loggingService LoggingService) *AccService {
	return &AccService{
		storage:        s,
		loggingService: loggingService,
	}
}

func (a *AccService) Get(c context.Context, id uuid.UUID) (model.Account, error) {
	a.loggingService.WriteLog(c, "AccService: Command Get received...")

	val, err := a.storage.Get(c, id)
	if err != nil {
		return model.Account{}, err
	}

	res, ok := val.(model.Account)
	if !ok {
		return model.Account{}, err
	}

	return res, nil

}

func (a *AccService) GetUser(c context.Context, userID uuid.UUID) ([]model.Account, error) {
	a.loggingService.WriteLog(c, "AccService: Command GetUser received...")

	val, err := a.storage.GetUserAccounts(c, userID)
	if err != nil {
		return []model.Account{}, err
	}

	res, ok := val.([]model.Account)
	if !ok {
		return []model.Account{}, err
	}

	return res, nil
}

func (a *AccService) GetAll(c context.Context) ([]model.Account, error) {
	a.loggingService.WriteLog(c, "AccService: Command GetAll received...")

	val, err := a.storage.GetAll(c)
	if err != nil {
		return nil, err
	}

	res, ok := val.([]model.Account)
	if !ok {
		return []model.Account{}, err
	}

	return res, nil
}

func (a *AccService) Create(c context.Context, userID uuid.UUID) (model.Account, error) {
	a.loggingService.WriteLog(c, "AccService: Command Create received...")

	acc := model.Account{
		UserID: userID,
	}

	val, err := a.storage.Create(c, acc)
	if err != nil {
		return model.Account{}, err
	}

	res, ok := val.(model.Account)
	if !ok {
		return model.Account{}, err
	}

	return res, nil
}

func (a *AccService) Update(c context.Context, account model.Account) (model.Account, error) {
	a.loggingService.WriteLog(c, "AccService: Command Update received...")

	val, err := a.storage.Update(c, account)
	if err != nil {
		return model.Account{}, err
	}

	res, ok := val.(model.Account)
	if !ok {
		return model.Account{}, err
	}

	return res, nil
}

func (a *AccService) Delete(c context.Context, id uuid.UUID) error {
	a.loggingService.WriteLog(c, "AccService: Command Delete received...")

	err := a.storage.Delete(c, id)
	if err != nil {
		return err
	}

	return nil
}
