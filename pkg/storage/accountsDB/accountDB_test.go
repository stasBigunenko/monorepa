package accountsDB

import (
	"context"
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAccountDB(t *testing.T) {
	acc := NewAccountStorage()

	id1 := uuid.New()
	userID1 := uuid.New()
	account1 := model.Account{
		ID:      id1,
		UserID:  userID1,
		Balance: 3000,
	}
	acc.MapAccount[id1] = account1

	id2 := uuid.New()
	userID2 := uuid.New()
	account2 := model.Account{
		ID:      id2,
		UserID:  userID2,
		Balance: 1500,
	}
	acc.MapAccount[id2] = account2

	id3 := uuid.New()
	account3 := model.Account{
		ID:      id3,
		UserID:  acc.MapAccount[id2].UserID,
		Balance: 100,
	}
	acc.MapAccount[id3] = account3

	t.Run("GetAccount", func(t *testing.T) {
		val, _ := acc.GetAccount(context.Background(), id1)
		require.Equal(t, account1, val, "account should be")

		_, err := acc.GetAccount(context.Background(), uuid.New())
		require.Errorf(t, err, "doesn't work")
	})

	t.Run("GetUserAccounts", func(t *testing.T) {
		accountsSet := []model.Account{
			account2,
			account3,
		}

		val, _ := acc.GetUserAccounts(context.Background(), userID2)

		require.Equal(t, accountsSet, val, "should be slice of model.Account with the same userID")
	})

	t.Run("GetAllUsers", func(t *testing.T) {
		accountsSet2 := []model.Account{
			account1,
			account2,
			account3,
		}

		val, _ := acc.GetAllUsers(context.Background())

		require.Equal(t, accountsSet2, val, "cann't show all users")
	})

}

func TestModificateAccountDB(t *testing.T) {
	acc := NewAccountStorage()

	id1 := uuid.New()
	userID1 := uuid.New()
	account1 := model.Account{
		ID:      id1,
		UserID:  userID1,
		Balance: 3000,
	}
	acc.MapAccount[id1] = account1

	id2 := uuid.New()
	userID2 := uuid.New()
	account2 := model.Account{
		ID:      id2,
		UserID:  userID2,
		Balance: 1500,
	}
	acc.MapAccount[id2] = account2

	id3 := uuid.New()
	account3 := model.Account{
		ID:      id3,
		UserID:  acc.MapAccount[id2].UserID,
		Balance: 100,
	}
	acc.MapAccount[id3] = account3
	t.Run("UpdateAccount", func(t *testing.T) {
		errAccount2 := account3
		errAccount2.Balance = 10000

		val, _ := acc.UpdateAccount(context.Background(), errAccount2)
		require.Equal(t, errAccount2, val, "should update account according id")

		errAccount := account3
		errAccount.ID = uuid.New()
		_, err := acc.UpdateAccount(context.Background(), errAccount)
		require.Errorf(t, err, "should be error not found")
	})

	t.Run("DeleteAccount", func(t *testing.T) {
		err := acc.DeleteAccount(context.Background(), account1.ID)
		require.Empty(t, model.Account{})

		err = acc.DeleteAccount(context.Background(), account1.ID)
		require.Errorf(t, err, "should be not found")
	})
}

func TestAccountStorage_CreateAccount(t *testing.T) {
	acc := NewAccountStorage()
	val, _ := acc.CreateAccount(context.Background(), uuid.New())
	require.NotNil(t, val, "should be no error")
}
