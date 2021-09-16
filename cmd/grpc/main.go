package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"net"
	"os/signal"
	"syscall"

	"github.com/stasBigunenko/monorepa/pkg/grpc/grpcstart"
	"google.golang.org/grpc"
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

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
}

func main() {
	config := getConfig()

	lis, err := net.Listen("tcp", config.gRPCServAddress)
	if err != nil {
		log.Fatal("failed to listen: %s", err)
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
		log.Error("error: grpc server failed: %s", err)
	}
}
