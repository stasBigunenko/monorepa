package httphandler

import (
	"github.com/gorilla/mux"

	tokenservice "github.com/stasBigunenko/monorepa/service/http"
	loggingservice "github.com/stasBigunenko/monorepa/service/loggingService"
)

type HTTPHandler struct {
	AccountsService AccountGrpcService
	UsersService    UserGrpcService
	TokenService    TokenService
	LoggingService  LoggingService
	JwtServiceAddr  string
}

func New(accountService AccountGrpcService, userService UserGrpcService, loggingService LoggingService, addr string) *HTTPHandler {
	return &HTTPHandler{
		AccountsService: accountService,
		UsersService:    userService,
		TokenService: tokenservice.HTTPService{
			JwtServiceAddr: addr,
		},
		LoggingService: loggingservice.LoggingService{},
	}
}

func (h HTTPHandler) GetRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/users", h.AddUser).Methods("POST")
	router.HandleFunc("/users/{id}", h.GetUser).Methods("GET")
	router.HandleFunc("/users", h.ListUsers).Methods("GET")
	router.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")

	router.HandleFunc("/accounts", h.AddAccount).Methods("POST")
	router.HandleFunc("/accounts/{id}", h.GetAccount).Methods("GET")
	router.HandleFunc("/accounts/{id}", h.UpdateAccount).Methods("PUT")
	router.HandleFunc("/accounts/{id}", h.DeleteAccount).Methods("DELETE")
	router.HandleFunc("/accounts", h.ListAccounts).Methods("GET")

	router.HandleFunc("/accounts_and_user/{id}", h.GetAggregate).Methods("GET")

	router.Use(h.AuthMiddleware)
	router.Use(h.RequestIDMiddleware)

	return router
}
