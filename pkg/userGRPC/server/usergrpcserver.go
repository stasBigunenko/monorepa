package usergrpcserver

import (
	"context"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/stasBigunenko/monorepa/model"
	pb "github.com/stasBigunenko/monorepa/pkg/userGRPC/proto"
	"github.com/stasBigunenko/monorepa/service/user"
)

// Server userGRPC

type LoggingService interface {
	WriteLog(ctx context.Context, message string)
}

type UserServerGRPC struct {
	pb.UnimplementedUserGRPCServiceServer

	service        user.User
	loggingService LoggingService
}

func NewUsersGRPCServer(s user.User, loggingService LoggingService) UserServerGRPC {
	return UserServerGRPC{
		service:        s,
		loggingService: loggingService,
	}
}

func (s UserServerGRPC) Get(c context.Context, in *pb.Id) (*pb.User, error) {

	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		log.Info("Cann't receive metada")
	}

	if ccc, ok := md["requestid"]; ok {
		c = context.WithValue(context.Background(), model.ContextKeyRequestID, ccc[0])
	}

	s.loggingService.WriteLog(c, "GRPC Server: Command Get received...")

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to parse uuid in grpc server")
	}

	res, err := s.service.Get(c, id)
	if err != nil {
		return &pb.User{}, status.Error(codes.Internal, "internal problem")
	}

	return &pb.User{
		Id:   res.ID.String(),
		Name: res.Name,
	}, nil
}
func (s UserServerGRPC) GetAllUsers(c context.Context, in *emptypb.Empty) (*pb.AllUsers, error) {
	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		log.Info("Cann't receive metada")
	}

	if ccc, ok := md["requestid"]; ok {
		c = context.WithValue(context.Background(), model.ContextKeyRequestID, ccc[0])
	}

	s.loggingService.WriteLog(c, "GRPC Server: Command GetAllUsers received...")

	users, err := s.service.GetAll(c)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get the list of users")
	}

	pbAllUsers := []*pb.User{}

	for _, val := range users {
		pbAllUsers = append(pbAllUsers, &pb.User{
			Id:   val.ID.String(),
			Name: val.Name,
		})
	}
	return &pb.AllUsers{
		AllUsers: pbAllUsers,
	}, nil
}

func (s UserServerGRPC) Create(c context.Context, in *pb.Name) (*pb.User, error) {

	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		log.Info("Cann't receive metada")
	}

	if ccc, ok := md["requestid"]; ok {
		c = context.WithValue(context.Background(), model.ContextKeyRequestID, ccc[0])
	}

	s.loggingService.WriteLog(c, "GRPC Server: Command Create received...")

	res, err := s.service.Create(c, in.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &pb.User{
		Id:   res.ID.String(),
		Name: res.Name,
	}, nil
}

func (s UserServerGRPC) Update(c context.Context, in *pb.User) (*pb.User, error) {
	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		log.Info("Cann't receive metada")
	}

	if ccc, ok := md["requestid"]; ok {
		c = context.WithValue(context.Background(), model.ContextKeyRequestID, ccc[0])
	}

	s.loggingService.WriteLog(c, "GRPC Server: Command Update received...")

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed to parse uuid")
	}

	m := model.UserHTTP{
		ID:   id,
		Name: in.Name,
	}

	res, err := s.service.Update(c, m)
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed to update user")
	}

	return &pb.User{
		Id:   res.ID.String(),
		Name: res.Name,
	}, nil
}
func (s UserServerGRPC) Delete(c context.Context, in *pb.Id) (*emptypb.Empty, error) {

	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		log.Info("Cann't receive metada")
	}

	if ccc, ok := md["requestid"]; ok {
		c = context.WithValue(context.Background(), model.ContextKeyRequestID, ccc[0])
	}

	s.loggingService.WriteLog(c, "GRPC Server: Command Delete received...")

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to parse uuid")
	}

	err = s.service.Delete(c, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to parse uuid")
	}

	return &emptypb.Empty{}, nil
}
