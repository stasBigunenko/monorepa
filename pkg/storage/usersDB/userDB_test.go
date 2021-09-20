package usersDB

import (
	"context"
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserDB(t *testing.T) {
	us := NewUserStorage()

	id1 := uuid.New()
	user1 := model.UserHTTP{
		ID:   id1,
		Name: "Jack",
	}
	us.MapUserHTTP[id1] = user1

	id2 := uuid.New()
	user2 := model.UserHTTP{
		ID:   id2,
		Name: "Johnie",
	}
	us.MapUserHTTP[id2] = user2

	id3 := uuid.New()
	user3 := model.UserHTTP{
		ID:   id3,
		Name: "Jim",
	}
	us.MapUserHTTP[id3] = user3

	t.Run("Get", func(t *testing.T) {
		val, _ := us.Get(context.Background(), id1)
		require.Equal(t, user1, val, "account should be")

		_, err := us.Get(context.Background(), uuid.New())
		require.Errorf(t, err, "doesn't work")
	})

	t.Run("GetAllUsers", func(t *testing.T) {
		usersSet := []model.UserHTTP{
			user1,
			user2,
			user3,
		}

		val, _ := us.GetAllUsers(context.Background())

		require.Equal(t, usersSet, val, "should be slice of users")
	})
}
func TestModificateUser(t *testing.T) {
	us := NewUserStorage()

	id1 := uuid.New()
	user1 := model.UserHTTP{
		ID:   id1,
		Name: "Jack",
	}
	us.MapUserHTTP[id1] = user1

	id2 := uuid.New()
	user2 := model.UserHTTP{
		ID:   id2,
		Name: "Johnie",
	}
	us.MapUserHTTP[id2] = user2

	id3 := uuid.New()
	user3 := model.UserHTTP{
		ID:   id3,
		Name: "Jim",
	}
	us.MapUserHTTP[id3] = user3

	t.Run("Update", func(t *testing.T) {
		user2.Name = "Glen"

		val, _ := us.Update(context.Background(), user2)
		require.Equal(t, user2, val, "should update account according id")

		errAccount := model.UserHTTP{
			ID:   uuid.New(),
			Name: "Empty",
		}
		_, err := us.Update(context.Background(), errAccount)
		require.Errorf(t, err, "should be error not found")
	})

	t.Run("Delete=", func(t *testing.T) {
		delUser := user1
		us.Delete(context.Background(), delUser.ID)
		require.Empty(t, model.Account{})

		err := us.Delete(context.Background(), delUser.ID)
		require.Errorf(t, err, "should be not found")
	})
}

func TestAccountStorage_CreateAccount(t *testing.T) {
	us := NewUserStorage()
	val, _ := us.Create(context.Background(), "newMan")
	require.NotNil(t, val, "should be no error")
}
