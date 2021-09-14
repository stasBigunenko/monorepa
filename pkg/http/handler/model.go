package httphandler

import (
	"github.com/google/uuid"
)

type Person struct {
	Name string `json:"name,omitempty"`
}

type Item struct {
	ID uuid.UUID `json:"id,omitempty"`
}

type Error struct {
	Message string `json:"message,omitempty"`
}

func (e Error) Error() string {
	return e.Message
}
