package main

import (
	"context"
	httphandler "github.com/stasBigunenko/monorepa/pkg/http/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	grpccontroller "github.com/stasBigunenko/monorepa/pkg/grpc/controller"
	pb "github.com/stasBigunenko/monorepa/pkg/grpc/proto"

	"google.golang.org/grpc"
)

type Config struct {
	HTTPAddress string
	JWTAddress  string
	GRPCAddress string
}

func getCfg() Config {
	httpAddr := os.Getenv("HTTP_ADDRESS")
	if httpAddr == "" {
		httpAddr = "127.0.0.1:8080"
	}

	jwtAddr := os.Getenv("JWT_ADDRESS")
	if jwtAddr == "" {
		jwtAddr = "127.0.0.1:8081/get-cert"
	}

	grpcAddr := os.Getenv("GRPC_ADDRESS")
	if grpcAddr == "" {
		grpcAddr = "127.0.0.1:50051"
	}

	return Config{
		HTTPAddress: httpAddr,
		JWTAddress:  jwtAddr,
		GRPCAddress: grpcAddr,
	}
}

func main() {
	cfg := getCfg()

	conn, err := grpc.Dial(cfg.GRPCAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to grpc: %v", err)
	}
	defer conn.Close()

	service := grpccontroller.New(pb.NewGrpcServiceClient(conn))

	h := httphandler.New(service, cfg.JWTAddress)

	srv := http.Server{
		Addr:    cfg.HTTPAddress,
		Handler: h.GetRouter(),
	}

	sigC := make(chan os.Signal, 1)
	defer close(sigC)
	go func() {
		<-sigC
		srv.Shutdown(context.TODO()) // nolint:errcheck
	}()
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("error: http server failed: %s", err)
	}
}
