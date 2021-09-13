package main

import (
	"context"
	"log"
	"monorepa/pkg/auth"
	"os"
	"os/signal"
)

func init() {
	// serever
	os.Setenv("SERVER_HOST", "")
	os.Setenv("SERVER_PORT", "8080")

	// shutdown
	os.Setenv("Server_Cancel_Timeout", "5")

	//JWT config
	os.Setenv("TOKEN_EXPIRE", "10") // minutes

	// certificates
	os.Setenv("CERT_VERSION", "1")
	os.Setenv("CERT_PATH", "./pkg/storage/certificates")
}

func main() {
	// init data
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx := context.Background()

	// init server and config

	// server := auth.New(ctx)
	server := auth.MockServerActions{}
	if err := server.ServerAddrConfig(); err != nil {
		log.Fatal("Can`t get config of server: ", err)
	}

	// add all routers endpoints
	server.GetRouters()

	// create shutdown
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	if err := server.Start(ctx); err != nil {
		log.Fatal("Problems with server run: ", err)
	}

}
