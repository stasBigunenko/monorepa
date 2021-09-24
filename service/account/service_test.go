package account

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/mocks/pkg/storage/mockNewStore"
	"github.com/stasBigunenko/monorepa/model"
	"github.com/stretchr/testify/assert"
)

type MockLoggingService struct {
}

func (s MockLoggingService) WriteLog(ctx context.Context, message string) {}

func Test_Create(t *testing.T) {
	loggingService := MockLoggingService{}
	ui := new(mockNewStore.NewStore)
	id := uuid.New()
	userID := uuid.New()
	m := model.Account{ID: id, UserID: userID, Balance: 0}
	mm := model.Account{UserID: userID}
	ui.On("Create", context.Background(), mm).Return(m, nil)

	tests := []struct {
		name    string
		param   uuid.UUID
		stor    *mockNewStore.NewStore
		want    model.Account
		wantErr string
	}{
		{
			name:    "Everything good",
			param:   userID,
			stor:    ui,
			want:    m,
			wantErr: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewAccService(tc.stor, loggingService)
			got, err := u.Create(context.Background(), tc.param)
			log.Info(got)
			if err != nil && err.Error() != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

//
func Test_Get(t *testing.T) {
	loggingService := MockLoggingService{}
	ui := new(mockNewStore.NewStore)
	id := uuid.New()
	userID := uuid.New()
	m := model.Account{ID: id, UserID: userID, Balance: 0}
	ui.On("Get", context.Background(), id).Return(m, nil)

	ui2 := new(mockNewStore.NewStore)
	ui2.On("Get", mock.Anything, mock.Anything).Return(nil, nil)

	tests := []struct {
		name    string
		stor    *mockNewStore.NewStore
		want    model.Account
		wantErr string
	}{
		{
			name: "Everything ok",
			stor: ui,
			want: m,
		},
		{
			name:    "Everything ok",
			stor:    ui2,
			want:    model.Account{},
			wantErr: "couldn't get account",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewAccService(tc.stor, loggingService)
			got, err := u.Get(context.Background(), id)
			if err != nil && err.Error() != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	loggingService := MockLoggingService{}
	ui := new(mockNewStore.NewStore)
	id := uuid.New()
	ui.On("Delete", context.Background(), id).Return(nil)

	tests := []struct {
		name    string
		stor    *mockNewStore.NewStore
		param   uuid.UUID
		result  []byte
		wantErr string
	}{
		{
			name:    "Everything ok",
			stor:    ui,
			param:   id,
			wantErr: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewAccService(tc.stor, loggingService)
			err := u.Delete(context.Background(), tc.param)
			if err != nil && err.Error() != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestUserService_GetAll(t *testing.T) {
	loggingService := MockLoggingService{}
	ui := new(mockNewStore.NewStore)
	m1 := model.Account{ID: uuid.New(), UserID: uuid.New(), Balance: 0}
	m2 := model.Account{ID: uuid.New(), UserID: uuid.New(), Balance: 12}
	m := []model.Account{
		m1,
		m2,
	}
	ui.On("GetAll", context.Background()).Return(m, nil)

	tests := []struct {
		name    string
		stor    *mockNewStore.NewStore
		want    []model.Account
		wantErr string
	}{
		{
			name: "Everything ok",
			stor: ui,
			want: m,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewAccService(tc.stor, loggingService)
			got, err := u.GetAll(context.Background())
			if err != nil {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestUserService_Update(t *testing.T) {
	loggingService := MockLoggingService{}
	ui := new(mockNewStore.NewStore)
	id := uuid.New()
	userID := uuid.New()
	m := model.Account{ID: id, UserID: userID, Balance: 0}
	ui.On("Update", context.Background(), m).Return(m, nil)

	ui2 := new(mockNewStore.NewStore)
	ui2.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("not found"))

	tests := []struct {
		name    string
		stor    *mockNewStore.NewStore
		param   model.Account
		want    model.Account
		wantErr string
	}{
		{
			name:  "Everything ok",
			stor:  ui,
			param: m,
			want:  m,
		},
		{
			name:    "Couldn't unmarshal",
			stor:    ui2,
			param:   m,
			want:    model.Account{},
			wantErr: "not found",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewAccService(tc.stor, loggingService)
			got, err := u.Update(context.Background(), m)
			if (err != nil) && err.Error() != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	loggingService := MockLoggingService{}
	ui := new(mockNewStore.NewStore)
	userID := uuid.New()
	m1 := model.Account{ID: uuid.New(), UserID: userID, Balance: 0}
	m2 := model.Account{ID: uuid.New(), UserID: userID, Balance: 12}
	m := []model.Account{
		m1,
		m2,
	}
	ui.On("GetUserAccounts", context.Background(), userID).Return(m, nil)

	tests := []struct {
		name    string
		stor    *mockNewStore.NewStore
		param   uuid.UUID
		want    []model.Account
		wantErr string
	}{
		{
			name:  "Everything ok",
			stor:  ui,
			param: userID,
			want:  m,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewAccService(tc.stor, loggingService)
			got, err := u.GetUser(context.Background(), tc.param)
			if (err != nil) && err.Error() != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
