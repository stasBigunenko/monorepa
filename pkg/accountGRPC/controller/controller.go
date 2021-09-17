package accountgrpccontroller

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	customerrors "github.com/stasBigunenko/monorepa/errors"
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

func (s AccountGRPCСontroller) formatError(err error, message string) error {
	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("%s, failed to parse status or not a grpc error type: %w", message, err)
	}

	if st.Code() == codes.NotFound {
		return fmt.Errorf("%s: %w", message, customerrors.NotFound)
	}

	if st.Code() == codes.AlreadyExists {
		return fmt.Errorf("%s: %w", message, customerrors.AlreadyExists)
	}

	if st.Code() == codes.DeadlineExceeded {
		return fmt.Errorf("%s: %w", message, customerrors.DeadlineExceeded)
	}

	return fmt.Errorf("%s: %s", message, err.Error())
}

func (s AccountGRPCСontroller) CreateAccount(userID uuid.UUID) (uuid.UUID, error) {
	resp, err := s.client.CreateAccount(context.Background(), &pb.UserID{
		UserID: userID.String(),
	})

	if err != nil {
		return uuid.Nil, s.formatError(err, "failed to create account")
	}

	id, err := uuid.Parse(resp.Id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse user ID: %s, %w", err.Error(), customerrors.ParseError)
	}

	return id, nil
}

func (s AccountGRPCСontroller) GetAccount(id uuid.UUID) (model.Account, error) {
	resp, err := s.client.GetAccount(context.Background(), &pb.AccountID{
		Id: id.String(),
	})

	if err != nil {
		return model.Account{}, s.formatError(err, "failed to get account")
	}

	accountID, err := uuid.Parse(resp.Id)
	if err != nil {
		return model.Account{}, fmt.Errorf("failed to parse account ID: %s, %w", err.Error(), customerrors.ParseError)
	}

	userID, err := uuid.Parse(resp.UserID)
	if err != nil {
		return model.Account{}, fmt.Errorf("failed to parse account ID: %s, %w", err.Error(), customerrors.ParseError)
	}

	return model.Account{
		ID:      accountID,
		UserID:  userID,
		Balance: resp.Balance,
	}, nil
}

func (s AccountGRPCСontroller) GetUserAccounts(userID uuid.UUID) ([]model.Account, error) {
	resp, err := s.client.GetUserAccounts(context.Background(), &pb.UserID{
		UserID: userID.String(),
	})

	if err != nil {
		return nil, s.formatError(err, "failed to get all user accounts")
	}

	accounts := []model.Account{}
	for _, account := range resp.Accounts {
		accountID, err := uuid.Parse(account.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to parse account ID: %s, %w", err.Error(), customerrors.ParseError)
		}

		userID, err := uuid.Parse(account.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse account ID: %s, %w", err.Error(), customerrors.ParseError)
		}

		accounts = append(accounts, model.Account{
			ID:      accountID,
			UserID:  userID,
			Balance: resp.Balance,
		})
	}

	return accounts, nil
}

func (s AccountGRPCСontroller) GetAllAccounts() ([]model.Account, error) {
	resp, err := s.client.GetAllUsers(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, s.formatError(err, "failed to get all accounts")
	}

	accounts := []model.Account{}
	for _, account := range resp.Accounts {
		accountID, err := uuid.Parse(account.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to parse account ID: %s, %w", err.Error(), customerrors.ParseError)
		}

		userID, err := uuid.Parse(account.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse account ID: %s, %w", err.Error(), customerrors.ParseError)
		}

		accounts = append(accounts, model.Account{
			ID:      accountID,
			UserID:  userID,
			Balance: resp.Balance,
		})
	}

	return accounts, nil
}

func (s AccountGRPCСontroller) UpdateAccount(account model.Account) error {
	_, err := s.client.UpdateAccount(context.Background(), &pb.Account{
		Id:      account.ID.String(),
		UserID:  account.UserID.String(),
		Balance: account.Balance,
	})

	if err != nil {
		return s.formatError(err, "failed to update account")
	}

	return nil
}

func (s AccountGRPCСontroller) DeleteAccount(id uuid.UUID) error {
	_, err := s.client.DeleteAccount(context.Background(), &pb.AccountID{
		Id: id.String(),
	})

	if err != nil {
		return s.formatError(err, "failed to delete account")
	}

	return nil
}
