package account

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/mocks/pkg/storage/mockNewStore"
	"github.com/stasBigunenko/monorepa/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Create(t *testing.T) {
	ui := new(mockNewStore.NewStore)
	id := uuid.New()
	userID := uuid.New()
	m := model.Account{ID: id, UserID: userID, Balance: 0}
	j, _ := json.Marshal(userID)
	ui.On("Create", context.Background(), j).Return(j, id, nil)

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
			u := NewAccService(tc.stor)
			got, err := u.Create(context.Background(), tc.param)
			if err != nil && err.Error() != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_Get(t *testing.T) {
	ui := new(mockNewStore.NewStore)
	id := uuid.New()
	userID := uuid.New()
	m := model.Account{ID: id, UserID: userID, Balance: 0}
	mj, _ := json.Marshal(m)
	ui.On("Get", context.Background(), id).Return(mj, nil)

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
			want: model.Account{
				ID:      id,
				UserID:  userID,
				Balance: 0,
			},
		},
		{
			name:    "Everything ok",
			stor:    ui2,
			want:    model.Account{},
			wantErr: "couldn't unmarshal data",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewAccService(tc.stor)
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
	ui := new(mockNewStore.NewStore)
	id := uuid.New()
	ui.On("Delete", context.Background(), id).Return(true, nil)

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
			u := NewAccService(tc.stor)
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
	ui := new(mockNewStore.NewStore)
	m1 := model.Account{ID: uuid.New(), UserID: uuid.New(), Balance: 0}
	m2 := model.Account{ID: uuid.New(), UserID: uuid.New(), Balance: 12}
	m := []model.Account{
		m1,
		m2,
	}
	var mj [][]byte

	for _, val := range m {
		j, _ := json.Marshal(val)
		mj = append(mj, j)
	}
	ui.On("GetAll", context.Background()).Return(mj, nil)

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
			u := NewAccService(tc.stor)
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
	ui := new(mockNewStore.NewStore)
	id := uuid.New()
	userID := uuid.New()
	m := model.Account{ID: id, UserID: userID, Balance: 0}
	mj, _ := json.Marshal(m)
	ui.On("Update", context.Background(), id, mj).Return(mj, nil)

	ui2 := new(mockNewStore.NewStore)
	ui2.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

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
			wantErr: "couldn't unmarshal data",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewAccService(tc.stor)
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
	ui := new(mockNewStore.NewStore)
	userID := uuid.New()
	m1 := model.Account{ID: uuid.New(), UserID: userID, Balance: 0}
	m2 := model.Account{ID: uuid.New(), UserID: userID, Balance: 12}
	m := []model.Account{
		m1,
		m2,
	}
	var mj [][]byte

	for _, val := range m {
		j, _ := json.Marshal(val)
		mj = append(mj, j)
	}
	ui.On("GetUser", context.Background(), userID).Return(mj, nil)

	ui2 := new(mockNewStore.NewStore)
	ui2.On("GetUser", mock.Anything, mock.Anything).Return(nil, nil)

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
		{
			name:    "Storage problem",
			stor:    ui2,
			param:   userID,
			want:    []model.Account{},
			wantErr: "storage problem",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewAccService(tc.stor)
			got, err := u.GetUser(context.Background(), tc.param)
			if (err != nil) && err.Error() != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
