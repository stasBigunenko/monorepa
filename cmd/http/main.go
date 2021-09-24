package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	accountscontroller "github.com/stasBigunenko/monorepa/pkg/accountGRPC/controller"
	pbaccounts "github.com/stasBigunenko/monorepa/pkg/accountGRPC/proto"
	httphandler "github.com/stasBigunenko/monorepa/pkg/http/handler"
	userscontroller "github.com/stasBigunenko/monorepa/pkg/userGRPC/controller"
	pbusers "github.com/stasBigunenko/monorepa/pkg/userGRPC/proto"
	loggingservice "github.com/stasBigunenko/monorepa/service/loggingService"
)

type Config struct {
	HTTPAddress        string
	JWTAddress         string
	GRPCAccountAddress string
	GRPCUserAddress    string
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

	grpcAccAddr := os.Getenv("GRPC_ACCOUNTS_ADDRESS")
	if grpcAccAddr == "" {
		grpcAccAddr = "127.0.0.1:50053"
	}

	grpcUserAddr := os.Getenv("GRPC_USERS_ADDRESS")
	if grpcUserAddr == "" {
		grpcUserAddr = "127.0.0.1:50052"
	}

	return Config{
		HTTPAddress:        httpAddr,
		JWTAddress:         jwtAddr,
		GRPCAccountAddress: grpcAccAddr,
		GRPCUserAddress:    grpcUserAddr,
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
	cfg := getCfg()

	connAcc, err := grpc.Dial(cfg.GRPCAccountAddress, grpc.WithInsecure())
	if err != nil {
		log.Info("did not connect to grpc: ", err)
		return
	}
	defer connAcc.Close()

	connUser, err := grpc.Dial(cfg.GRPCUserAddress, grpc.WithInsecure())
	if err != nil {
		log.Info("did not connect to grpc: ", err)
		return
	}
	defer connUser.Close()

	loggingService := loggingservice.New()
	userService := userscontroller.New(pbusers.NewUserGRPCServiceClient(connUser), loggingService)
	accountService := accountscontroller.New(pbaccounts.NewAccountGRPCServiceClient(connAcc), loggingService)

	h := httphandler.New(accountService, userService, loggingService, cfg.JWTAddress)

	srv := http.Server{
		Addr:    cfg.HTTPAddress,
		Handler: h.AccessControlMiddleware(h.GetRouter()),
	}

	sigC := make(chan os.Signal, 1)
	defer close(sigC)
	go func() {
		<-sigC
		srv.Shutdown(context.TODO()) // nolint:errcheck
	}()
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("error: http server failed: ", err)
	}
}
