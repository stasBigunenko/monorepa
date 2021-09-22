package loggingservice

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type ContextKey string

const ContextKeyRequestID ContextKey = "requestID"

type LoggingService struct {
}

func New() LoggingService {
	return LoggingService{}
}

func (h LoggingService) WriteLog(ctx context.Context, message string) {
	id, ok := ctx.Value(ContextKeyRequestID).(string)
	if !ok {
		log.Info("failed to onvert context value and get context id")
	}

	log.Info(fmt.Sprintf("Request ID: %s; Message: %s", id, message))
}
