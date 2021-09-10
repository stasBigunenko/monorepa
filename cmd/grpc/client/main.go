package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	pb "monorepa/pkg/grpc/proto"
	my "monorepa/service/items"
)

func main() {
	//Get addr
	addr := ":8080"

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to grpc: %v", err)
	}
	defer conn.Close()

	my.New(pb.NewGrpcServiceClient(conn))
	fmt.Println("hello")
}
