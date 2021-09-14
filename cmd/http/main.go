package main

import (
	"context"
	"log"
	httphandler "monorepa/pkg/http/handler"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	grpccontroller "monorepa/pkg/grpc/controller"
	pb "monorepa/pkg/grpc/proto"

	"google.golang.org/grpc"
)

type Config struct {
	HttpAddress string
	JwtAddress  string
	GrpcAddress string
}

func getCfg() Config {
	httpAddr := os.Getenv("HTTP_ADDRESS")
	if httpAddr == "" {
		httpAddr = "127.0.0.1:8080"
	}

	jwtAddr := os.Getenv("JWT_ADDRESS")
	if jwtAddr == "" {
		jwtAddr = "127.0.0.1:8081"
	}

	grpcAddr := os.Getenv("GRPC_ADDRESS")
	if jwtAddr == "" {
		jwtAddr = "127.0.0.1:50051"
	}

	return Config{
		HttpAddress: httpAddr,
		JwtAddress:  jwtAddr,
		GrpcAddress: grpcAddr,
	}
}

func main() {
	cfg := getCfg()

	conn, err := grpc.Dial(cfg.GrpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to grpc: %v", err)
	}
	defer conn.Close()

	service := grpccontroller.New(pb.NewGrpcServiceClient(conn))

	h := httphandler.New(service, cfg.JwtAddress)

	srv := http.Server{
		Addr:    cfg.HttpAddress,
		Handler: h.GetRouter(),
	}

	sigC := make(chan os.Signal, 1)
	defer close(sigC)
	go func() {
		<-sigC
		srv.Shutdown(context.TODO())
	}()
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("error: http server failed: %s", err)
	}
}
