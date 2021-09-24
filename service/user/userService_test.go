package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/mocks/pkg/storage/mockNewStore"
	"github.com/stasBigunenko/monorepa/model"
)

type MockLoggingService struct {
}

func (s MockLoggingService) WriteLog(ctx context.Context, message string) {}

func Test_Create(t *testing.T) {
	loggingService := MockLoggingService{}
	ui := new(mockNewStore.NewStore)
	id := uuid.New()
	mm := model.UserHTTP{Name: "Andrew"}
	m := model.UserHTTP{ID: id, Name: "Andrew"}
	ui.On("Create", mock.Anything, mm).Return(m, nil)

	ui2 := new(mockNewStore.NewStore)
	ui2.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("invalid data"))

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
			wantErr: "invalid data",
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
	loggingService := MockLoggingService{}
	ui := new(mockNewStore.NewStore)
	id := uuid.New()
	m := model.UserHTTP{ID: id, Name: "Andrew"}
	ui.On("Get", context.Background(), id).Return(m, nil)

	ui2 := new(mockNewStore.NewStore)
	ui2.On("Get", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))

	tests := []struct {
		name    string
		stor    *mockNewStore.NewStore
		want    model.UserHTTP
		wantErr string
	}{
		{
			name: "Everything ok",
			stor: ui,
			want: m,
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
	loggingService := MockLoggingService{}
	ui := new(mockNewStore.NewStore)
	id := uuid.New()
	ui.On("Delete", context.Background(), id).Return(nil)

	ui2 := new(mockNewStore.NewStore)
	ui2.On("Delete", mock.Anything, mock.Anything).Return(errors.New("not found"))

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
	loggingService := MockLoggingService{}
	ui := new(mockNewStore.NewStore)
	m1 := model.UserHTTP{ID: uuid.New(), Name: "Andrew"}
	m2 := model.UserHTTP{ID: uuid.New(), Name: "Ivan"}
	m := []model.UserHTTP{
		m1,
		m2,
	}
	ui.On("GetAll", context.Background()).Return(m, nil)

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
	loggingService := MockLoggingService{}
	ui := new(mockNewStore.NewStore)
	id := uuid.New()
	m := model.UserHTTP{ID: id, Name: "Abdula"}
	ui.On("Update", context.Background(), m).Return(m, nil)

	ui2 := new(mockNewStore.NewStore)
	ui2.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("not found in DB"))

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
			wantErr: "not found in DB",
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
