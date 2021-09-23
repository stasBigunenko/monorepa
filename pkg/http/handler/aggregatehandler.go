package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/stasBigunenko/monorepa/customErrors"
	"github.com/stasBigunenko/monorepa/model"
)

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
