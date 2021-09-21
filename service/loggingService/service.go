package loggingservice

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ContextKey string

const ContextKeyRequestID ContextKey = "requestID"

type LoggingService struct {
}

func New() LoggingService {
	return LoggingService{}
}

func (h LoggingService) WriteLog(w http.ResponseWriter, req *http.Request, message string) {
	id, ok := req.Context().Value(ContextKeyRequestID).(string)
	if !ok {
		log.Info("failed to onvert context value and get context id")
	}

	log.Info(fmt.Sprintf("Request ID: %s; Message: %s", id, message))
}
