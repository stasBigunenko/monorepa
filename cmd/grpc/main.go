package main

import (
	"log"
	"os"

	"github.com/stasBigunenko/monorepa/pkg/grpc/grpcstart"
	"google.golang.org/grpc"
	"net"
	"os/signal"
	"syscall"
)

type Config struct {
	gRPCServAddress string
}

func getConfig() Config {
	grpcServAddr := os.Getenv("GRPC_SERV_ADDRESS")
	if grpcServAddr == "" {
		grpcServAddr = "127.0.0.1:50052"
	}

	return Config{
		gRPCServAddress: grpcServAddr,
	}
}

func main() {
	config := getConfig()

	lis, err := net.Listen("tcp", config.gRPCServAddress)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	srv := grpcstart.GrpcStart()

	sigC := make(chan os.Signal, 1)
	defer close(sigC)
	go func() {
		<-sigC
		srv.GracefulStop()
	}()
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

	// start server
	if err := srv.Serve(lis); err != nil && err != grpc.ErrServerStopped {
		log.Printf("error: grpc server failed: %s", err)
	}
}
