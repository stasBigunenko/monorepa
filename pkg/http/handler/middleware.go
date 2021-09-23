package httphandler

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// ********** //
// Middleware //
// ********** //

type ContextKey string

const (
	NameKey             ContextKey = "name"
	ContextKeyRequestID ContextKey = "requestID"
)

func (h HTTPHandler) AuthMiddleware(next http.Handler) http.Handler {
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

func (h HTTPHandler) RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		name, ok := ctx.Value(NameKey).(string)
		if !ok {
			h.reportError(w, errors.New("failed to generate context value"))
			return
		}

		requestID := name + "_" + uuid.New().String()

		ctx = context.WithValue(ctx, ContextKeyRequestID, requestID)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func (h HTTPHandler) AccessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
