package httphandler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ItemsGrpcService interface {
	GetItems() ([]Item, error)
}

type HTTPHandler struct {
	ItemsService   ItemsGrpcService
	JwtServiceAddr string
}

func New(service ItemsGrpcService, addr string) HTTPHandler {
	return HTTPHandler{
		ItemsService:   service,
		JwtServiceAddr: addr,
	}
}

// if not internal server error, we provide err message to the user
func (h HTTPHandler) reportError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")

	if status == http.StatusInternalServerError {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(Error{
		Message: err.Error(),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(res) //nolint:errcheck
}

func (h HTTPHandler) GetRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/items", h.ListItems).Methods("GET")
	router.Use(h.authMiddleware)

	return router
}

// Returns the list of items
func (h HTTPHandler) ListItems(w http.ResponseWriter, req *http.Request) {
	items, err := h.ItemsService.GetItems()
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

func (h HTTPHandler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("Auth starting...")

		// here goes auth

		next.ServeHTTP(w, req)
	})
}
