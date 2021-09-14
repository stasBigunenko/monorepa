package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "monorepa/pkg/grpc/proto"
	"monorepa/pkg/storage"
)

// Server GRPC

type ServerGRPC struct {
	pb.UnimplementedGrpcServiceServer

	storage storage.ItemService
}

func NewGRPC(s storage.ItemService) ServerGRPC {
	return ServerGRPC{
		storage: s,
	}
}

func (s ServerGRPC) GetItems(c context.Context, in *pb.Username) (*pb.Items, error) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	if in.Username == "" {
		return &pb.Items{}, status.Error(codes.InvalidArgument, "invalid username")
	}
	items := []*pb.Obj{}

	data, err := s.storage.GetItems(ctx, in.Username)
	if err != nil {
		return &pb.Items{}, status.Error(codes.Internal, "internal problem")
	}

	for _, val := range data {
		items = append(items, &pb.Obj{
			Id:          val.ID,
			Title:       val.Title,
			Description: val.Description,
		})
	}

	return &pb.Items{
		Items: items,
	}, nil
}
