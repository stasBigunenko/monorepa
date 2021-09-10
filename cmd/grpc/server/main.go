package main

import (
	"google.golang.org/grpc"
	"log"
	grpc2 "monorepa/pkg/grpc"
	pb "monorepa/pkg/grpc/proto"
	"monorepa/pkg/storage"
	"net"
)

func main() {
	//Get addr
	addr := ":8080"

	//start listening on tcp
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	d := storage.NewDataProvider()

	data := storage.StorageInterface(d)

	//create GRPC server
	s := grpc.NewServer()
	pb.RegisterGrpcServiceServer(s, grpc2.NewGRPC(data))

	//start server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
