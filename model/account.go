package model

import "github.com/google/uuid"

type Account struct {
	ID      uuid.UUID `json:"id,omitempty"`
	UserID  uuid.UUID `json:"user_id,omitempty"`
	Balance int       `json:"balance"`
}
