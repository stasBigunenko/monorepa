package grpcstart

import (
	grpcservice "github.com/stasBigunenko/monorepa/pkg/grpc"
	pb "github.com/stasBigunenko/monorepa/pkg/grpc/proto"

	"github.com/stasBigunenko/monorepa/pkg/storage"
	"google.golang.org/grpc"
)

func GrpcStart() *grpc.Server {
	d := storage.NewStorage()
	data := storage.ItemService(d)

	s := grpc.NewServer()
	pb.RegisterGrpcServiceServer(s, grpcservice.NewGRPC(data))

	return s
}
