package httphandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/stasBigunenko/monorepa/customErrors"
	"github.com/stasBigunenko/monorepa/model"
)

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
