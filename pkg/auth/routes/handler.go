package routes

import (
	"context"
	"encoding/json"
	"monorepa/model"
	"monorepa/service/auth"
	"net/http"

	"github.com/gorilla/mux"
)

type HandlerItemsServ struct {
	router   *mux.Router
	services *auth.MockAuthService
	ctx      context.Context
}

func New(ctx context.Context, router *mux.Router) *HandlerItemsServ {
	return &HandlerItemsServ{
		router:   router,
		services: &auth.MockAuthService{},
		ctx:      ctx,
	}
}

func (h *HandlerItemsServ) HandlerItems() {
	h.router.HandleFunc("/login", h.GetJWTToken).Methods("GET")
	h.router.HandleFunc("/get-cert/{version}", h.GetCertKey).Methods("GET")
}

func (h *HandlerItemsServ) GetJWTToken(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := h.services.Login(user)
	if err != nil {
		/**
		 * change to global err
		 */
		if err.Error() == "wrong password" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("token", token)
	w.WriteHeader(http.StatusCreated)
}

func (h *HandlerItemsServ) GetCertKey(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	versionCert := params["version"]
	if versionCert == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	byteCertPbKey, err := h.services.GetCert(versionCert)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// resp struct with token
	// decode and return json
	type tokenResp struct {
		PbKey []byte `json:"publicKey"`
	}

	resp := tokenResp{
		PbKey: byteCertPbKey,
	}

	respByte, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(respByte)
	w.WriteHeader(http.StatusOK)
}
