package newStorage

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStorageDB_Get(t *testing.T) {
	acc := NewDB()
	id := uuid.New()
	userID := uuid.New()
	m := model.Account{ID: id, UserID: userID, Balance: 0}
	j, _ := json.Marshal(m)
	acc.Data[id] = j
	nf := uuid.New()

	tests := []struct {
		name  string
		param uuid.UUID
		want  []byte
	}{
		{
			name:  "Everything ok",
			param: id,
			want:  j,
		},
		{
			name:  "not found",
			param: nf,
			want:  nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := acc.Get(context.Background(), tc.param)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStorageDB_GetUser(t *testing.T) {
	acc := NewDB()
	id := uuid.New()
	userID := uuid.New()
	m := model.Account{ID: id, UserID: userID, Balance: 0}
	j, _ := json.Marshal(m)
	acc.Data[id] = j
	res := [][]byte{
		j,
	}

	tests := []struct {
		name string
		want [][]byte
	}{
		{
			name: "Everything ok",
			want: res,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := acc.GetUser(context.Background(), userID)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStorageDB_GetAll(t *testing.T) {
	acc := NewDB()
	id := uuid.New()
	userID := uuid.New()
	m := model.Account{ID: id, UserID: userID, Balance: 0}
	j, _ := json.Marshal(m)
	acc.Data[id] = j
	res := [][]byte{
		j,
	}

	tests := []struct {
		name string
		want [][]byte
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
	acc := NewDB()
	username := "Jim"
	r, _ := json.Marshal(username)

	tests := []struct {
		name  string
		param []byte
		want  [][]byte
	}{
		{
			name:  "Everything ok",
			param: r,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			val, _, _ := acc.Create(context.Background(), tc.param)
			require.NotNil(t, val, "should be no error")
		})
	}
}

func TestStorageDB_Update(t *testing.T) {
	acc := NewDB()
	id := uuid.New()
	userID := uuid.New()
	m := model.Account{ID: id, UserID: userID, Balance: 0}
	j, _ := json.Marshal(m)
	acc.Data[id] = j
	m.Balance = 10
	j2, _ := json.Marshal(m)

	tests := []struct {
		name   string
		param1 uuid.UUID
		param2 []byte
		want   []byte
	}{
		{
			name:   "Everything ok",
			param1: id,
			param2: j2,
			want:   j2,
		},
		{
			name:   "Not forun (error)",
			param1: uuid.New(),
			param2: j2,
			want:   nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := acc.Update(context.Background(), tc.param1, tc.param2)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStorageDB_Delete(t *testing.T) {
	acc := NewDB()
	id := uuid.New()
	userID := uuid.New()
	m := model.Account{ID: id, UserID: userID, Balance: 0}
	j, _ := json.Marshal(m)
	acc.Data[id] = j

	tests := []struct {
		name  string
		param uuid.UUID
		want  bool
	}{
		{
			name:  "Everything ok",
			param: id,
			want:  true,
		},
		{
			name:  "not found",
			param: id,
			want:  false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := acc.Delete(context.Background(), id)
			assert.Equal(t, tc.want, got)
		})
	}
}
