package user

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/mocks/pkg/storage/mockNewStore"
	"github.com/stasBigunenko/monorepa/model"
	loggingservice "github.com/stasBigunenko/monorepa/service/loggingService"
)

func Test_Create(t *testing.T) {
	loggingService := loggingservice.New()
	ui := new(mockNewStore.NewStore)
	id := uuid.New()
	mm := model.UserHTTP{Name: "Andrew"}
	mj, _ := json.Marshal(mm)
	m := model.UserHTTP{ID: id, Name: "Andrew"}
	ui.On("Create", mock.Anything, mj).Return(mj, id, nil)

	ui2 := new(mockNewStore.NewStore)
	ui2.On("Create", mock.Anything, mock.Anything).Return(nil, nil, nil)

	tests := []struct {
		name    string
		param   string
		stor    *mockNewStore.NewStore
		want    model.UserHTTP
		wantErr string
	}{
		{
			name:  "Everything good",
			param: "Andrew",
			stor:  ui,
			want:  m,
		},
		{
			name:    "Everything bad",
			param:   "sdda",
			stor:    ui2,
			want:    model.UserHTTP{},
			wantErr: "unmarshal problem",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewUsrService(tc.stor, loggingService)
			got, err := u.Create(context.Background(), tc.param)
			if (err != nil) && err.Error() != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_Get(t *testing.T) {
	loggingService := loggingservice.New()
	ui := new(mockNewStore.NewStore)
	id := uuid.New()
	m := model.UserHTTP{ID: id, Name: "Andrew"}
	mj, _ := json.Marshal(m)
	ui.On("Get", context.Background(), id).Return(mj, nil)

	ui2 := new(mockNewStore.NewStore)
	ui2.On("Get", mock.Anything, mock.Anything).Return(nil, nil)

	tests := []struct {
		name    string
		stor    *mockNewStore.NewStore
		want    model.UserHTTP
		wantErr string
	}{
		{
			name: "Everything ok",
			stor: ui,
			want: model.UserHTTP{
				ID:   id,
				Name: "Andrew",
			},
		},
		{
			name:    "Everything bad",
			stor:    ui2,
			want:    model.UserHTTP{},
			wantErr: "not found",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewUsrService(tc.stor, loggingService)
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
	loggingService := loggingservice.New()
	ui := new(mockNewStore.NewStore)
	id := uuid.New()
	ui.On("Delete", context.Background(), id).Return(true, nil)

	ui2 := new(mockNewStore.NewStore)
	ui2.On("Delete", mock.Anything, mock.Anything).Return(false, nil)

	tests := []struct {
		name    string
		stor    *mockNewStore.NewStore
		param   uuid.UUID
		result  []byte
		wantErr string
	}{
		{
			name:  "Everything ok",
			stor:  ui,
			param: id,
		},
		{
			name:    "false",
			stor:    ui2,
			wantErr: "not found",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewUsrService(tc.stor, loggingService)
			err := u.Delete(context.Background(), tc.param)
			if err != nil {
				assert.Error(t, err)
				return
			}
			assert.Empty(t, err)
		})
	}
}

//
func TestUserService_GetAll(t *testing.T) {
	loggingService := loggingservice.New()
	ui := new(mockNewStore.NewStore)
	m1 := model.UserHTTP{ID: uuid.New(), Name: "Andrew"}
	m2 := model.UserHTTP{ID: uuid.New(), Name: "Ivan"}
	m := []model.UserHTTP{
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
		want    []model.UserHTTP
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
			u := NewUsrService(tc.stor, loggingService)
			got, err := u.GetAll(context.Background())
			if (err != nil) && err.Error() != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestUserService_Update(t *testing.T) {
	loggingService := loggingservice.New()
	ui := new(mockNewStore.NewStore)
	id := uuid.New()
	m := model.UserHTTP{ID: id, Name: "Abdula"}
	mj, _ := json.Marshal(m)
	ui.On("Update", context.Background(), id, mj).Return(mj, nil)

	ui2 := new(mockNewStore.NewStore)
	ui2.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	tests := []struct {
		name    string
		stor    *mockNewStore.NewStore
		param   model.UserHTTP
		want    model.UserHTTP
		wantErr string
	}{
		{
			name:  "Everything ok",
			stor:  ui,
			param: m,
			want:  m,
		},
		{
			name:    "Not found",
			stor:    ui2,
			param:   m,
			want:    model.UserHTTP{},
			wantErr: "not found",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewUsrService(tc.stor, loggingService)
			got, err := u.Update(context.Background(), m)
			if (err != nil) && err.Error() != tc.wantErr {
				t.Errorf("SomeLogic error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
