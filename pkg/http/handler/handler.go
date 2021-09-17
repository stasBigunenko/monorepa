package httphandler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/stasBigunenko/monorepa/model"

	"github.com/gorilla/mux"

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

// if internal server error, we provide err message to log, else to the user
func (h HTTPHandler) reportError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")

	if status == http.StatusInternalServerError {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

	res, err := json.Marshal(Error{
		Message: err.Error(),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
	router.HandleFunc("/accounts/me", h.MyAccounts).Methods("GET")
	router.HandleFunc("/accounts/{id}", h.GetAccount).Methods("GET")
	router.HandleFunc("/accounts/{id}", h.UpdateAccount).Methods("PUT")
	router.HandleFunc("/accounts/{id}", h.DeleteAccount).Methods("DELETE")
	router.HandleFunc("/accounts", h.ListAccounts).Methods("GET") //--userid =
	//router.HandleFunc("/accounts/{userid}", h.UserAccounts).Methods("GET")

	router.Use(h.authMiddleware)

	return router
}

// ***** //
// Users //
// ***** //

func (h HTTPHandler) AddUser(w http.ResponseWriter, req *http.Request) {
}

func (h HTTPHandler) GetUser(w http.ResponseWriter, req *http.Request) {
}

func (h HTTPHandler) UpdateUser(w http.ResponseWriter, req *http.Request) {
}

func (h HTTPHandler) DeleteUser(w http.ResponseWriter, req *http.Request) {
}

func (h HTTPHandler) ListUsers(w http.ResponseWriter, req *http.Request) {
	// name := req.Context().Value(nameKey)
	// nameStr, ok := name.(string)
	// if !ok {
	// 	h.reportError(w, http.StatusInternalServerError, errors.New("failed to convert name to string"))
	// 	return
	// }

	items, err := h.UsersService.GetAllUsers()
	if err != nil {
		h.reportError(w, http.StatusInternalServerError, err)
		return
	}

	res, err := json.Marshal(items)
	if err != nil {
		h.reportError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res) //nolint:errcheck

}

// ******** //
// Accounts //
// ******** //

func (h HTTPHandler) AddAccount(w http.ResponseWriter, req *http.Request) {
}

func (h HTTPHandler) GetAccount(w http.ResponseWriter, req *http.Request) {
}

func (h HTTPHandler) UpdateAccount(w http.ResponseWriter, req *http.Request) {
}

func (h HTTPHandler) DeleteAccount(w http.ResponseWriter, req *http.Request) {
}

func (h HTTPHandler) ListAccounts(w http.ResponseWriter, req *http.Request) {
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
			h.reportError(w, http.StatusInternalServerError, err)
			return
		}

		ctx := context.WithValue(req.Context(), nameKey, name)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
