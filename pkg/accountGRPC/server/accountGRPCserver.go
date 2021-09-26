package accountgrpcserver

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/stasBigunenko/monorepa/model"
	pb "github.com/stasBigunenko/monorepa/pkg/accountGRPC/proto"
	"github.com/stasBigunenko/monorepa/service/account"
)

type LoggingService interface {
	WriteLog(ctx context.Context, message string)
}

type AccountServerGRPC struct {
	pb.UnimplementedAccountGRPCServiceServer

	service        account.AccInterface
	loggingService LoggingService
}

func NewAccountGRPCServer(s account.AccInterface, loggingService LoggingService) AccountServerGRPC {
	return AccountServerGRPC{
		service:        s,
		loggingService: loggingService,
	}
}

func (s AccountServerGRPC) GetAccount(c context.Context, in *pb.AccountID) (*pb.Account, error) {

	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		log.Info("Cann't receive metada")
	}

	ccc, ok := md["requestid"]
	if !ok {
		log.Info("cann't receive request id")
	}

	ctx := context.WithValue(context.Background(), model.ContextKeyRequestID, ccc[0])

	s.loggingService.WriteLog(ctx, "GRPC Server: Command GetAccount received...")

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse uuid in grpc server")
	}

	res, err := s.service.Get(ctx, id)
	if err != nil {
		return &pb.Account{}, status.Error(codes.DataLoss, "internal problems")
	}

	return &pb.Account{
		Id:      res.ID.String(),
		UserID:  res.UserID.String(),
		Balance: int32(res.Balance),
	}, nil
}

func (s AccountServerGRPC) GetUserAccounts(c context.Context, in *pb.UserID) (*pb.AllAccounts, error) {

	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		log.Info("Cann't receive metada")
	}

	ccc, ok := md["requestid"]
	if !ok {
		log.Info("cann't receive request id")
	}

	ctx := context.WithValue(context.Background(), model.ContextKeyRequestID, ccc[0])

	s.loggingService.WriteLog(ctx, "GRPC Server: Command GetUserAccounts received...")

	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse uuid in grpc server")
	}

	users, err := s.service.GetUser(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.DataLoss, "failed to get the list of users")
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

	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		log.Info("Cann't receive metada")
	}

	ccc, ok := md["requestid"]
	if !ok {
		log.Info("cann't receive request id")
	}

	ctx := context.WithValue(context.Background(), model.ContextKeyRequestID, ccc[0])

	s.loggingService.WriteLog(ctx, "GRPC Server: Command GetAllUsers received...")

	users, err := s.service.GetAll(ctx)
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

	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		log.Info("Cann't receive metada")
	}

	ccc, ok := md["requestid"]
	if !ok {
		log.Info("cann't receive request id")
	}

	ctx := context.WithValue(context.Background(), model.ContextKeyRequestID, ccc[0])

	s.loggingService.WriteLog(ctx, "GRPC Server: Command CreateAccount received...")

	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse uuid in grpc server")
	}

	res, err := s.service.Create(ctx, userID)
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

	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		log.Info("Cann't receive metada")
	}

	ccc, ok := md["requestid"]
	if !ok {
		log.Info("cann't receive request id")
	}

	ctx := context.WithValue(context.Background(), model.ContextKeyRequestID, ccc[0])

	s.loggingService.WriteLog(ctx, "GRPC Server: Command UpdateAccount received...")

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

	res, err := s.service.Update(ctx, m)
	if err != nil {
		return nil, status.Error(codes.DataLoss, "failed to update user")
	}

	return &pb.Account{
		Id:      res.ID.String(),
		UserID:  res.UserID.String(),
		Balance: int32(res.Balance),
	}, nil
}
func (s AccountServerGRPC) DeleteAccount(c context.Context, in *pb.AccountID) (*emptypb.Empty, error) {

	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		log.Info("Cann't receive metada")
	}

	ccc, ok := md["requestid"]
	if !ok {
		log.Info("cann't receive request id")
	}

	ctx := context.WithValue(context.Background(), model.ContextKeyRequestID, ccc[0])

	s.loggingService.WriteLog(ctx, "GRPC Server: Command DeleteAccount received...")

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to parse uuid")
	}

	err = s.service.Delete(ctx, id)
	if err != nil {
		return nil, status.Error(codes.DataLoss, "failed to delete")
	}

	return &emptypb.Empty{}, nil
}
