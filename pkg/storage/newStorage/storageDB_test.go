package newStorage

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStorageDB(t *testing.T) {
	acc := DB()

	id1 := uuid.New()
	userID1 := uuid.New()
	account1 := model.Account{
		ID:      id1,
		UserID:  userID1,
		Balance: 3000,
	}
	res1, _ := json.Marshal(account1)
	acc.Data[id1] = res1

	id2 := uuid.New()
	userID2 := uuid.New()
	account2 := model.Account{
		ID:      id2,
		UserID:  userID2,
		Balance: 1500,
	}
	res2, _ := json.Marshal(account2)
	acc.Data[id2] = res2

	id3 := uuid.New()
	account3 := model.Account{
		ID:      id3,
		UserID:  userID2,
		Balance: 100,
	}
	res3, _ := json.Marshal(account3)
	acc.Data[id3] = res3

	accountsSet := []model.Account{
		account2,
		account3,
	}
	resSetUser := [][]byte{}
	for _, val := range accountsSet {
		res, _ := json.Marshal(val)
		resSetUser = append(resSetUser, res)
	}

	accountsSet2 := []model.Account{
		account1,
		account2,
		account3,
	}
	resSetAll := [][]byte{}
	for _, val := range accountsSet2 {
		res, _ := json.Marshal(val)
		resSetAll = append(resSetAll, res)
	}

	t.Run("Get", func(t *testing.T) {
		val, _ := acc.Get(context.Background(), id1)
		require.Equal(t, res1, val, "account should be")

		_, err := acc.Get(context.Background(), uuid.New())
		require.Empty(t, err)
	})

	t.Run("GetUser", func(t *testing.T) {

		val, _ := acc.GetUser(context.Background(), userID2)
		require.Equal(t, resSetUser, val, "should be slice of model.Account with the same userID")
	})

	t.Run("GetAll", func(t *testing.T) {

		val, _ := acc.GetAll(context.Background())
		require.Equal(t, resSetAll, val, "cann't show all users")
	})
}

func TestNewStoreModificate(t *testing.T) {
	acc2 := DB()

	id1 := uuid.New()
	userID1 := uuid.New()
	account1 := model.Account{
		ID:      id1,
		UserID:  userID1,
		Balance: 3000,
	}
	res1, _ := json.Marshal(account1)
	acc2.Data[id1] = res1

	id2 := uuid.New()
	userID2 := uuid.New()
	account2 := model.Account{
		ID:      id2,
		UserID:  userID2,
		Balance: 1500,
	}
	res2, _ := json.Marshal(account2)
	acc2.Data[id2] = res2

	t.Run("Update", func(t *testing.T) {
		account1.Balance = 20
		resUpd, _ := json.Marshal(account1)

		val, _ := acc2.Update(context.Background(), id1, resUpd)
		require.Equal(t, resUpd, val, "should update account according id")

		errUpd := account1
		errUpd.ID = uuid.New()
		resErrUpd, _ := json.Marshal(errUpd)
		_, err := acc2.Update(context.Background(), errUpd.ID, resErrUpd)
		require.Empty(t, err, "should be nil")
	})

	t.Run("Delete", func(t *testing.T) {
		err := acc2.Delete(context.Background(), id2)
		require.Empty(t, err)

		err = acc2.Delete(context.Background(), id2)
		require.Empty(t, err, "should be nil")
	})

	t.Run("Create", func(t *testing.T) {
		username := "Jim"
		r, _ := json.Marshal(username)
		val, _, _ := acc2.Create(context.Background(), r)
		require.NotNil(t, val, "should be no error")
	})
}
