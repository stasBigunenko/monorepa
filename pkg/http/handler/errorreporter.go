package httphandler

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/stasBigunenko/monorepa/customErrors"
)

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
