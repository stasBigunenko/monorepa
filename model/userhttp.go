package model

import "github.com/google/uuid"

type UserHTTP struct {
	ID   uuid.UUID `json:"id,omitempty"`
	Name string    `json:"name,omitempty"`
}
