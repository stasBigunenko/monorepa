package accountgrpccontroller

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	customerrors "github.com/stasBigunenko/monorepa/customErrors"
	"github.com/stasBigunenko/monorepa/model"
	pb "github.com/stasBigunenko/monorepa/pkg/accountGRPC/proto"
)

type LoggingService interface {
	WriteLog(ctx context.Context, message string)
}

type AccountGRPCСontroller struct {
	client         pb.AccountGRPCServiceClient
	loggingService LoggingService
}

func New(cli pb.AccountGRPCServiceClient, loggingService LoggingService) *AccountGRPCСontroller {
	return &AccountGRPCСontroller{
		client:         cli,
		loggingService: loggingService,
	}
}

func (s AccountGRPCСontroller) formatError(err error, message string) error {
	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("%s, failed to parse status or not a grpc error type: %w", message, err)
	}

	switch st.Code() {
	case codes.NotFound:
		return fmt.Errorf("%s: %w", message, customerrors.NotFound)
	case codes.AlreadyExists:
		return fmt.Errorf("%s: %w", message, customerrors.AlreadyExists)
	case codes.DeadlineExceeded:
		return fmt.Errorf("%s: %w", message, customerrors.DeadlineExceeded)
	}

	return fmt.Errorf("%s: %s", message, err.Error())
}

func (s AccountGRPCСontroller) CreateAccount(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	s.loggingService.WriteLog(ctx, "GRPC Client: Command CreateAccount received...")
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

func (s AccountGRPCСontroller) GetAccount(ctx context.Context, id uuid.UUID) (model.Account, error) {
	s.loggingService.WriteLog(ctx, "GRPC Client: Command GetAccount received...")
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
		Balance: int(resp.Balance),
	}, nil
}

func (s AccountGRPCСontroller) GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]model.Account, error) {
	s.loggingService.WriteLog(ctx, "GRPC Client: Command GetUserAccounts received...")
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
			Balance: int(account.Balance),
		})
	}

	return accounts, nil
}

func (s AccountGRPCСontroller) GetAllAccounts(ctx context.Context) ([]model.Account, error) {
	s.loggingService.WriteLog(ctx, "GRPC Client: Command GetAllAccounts received...")
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
			Balance: int(account.Balance),
		})
	}

	return accounts, nil
}

func (s AccountGRPCСontroller) UpdateAccount(ctx context.Context, account model.Account) error {
	s.loggingService.WriteLog(ctx, "GRPC Client: Command UpdateAccount received...")
	_, err := s.client.UpdateAccount(context.Background(), &pb.Account{
		Id:      account.ID.String(),
		UserID:  account.UserID.String(),
		Balance: int32(account.Balance),
	})

	if err != nil {
		return s.formatError(err, "failed to update account")
	}

	return nil
}

func (s AccountGRPCСontroller) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	s.loggingService.WriteLog(ctx, "GRPC Client: Command DeleteAccount received...")
	_, err := s.client.DeleteAccount(context.Background(), &pb.AccountID{
		Id: id.String(),
	})

	if err != nil {
		return s.formatError(err, "failed to delete account")
	}

	return nil
}
