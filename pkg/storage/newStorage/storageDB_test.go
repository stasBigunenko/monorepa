package newStorage

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
	"github.com/stretchr/testify/assert"
)

type MockLoggingService struct {
}

func (s MockLoggingService) WriteLog(ctx context.Context, message string) {}

func TestStorageDB_Get(t *testing.T) {
	loggingService := MockLoggingService{}
	acc := NewDB(loggingService)
	id := uuid.New()
	userID := uuid.New()
	m := model.Account{ID: id, UserID: userID, Balance: 0}
	acc.Data[id] = m

	tests := []struct {
		name    string
		param   uuid.UUID
		want    model.Account
		wantErr string
	}{
		{
			name:  "Everything ok",
			param: id,
			want:  m,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := acc.Get(context.Background(), tc.param)
			if (err != nil) && err.Error() != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

//
func TestStorageDB_GetUser(t *testing.T) {
	loggingService := MockLoggingService{}
	acc := NewDB(loggingService)
	id := uuid.New()
	userID := uuid.New()
	m := model.Account{ID: id, UserID: userID, Balance: 0}
	acc.Data[id] = m
	res := []model.Account{
		m,
	}

	tests := []struct {
		name string
		want []model.Account
	}{
		{
			name: "Everything ok",
			want: res,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := acc.GetUserAccounts(context.Background(), userID)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStorageDB_GetAll(t *testing.T) {
	loggingService := MockLoggingService{}
	acc := NewDB(loggingService)
	id := uuid.New()
	userID := uuid.New()
	m := model.Account{ID: id, UserID: userID, Balance: 0}
	acc.Data[id] = m
	res := []model.Account{
		m,
	}

	tests := []struct {
		name string
		want []model.Account
	}{
		{
			name: "Everything ok",
			want: res,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := acc.GetAll(context.Background())
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStorageDB_Create(t *testing.T) {
	loggingService := MockLoggingService{}
	acc := NewDB(loggingService)
	m := model.UserHTTP{Name: "Jim"}

	tests := []struct {
		name  string
		param model.UserHTTP
	}{
		{
			name:  "Everything ok",
			param: m,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			val, _ := acc.Create(context.Background(), tc.param)
			require.NotNil(t, val, "should be no error")
		})
	}
}

func TestStorageDB_Delete(t *testing.T) {
	loggingService := MockLoggingService{}
	acc := NewDB(loggingService)
	id := uuid.New()
	userID := uuid.New()
	m := model.Account{ID: id, UserID: userID, Balance: 0}
	acc.Data[id] = m

	tests := []struct {
		name    string
		param   uuid.UUID
		want    []byte
		wantErr string
	}{
		{
			name:  "Everything ok",
			param: id,
			want:  nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := acc.Delete(context.Background(), id)
			assert.Empty(t, err)
		})
	}
}
