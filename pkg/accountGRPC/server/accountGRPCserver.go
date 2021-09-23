package accountgrpcserver

import (
	"context"
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
	pb "github.com/stasBigunenko/monorepa/pkg/accountGRPC/proto"
	"github.com/stasBigunenko/monorepa/service/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AccountServerGRPC struct {
	pb.UnimplementedAccountGRPCServiceServer

	service account.AccInterface
}

func NewAccountGRPCServer(s account.AccInterface) AccountServerGRPC {
	return AccountServerGRPC{
		service: s,
	}
}

func (s AccountServerGRPC) GetAccount(c context.Context, in *pb.AccountID) (*pb.Account, error) {

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse uuid in grpc server")
	}

	res, err := s.service.Get(c, id)
	if err != nil {
		return &pb.Account{}, status.Error(codes.Internal, "internal problems")
	}

	return &pb.Account{
		Id:      res.ID.String(),
		UserID:  res.UserID.String(),
		Balance: int32(res.Balance),
	}, nil
}

func (s AccountServerGRPC) GetUserAccounts(c context.Context, in *pb.UserID) (*pb.AllAccounts, error) {

	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse uuid in grpc server")
	}

	users, err := s.service.GetUser(c, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get the list of users")
	}

	all := []*pb.Account{}

	for _, val := range users {
		all = append(all, &pb.Account{
			Id:      val.ID.String(),
			UserID:  val.UserID.String(),
			Balance: int32(val.Balance),
		})
	}
	return &pb.AllAccounts{
		Accounts: all,
	}, nil
}

func (s AccountServerGRPC) GetAllUsers(c context.Context, in *emptypb.Empty) (*pb.AllAccounts, error) {

	users, err := s.service.GetAll(c)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get the list of users")
	}

	all := []*pb.Account{}

	for _, val := range users {
		all = append(all, &pb.Account{
			Id:      val.ID.String(),
			UserID:  val.UserID.String(),
			Balance: int32(val.Balance),
		})
	}
	return &pb.AllAccounts{
		Accounts: all,
	}, nil
}
func (s AccountServerGRPC) CreateAccount(c context.Context, in *pb.UserID) (*pb.Account, error) {

	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse uuid in grpc server")
	}

	res, err := s.service.Create(c, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &pb.Account{
		Id:      res.ID.String(),
		UserID:  res.UserID.String(),
		Balance: int32(res.Balance),
	}, nil
}
func (s AccountServerGRPC) UpdateAccount(c context.Context, in *pb.Account) (*pb.Account, error) {

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to parse uuid")
	}

	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to parse uuid")
	}

	m := model.Account{
		ID:      id,
		UserID:  userID,
		Balance: int(in.Balance),
	}

	res, err := s.service.Update(c, m)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	return &pb.Account{
		Id:      res.ID.String(),
		UserID:  res.UserID.String(),
		Balance: int32(res.Balance),
	}, nil
}
func (s AccountServerGRPC) DeleteAccount(c context.Context, in *pb.AccountID) (*emptypb.Empty, error) {

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to parse uuid")
	}

	err = s.service.Delete(c, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete")
	}

	return &emptypb.Empty{}, nil
}
