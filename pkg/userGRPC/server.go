package userGRPC

import (
	"context"
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
	pb "github.com/stasBigunenko/monorepa/pkg/UserGRPC/proto"
	"github.com/stasBigunenko/monorepa/pkg/storage/UsersDB"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserGRPC struct {
	pb.UnimplementedUserGRPCServiceServer

	storage usersDB.UsersService
}

func NewUserGRPC(s usersDB.UsersService) UserGRPC {
	return UserGRPC{
		storage: s,
	}
}

func (s *UserGRPC) Get(c context.Context, in *pb.Id) (*pb.User, error) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return &pb.User{}, status.Error(codes.InvalidArgument, "invalid ID")
	}

	user, err := s.storage.Get(ctx, id)
	if err != nil {
		return &pb.User{}, status.Error(codes.Internal, "internal problem")
	}

	return &pb.User{
		Id:   in.Id,
		Name: user.Name,
	}, nil
}

func (s *UserGRPC) GetAllUsers(c context.Context, in *emptypb.Empty) (*pb.AllUsers, error) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	users, err := s.storage.GetAllUsers(ctx)
	if err != nil {
		return &pb.AllUsers{}, status.Error(codes.Internal, "storage problem")
	}

	pbUsers := []*pb.User{}

	for _, u := range users {
		id := u.ID.String()
		pbUsers = append(pbUsers, &pb.User{
			Id:   id,
			Name: u.Name,
		})
	}

	return &pb.AllUsers{
		AllUsers: pbUsers,
	}, nil
}

func (s *UserGRPC) Create(c context.Context, in *pb.Name) (*pb.User, error) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	user, err := s.storage.Create(ctx, in.Name)
	if err != nil {
		return &pb.User{}, status.Error(codes.Internal, "storage problem")
	}

	return &pb.User{
		Id:   user.ID.String(),
		Name: user.Name,
	}, nil

}
func (s *UserGRPC) Update(c context.Context, in *pb.User) (*pb.User, error) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid ID")
	}

	user := model.UserHTTP{
		ID:   id,
		Name: in.Name,
	}

	newUser, err := s.storage.Update(ctx, user)
	if err != nil {
		return &pb.User{}, status.Error(codes.Internal, "storage problem")
	}

	return &pb.User{
		Id:   newUser.ID.String(),
		Name: newUser.Name,
	}, nil
}

func (s *UserGRPC) Delete(c context.Context, in *pb.Id) (*emptypb.Empty, error) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid ID")
	}

	err = s.storage.Delete(ctx, id)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "storage problem")
	}

	return &emptypb.Empty{}, nil
}
