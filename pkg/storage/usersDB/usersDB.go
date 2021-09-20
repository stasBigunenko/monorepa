package usersDB

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
	"sync"
)

type UserStorage struct {
	MapUserHTTP map[uuid.UUID]model.UserHTTP
	mu          sync.Mutex
}

func NewUserStorage() *UserStorage {
	us := UserStorage{}
	us.MapUserHTTP = make(map[uuid.UUID]model.UserHTTP)
	return &us
}

func (us *UserStorage) Get(_ context.Context, id uuid.UUID) (model.UserHTTP, error) {
	us.mu.Lock()
	val, ok := us.MapUserHTTP[id]
	us.mu.Unlock()
	if !ok {
		return model.UserHTTP{}, fmt.Errorf("user with id %d not found", id)
	}

	return val, nil
}
func (us *UserStorage) GetAllUsers(_ context.Context) ([]model.UserHTTP, error) {
	us.mu.Lock()

	var users []model.UserHTTP

	for _, user := range us.MapUserHTTP {
		users = append(users, user)
	}

	us.mu.Unlock()

	return users, nil
}
func (us *UserStorage) Create(_ context.Context, name string) (model.UserHTTP, error) {
	us.mu.Lock()

	id := uuid.New()

	us.mu.Unlock()

	user := model.UserHTTP{
		ID:   id,
		Name: name,
	}

	us.MapUserHTTP[id] = user

	return user, nil
}
func (us *UserStorage) Update(_ context.Context, user model.UserHTTP) (model.UserHTTP, error) {
	us.mu.Lock()
	if _, ok := us.MapUserHTTP[user.ID]; !ok {
		return model.UserHTTP{}, fmt.Errorf("user with id %d not found", user.ID)
	}

	us.MapUserHTTP[user.ID] = user
	us.mu.Unlock()
	return user, nil
}
func (us *UserStorage) Delete(_ context.Context, id uuid.UUID) error {

	if _, ok := us.MapUserHTTP[id]; !ok {
		return fmt.Errorf("user with id %d not found", id)
	}

	delete(us.MapUserHTTP, id)

	return nil
}
