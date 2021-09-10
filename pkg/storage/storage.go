package storage

import (
	"errors"
	"github.com/google/uuid"
	"math/rand"
)

type Storage struct {
	Id          string
	Title       string
	Description string
}

func NewDataProvider() *Storage {
	return &Storage{}
}

const letterBytes = "abcdefgh ijklmnop qrstuvw xyz"
const letterBytesUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func RandStringBytesUpper(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytesUpper[rand.Intn(len(letterBytesUpper))]
	}
	return string(b)
}

func (d *Storage) GetItems(u string) ([]Storage, error) {
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
