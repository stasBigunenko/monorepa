package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/stasBigunenko/monorepa/pkg/accountGRPC/proto"
	accountgrpcserver "github.com/stasBigunenko/monorepa/pkg/accountGRPC/server"
	"github.com/stasBigunenko/monorepa/pkg/storage/newStorage"
	"github.com/stasBigunenko/monorepa/service/account"
)

type Config struct {
	accountGRPCServAddress string
}

func getConfig() Config {
	accountGrpcServAddr := os.Getenv("ACCOUNT_GRPC_SERV_ADDRESS")
	if accountGrpcServAddr == "" {
		accountGrpcServAddr = "127.0.0.1:50053"
	}

	return Config{
		accountGRPCServAddress: accountGrpcServAddr,
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

	lis, err := net.Listen("tcp", config.accountGRPCServAddress)
	if err != nil {
		log.Fatal("failed to listen: ", err)
	}

	db := newStorage.NewDB()
	dbInt := newStorage.NewStore(db)
	asi := account.NewAccService(dbInt)

	s := grpc.NewServer()
	pb.RegisterAccountGRPCServiceServer(s, accountgrpcserver.NewAccountGRPCServer(asi))

	sigC := make(chan os.Signal, 1)
	defer close(sigC)
	go func() {
		<-sigC
		s.GracefulStop()
	}()
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

	// start server
	if err := s.Serve(lis); err != nil && err != grpc.ErrServerStopped {
		log.Error("error: grpc server failed: ", err)
	}
}
