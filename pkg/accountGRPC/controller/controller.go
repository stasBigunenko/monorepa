package accountgrpccontroller

import (
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
	pb "github.com/stasBigunenko/monorepa/pkg/accountGRPC/proto"
)

type AccountGRPCСontroller struct {
	client pb.AccountGRPCServiceClient
}

func New(cli pb.AccountGRPCServiceClient) AccountGRPCСontroller {
	return AccountGRPCСontroller{
		client: cli,
	}
}

func (s AccountGRPCСontroller) CreateAccount(userID uuid.UUID) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (s AccountGRPCСontroller) GetAccount(id uuid.UUID) (model.Account, error) {
	return model.Account{}, nil
}

func (s AccountGRPCСontroller) GetUserAccounts(userID uuid.UUID) ([]model.Account, error) {
	return nil, nil
}

func (s AccountGRPCСontroller) GetAllAccounts() ([]model.Account, error) {
	return nil, nil
}

func (s AccountGRPCСontroller) UpdateAccount(model.Account) error {
	return nil
}

func (s AccountGRPCСontroller) DeleteAccount(id uuid.UUID) error {
	return nil
}
