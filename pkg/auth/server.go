package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/stasBigunenko/monorepa/pkg/auth/middleware"
	"github.com/stasBigunenko/monorepa/pkg/auth/routes"

	"github.com/gorilla/mux"
	er "github.com/stasBigunenko/monorepa/customErrors"
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
		return er.WrongPort
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
		Handler:      middleware.JSONRespHeaders(middleware.AccessControlMiddleware(s.router)),
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

	log.Info("Server is running on: " + s.getHTTPAddress())

	<-ctx.Done() // wait end of work

	/**
	 * start graceful shutdown
	 */

	log.Warn("Server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), time.Duration(cancelTimeout)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutDown); err != nil {
		return err
	}

	log.Warn("Server exited properly")
	return nil
}
