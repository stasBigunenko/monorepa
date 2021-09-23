package httphandler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/stasBigunenko/monorepa/customErrors"
	"github.com/stasBigunenko/monorepa/model"
)

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

	forUser := true
	p, err := ioutil.ReadAll(req.Body)
	if err == io.EOF {
		forUser = false
	} else if err != nil {
		h.reportError(w, err)
		return
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
