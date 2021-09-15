package model

import "github.com/google/uuid"

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Item struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
}
