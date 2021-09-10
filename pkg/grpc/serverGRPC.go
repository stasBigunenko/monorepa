package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "monorepa/pkg/grpc/proto"
	"monorepa/pkg/storage"
)

type ServerGRPC struct {
	pb.UnimplementedGrpcServiceServer

	Data storage.StorageInterface
}

func NewGRPC(data storage.StorageInterface) *ServerGRPC {
	return &ServerGRPC{
		Data: data,
	}
}

func (s *ServerGRPC) GetItems(_ context.Context, in *pb.Username) (*pb.Items, error) {

	if in.Username == "" {
		return &pb.Items{}, status.Error(codes.InvalidArgument, "invalid username")
	}
	items := []*pb.Obj{}

	data, err := s.Data.GetItems(in.Username)
	if err != nil {
		return &pb.Items{}, status.Error(codes.Internal, "internal problem")
	}

	for _, val := range data {
		items = append(items, &pb.Obj{
			Id:          val.Id,
			Title:       val.Title,
			Description: val.Description,
		})
	}

	return &pb.Items{
		Items: items,
	}, nil
}
