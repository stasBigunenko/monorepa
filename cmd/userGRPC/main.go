package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/stasBigunenko/monorepa/pkg/storage/newStorage"
	pb "github.com/stasBigunenko/monorepa/pkg/userGRPC/proto"
	usergrpcserver "github.com/stasBigunenko/monorepa/pkg/userGRPC/server"
	"github.com/stasBigunenko/monorepa/service/user"
)

type Config struct {
	userGRPCServAddress string
}

func getConfig() Config {
	userGrpcServAddr := os.Getenv("USER_GRPC_SERV_ADDRESS")
	if userGrpcServAddr == "" {
		userGrpcServAddr = "127.0.0.1:50052"
	}

	return Config{
		userGRPCServAddress: userGrpcServAddr,
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

	lis, err := net.Listen("tcp", config.userGRPCServAddress)
	if err != nil {
		log.Fatal("failed to listen: ", err)
	}

	db := newStorage.NewDB()
	dbInt := newStorage.NewStore(db)
	usi := user.NewUsrService(dbInt)

	s := grpc.NewServer()
	pb.RegisterUserGRPCServiceServer(s, usergrpcserver.NewUsersGRPCServer(usi))

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
