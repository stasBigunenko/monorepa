package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"monorepa/pkg/auth/middleware"
	"monorepa/pkg/auth/routes"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	ctx    context.Context
	config struct { // can be used as general model
		Host string
		Port string
	}
	router *mux.Router
}

func New(ctx context.Context) *Server {
	return &Server{
		router: mux.NewRouter(),
		ctx:    ctx,
	}
}

func (s *Server) ServerAddrConfig() error {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		return errors.New("wrong port")
	}

	s.config.Host = os.Getenv("SERVER_HOST")
	s.config.Port = port
	return nil
}

/**
 * all routes of http server
 * also contain connection to graphql server
 */

func (s *Server) GetRouters() {
	itemsHandler := routes.New(s.ctx, s.router)
	itemsHandler.HandlerItems()
}

func (s *Server) getHTTPAddress() string {
	return fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
}

func (s *Server) Start(ctx context.Context) error {
	// init server
	server := &http.Server{
		Addr:         s.getHTTPAddress(),
		Handler:      middleware.JSONRespHeaders(s.router),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	cancelTimeout, errTimeout := strconv.Atoi(os.Getenv("Server_Cancel_Timeout"))
	if errTimeout != nil {
		return errTimeout
	}

	// run server
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	fmt.Println("Server is running on: " + s.getHTTPAddress())

	<-ctx.Done() // wait end of work

	/**
	 * start graceful shutdown
	 */

	log.Println("Server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), time.Duration(cancelTimeout)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutDown); err != nil {
		return err
	}

	log.Printf("Server exited properly")
	return nil
}
