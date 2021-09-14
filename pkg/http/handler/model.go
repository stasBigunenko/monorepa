package httphandler

import (
	"github.com/google/uuid"
)

type Item struct {
	id uuid.UUID
}

type Error struct {
	Message string `json:"message,omitempty"`
}

func (e Error) Error() string {
	return e.Message
}
