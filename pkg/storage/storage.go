package storage

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"math/rand"
)

type Storage struct {
	Id          string
	Title       string
	Description string
}

func NewStorage() *Storage {
	return &Storage{}
}

const letterBytes = "abcdefgh ijklmnop qrstuvw xyz"
const letterBytesUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandStringBytesfunction generate lower runes and spaces
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// RandStringBytesUpper function generate upper runes
func RandStringBytesUpper(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytesUpper[rand.Intn(len(letterBytesUpper))]
	}
	return string(b)
}

// GetItems generate slice of Storage according the username and send back
func (d *Storage) GetItems(c context.Context, u string) ([]Storage, error) {
	items := []Storage{}

	if u == "" {
		return nil, errors.New("invalid user name")
	}

	for i := range u {

		id := uuid.New()
		title := RandStringBytesUpper(len(id))
		descr := string(u[i]) + title
		description := RandStringBytes(len(descr))

		items = append(items, Storage{
			Id:          id.String(),
			Title:       title,
			Description: description,
		})
	}
	return items, nil
}
