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
	loggingservice "github.com/stasBigunenko/monorepa/service/loggingService"
)

type AccountGrpcService interface {
	CreateAccount(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	GetAccount(ctx context.Context, id uuid.UUID) (model.Account, error)
	GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]model.Account, error)
	GetAllAccounts(ctx context.Context) ([]model.Account, error)
	UpdateAccount(ctx context.Context, account model.Account) error
	DeleteAccount(ctx context.Context, id uuid.UUID) error
}

type UserGrpcService interface {
	CreateUser(ctx context.Context, name string) (uuid.UUID, error)
	GetUser(ctx context.Context, id uuid.UUID) (model.UserHTTP, error)
	GetAllUsers(ctx context.Context) ([]model.UserHTTP, error)
	UpdateUser(ctx context.Context, user model.UserHTTP) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type TokenService interface {
	ParseToken(tokenPart string) (string, error)
}

type LoggingService interface {
	WriteLog(ctx context.Context, message string)
}

type ContextKey string

const (
	NameKey             ContextKey = "name"
	ContextKeyRequestID ContextKey = "requestID"
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

	router.HandleFunc("/accounts_and_user/{id}", h.GetAggregate).Methods("GET")

	router.Use(h.authMiddleware)
	router.Use(h.requestIDMiddleware)

	return router
}

// ***** //
// Users //
// ***** //

func (h HTTPHandler) AddUser(w http.ResponseWriter, req *http.Request) {
	h.LoggingService.WriteLog(req.Context(), "HTTTP: Command AddUSer received...")

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

	accountID, err := h.UsersService.CreateUser(req.Context(), user.Name)
	if err != nil {
		h.reportError(w, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/accounts/%d", accountID))
	w.WriteHeader(http.StatusCreated)
}

func (h HTTPHandler) GetUser(w http.ResponseWriter, req *http.Request) {
	h.LoggingService.WriteLog(req.Context(), "HTTTP: Command GetUser received...")

	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.UUIDError))
		return
	}

	user, err := h.UsersService.GetUser(req.Context(), id)
	if err != nil {
		h.reportError(w, err)
		return
	}

	u, err := json.Marshal(user)
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.JSONError))
		return
	}

	w.Write(u) //nolint:errcheck

}

func (h HTTPHandler) UpdateUser(w http.ResponseWriter, req *http.Request) {
	h.LoggingService.WriteLog(req.Context(), "HTTTP: Command UpdateUser received...")

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

	err = h.UsersService.UpdateUser(req.Context(), user)
	if err != nil {
		h.reportError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h HTTPHandler) DeleteUser(w http.ResponseWriter, req *http.Request) {
	h.LoggingService.WriteLog(req.Context(), "HTTTP: Command DeleteUser received...")

	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.UUIDError))
		return
	}

	if err := h.UsersService.DeleteUser(req.Context(), id); err != nil {
		h.reportError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h HTTPHandler) ListUsers(w http.ResponseWriter, req *http.Request) {
	h.LoggingService.WriteLog(req.Context(), "HTTTP: Command ListUsers received...")

	users, err := h.UsersService.GetAllUsers(req.Context())
	if err != nil {
		h.reportError(w, err)
		return
	}

	res, err := json.Marshal(users)
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.JSONError))
		return
	}

	w.Write(res) //nolint:errcheck
}

// ******** //
// Accounts //
// ******** //

func (h HTTPHandler) AddAccount(w http.ResponseWriter, req *http.Request) {
	h.LoggingService.WriteLog(req.Context(), "HTTTP: Command AddAccount received...")

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

	accountID, err := h.AccountsService.CreateAccount(req.Context(), account.UserID)
	if err != nil {
		h.reportError(w, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/accounts/%d", accountID))
	w.WriteHeader(http.StatusCreated)
}

func (h HTTPHandler) GetAccount(w http.ResponseWriter, req *http.Request) {
	h.LoggingService.WriteLog(req.Context(), "HTTTP: Command GetAccount received...")

	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.UUIDError))
		return
	}

	account, err := h.AccountsService.GetAccount(req.Context(), id)
	if err != nil {
		h.reportError(w, err)
		return
	}

	a, err := json.Marshal(account)
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.JSONError))
		return
	}

	w.Write(a) //nolint:errcheck
}

func (h HTTPHandler) UpdateAccount(w http.ResponseWriter, req *http.Request) {
	h.LoggingService.WriteLog(req.Context(), "HTTTP: Command UpdateAccount received...")

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

	err = h.AccountsService.UpdateAccount(req.Context(), account)
	if err != nil {
		h.reportError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h HTTPHandler) DeleteAccount(w http.ResponseWriter, req *http.Request) {
	h.LoggingService.WriteLog(req.Context(), "HTTTP: Command DeleteAccount received...")

	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.UUIDError))
		return
	}

	if err := h.AccountsService.DeleteAccount(req.Context(), id); err != nil {
		h.reportError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h HTTPHandler) ListAccounts(w http.ResponseWriter, req *http.Request) {
	h.LoggingService.WriteLog(req.Context(), "HTTTP: Command ListAccount received...")

	forUser := false
	p, err := ioutil.ReadAll(req.Body)
	if err != nil {
		if err != io.EOF {
			forUser = true
		} else {
			h.reportError(w, err)
			return
		}
	}

	var accounts []model.Account
	if forUser {
		account := model.Account{}
		if err = json.Unmarshal(p, &account); err != nil {
			h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.JSONError))
			return
		}
		accounts, err = h.AccountsService.GetUserAccounts(req.Context(), account.UserID)
	} else {
		accounts, err = h.AccountsService.GetAllAccounts(req.Context())
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

	w.Write(res) //nolint:errcheck
}

func (h HTTPHandler) GetAggregate(w http.ResponseWriter, req *http.Request) {
	h.LoggingService.WriteLog(req.Context(), "HTTTP: Command GetAggregate received...")

	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.UUIDError))
		return
	}

	accounts, err := h.AccountsService.GetUserAccounts(req.Context(), id)
	if err != nil {
		h.reportError(w, err)
		return
	}

	user, err := h.UsersService.GetUser(req.Context(), id)
	if err != nil {
		h.reportError(w, err)
		return
	}

	aggregated := model.UserAndAccounts{
		User:     user,
		Accounts: accounts,
	}

	res, err := json.Marshal(aggregated)
	if err != nil {
		h.reportError(w, fmt.Errorf("%s: %w", err, customErrors.JSONError))
		return
	}

	w.Write(res) //nolint:errcheck
}

// ******** //
// Middleware //
// ******** //

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

		w.Header().Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), NameKey, name)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func (h HTTPHandler) requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		name, ok := ctx.Value(NameKey).(string)
		if !ok {
			h.reportError(w, errors.New("failed to generate context value"))
		}

		requestID := name + "_" + uuid.New().String()

		ctx = context.WithValue(ctx, ContextKeyRequestID, requestID)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
