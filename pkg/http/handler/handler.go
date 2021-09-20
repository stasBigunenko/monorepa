package httphandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/stasBigunenko/monorepa/customErrors"
	"github.com/stasBigunenko/monorepa/model"
	tokenservice "github.com/stasBigunenko/monorepa/service/http"
)

type AccountGrpcService interface {
	CreateAccount(userID uuid.UUID) (uuid.UUID, error)
	GetAccount(id uuid.UUID) (model.Account, error)
	GetUserAccounts(userID uuid.UUID) ([]model.Account, error)
	GetAllAccounts() ([]model.Account, error)
	UpdateAccount(account model.Account) error
	DeleteAccount(id uuid.UUID) error
}

type UserGrpcService interface {
	CreateUser(name string) (uuid.UUID, error)
	GetUser(id uuid.UUID) (model.UserHTTP, error)
	GetAllUsers() ([]model.UserHTTP, error)
	UpdateUser(user model.UserHTTP) error
	DeleteUser(id uuid.UUID) error
}

type TokenService interface {
	ParseToken(tokenPart string) (string, error)
}

type key int

const (
	nameKey key = iota
)

type HTTPHandler struct {
	AccountsService AccountGrpcService
	UsersService    UserGrpcService
	TokenService    TokenService
	JwtServiceAddr  string
}

func New(accountService AccountGrpcService, userService UserGrpcService, addr string) HTTPHandler {
	return HTTPHandler{
		AccountsService: accountService,
		UsersService:    userService,
		TokenService: tokenservice.HTTPService{
			JwtServiceAddr: addr,
		},
	}
}

func (h HTTPHandler) reportError(w http.ResponseWriter, err error) {
	var status int

	switch {
	case errors.Is(err, customErrors.UUIDError) || errors.Is(err, customErrors.JSONError) || errors.Is(err, customErrors.AlreadyExists):
		status = http.StatusBadRequest
	case errors.Is(err, customErrors.NotFound):
		status = http.StatusNotFound
	case errors.Is(err, customErrors.DeadlineExceeded):
		status = http.StatusGatewayTimeout
	default:
		status = http.StatusInternalServerError
	}

	w.WriteHeader(status)

	if status == http.StatusInternalServerError {
		log.Error(err)
		return
	}

	res, err := json.Marshal(err)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res) //nolint:errcheck
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

	// TODO: router.HandleFunc("/accounts/me", h.MyAccounts).Methods("GET")
	// TODO: /accounts add param --userid =
	// TODO: get user with all her accounts (aggregation query)

	router.Use(h.authMiddleware)

	return router
}

// ***** //
// Users //
// ***** //

func (h HTTPHandler) AddUser(w http.ResponseWriter, req *http.Request) {
	log.Info("Command AddUSer received...")

	name, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.reportError(w, err)
		return
	}

	var user model.UserHTTP
	if err = json.Unmarshal(name, &user); err != nil {
		h.reportError(w, err)
		return
	}

	accountID, err := h.UsersService.CreateUser(user.Name)
	if err != nil {
		h.reportError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/accounts/%d", accountID))
	w.WriteHeader(http.StatusCreated)

}

func (h HTTPHandler) GetUser(w http.ResponseWriter, req *http.Request) {
	log.Info("Command GetUser received...")

	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.UUIDError))
		return
	}

	user, err := h.UsersService.GetUser(id)
	if err != nil {
		h.reportError(w, err)
		return
	}

	u, err := json.Marshal(user)
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.JSONError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(u) //nolint:errcheck

}

func (h HTTPHandler) UpdateUser(w http.ResponseWriter, req *http.Request) {
	log.Info("Command UpdateUser received")

	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.UUIDError))
		return
	}

	p, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.reportError(w, err)
		return
	}

	user := model.UserHTTP{}
	if err = json.Unmarshal(p, &user); err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.JSONError))
		return
	}

	user.ID = id

	err = h.UsersService.UpdateUser(user)
	if err != nil {
		h.reportError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h HTTPHandler) DeleteUser(w http.ResponseWriter, req *http.Request) {
	log.Info("Command DeleteUser received...")

	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.UUIDError))
		return
	}

	if err := h.UsersService.DeleteUser(id); err != nil {
		h.reportError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h HTTPHandler) ListUsers(w http.ResponseWriter, req *http.Request) {
	log.Info("Command ListUsers received...")

	users, err := h.UsersService.GetAllUsers()
	if err != nil {
		h.reportError(w, err)
		return
	}

	res, err := json.Marshal(users)
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.JSONError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res) //nolint:errcheck
}

// ******** //
// Accounts //
// ******** //

func (h HTTPHandler) AddAccount(w http.ResponseWriter, req *http.Request) {
	log.Info("Command AddAccount received...")

	userID, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.reportError(w, err)
		return
	}

	var account model.Account
	if err = json.Unmarshal(userID, &account); err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.JSONError))
		return
	}

	accountID, err := h.AccountsService.CreateAccount(account.UserID)
	if err != nil {
		h.reportError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/accounts/%d", accountID))
	w.WriteHeader(http.StatusCreated)
}

func (h HTTPHandler) GetAccount(w http.ResponseWriter, req *http.Request) {
	log.Info("Command GetAccount received...")

	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.UUIDError))
		return
	}

	account, err := h.AccountsService.GetAccount(id)
	if err != nil {
		h.reportError(w, err)
		return
	}

	a, err := json.Marshal(account)
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.JSONError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(a) //nolint:errcheck
}

func (h HTTPHandler) UpdateAccount(w http.ResponseWriter, req *http.Request) {
	log.Info("Command UpdateAccount received...")

	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.UUIDError))
		return
	}

	p, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.reportError(w, err)
		return
	}

	account := model.Account{}
	if err = json.Unmarshal(p, &account); err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.JSONError))
		return
	}

	account.ID = id

	err = h.AccountsService.UpdateAccount(account)
	if err != nil {
		h.reportError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h HTTPHandler) DeleteAccount(w http.ResponseWriter, req *http.Request) {
	log.Info("Command DeleteAccount received...")

	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.UUIDError))
		return
	}

	if err := h.AccountsService.DeleteAccount(id); err != nil {
		h.reportError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h HTTPHandler) ListAccounts(w http.ResponseWriter, req *http.Request) {
	log.Info("Command ListAccount received...")

	withUser := true
	p, err := ioutil.ReadAll(req.Body)
	if err == io.EOF {
		withUser = false
	} else if err != nil {
		h.reportError(w, err)
		return
	}

	var accounts []model.Account
	if withUser {
		account := model.Account{}
		if err = json.Unmarshal(p, &account); err != nil {
			h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.JSONError))
			return
		}
		accounts, err = h.AccountsService.GetUserAccounts(account.UserID)
	} else {
		accounts, err = h.AccountsService.GetAllAccounts()
	}

	if err != nil {
		h.reportError(w, err)
		return
	}

	res, err := json.Marshal(accounts)
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.JSONError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res) //nolint:errcheck
}

func (h HTTPHandler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Info("Auth starting...")

		tokenHeader := req.Header.Get("Authorization")
		if tokenHeader == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		name, err := h.TokenService.ParseToken(tokenHeader)
		if err != nil {
			h.reportError(w, err)
			return
		}

		ctx := context.WithValue(req.Context(), nameKey, name)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
