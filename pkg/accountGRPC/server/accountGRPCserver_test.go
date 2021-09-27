package accountgrpcserver

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	mockAccInt "github.com/stasBigunenko/monorepa/mocks/service/account"
	"github.com/stasBigunenko/monorepa/model"
	pb "github.com/stasBigunenko/monorepa/pkg/accountGRPC/proto"
	loggingservice "github.com/stasBigunenko/monorepa/service/loggingService"
)

func Test_Create(t *testing.T) {
	loggingService := loggingservice.New()
	ui := new(mockAccInt.AccInterface)
	uuidS := "00000000-0000-0000-0000-000000000000"
	id, _ := uuid.Parse(uuidS)
	m := model.Account{ID: id, UserID: id, Balance: 0}
	ui.On("Create", mock.Anything, id).Return(m, nil)

	ui2 := new(mockAccInt.AccInterface)
	ui.On("Create", mock.Anything, id).Return(nil, errors.New("err"))

	tests := []struct {
		name    string
		stor    *mockAccInt.AccInterface
		param   *pb.UserID
		want    *pb.Account
		wantErr codes.Code
	}{
		{
			name:  "Everything good",
			param: &pb.UserID{UserID: uuidS},
			stor:  ui,
			want:  &pb.Account{Id: uuidS, UserID: uuidS, Balance: 0},
		},
		{
			name:    "Wrong uuid",
			param:   &pb.UserID{UserID: "000000-0000-0000-0000-000"},
			stor:    ui2,
			want:    nil,
			wantErr: codes.InvalidArgument,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewAccountGRPCServer(tc.stor, loggingService)
			got, err := u.CreateAccount(context.Background(), tc.param)
			if (err != nil) && status.Code(err) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_Get(t *testing.T) {
	loggingService := loggingservice.New()
	ui := new(mockAccInt.AccInterface)
	uuidS := "00000000-0000-0000-0000-000000000000"
	id, _ := uuid.Parse(uuidS)
	m := model.Account{ID: id, UserID: id, Balance: 0}
	ui.On("Get", context.Background(), id).Return(m, nil)

	ui2 := new(mockAccInt.AccInterface)
	m1 := model.UserHTTP{}
	ui.On("Get", mock.Anything, mock.Anything).Return(m1, errors.New("ss"))

	tests := []struct {
		name    string
		stor    *mockAccInt.AccInterface
		param   *pb.AccountID
		want    *pb.Account
		wantErr codes.Code
	}{
		{
			name:  "Everything ok",
			stor:  ui,
			param: &pb.AccountID{Id: "00000000-0000-0000-0000-000000000000"},
			want:  &pb.Account{Id: "00000000-0000-0000-0000-000000000000", UserID: "00000000-0000-0000-0000-000000000000", Balance: 0},
		},
		{
			name:  "Wrong uuid",
			stor:  ui2,
			param: &pb.AccountID{Id: "00000000-0000-00-000000000000"},
			want:  &pb.Account{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewAccountGRPCServer(tc.stor, loggingService)
			got, err := u.GetAccount(context.Background(), tc.param)
			if err != nil && status.Code(err) != tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAccount_Delete(t *testing.T) {
	loggingService := loggingservice.New()
	ui := new(mockAccInt.AccInterface)
	uuidS := "00000000-0000-0000-0000-000000000000"
	id, _ := uuid.Parse(uuidS)
	ui.On("Delete", context.Background(), id).Return(nil)

	tests := []struct {
		name  string
		stor  *mockAccInt.AccInterface
		param *pb.AccountID
		want  *emptypb.Empty
	}{
		{
			name:  "Everything ok",
			stor:  ui,
			param: &pb.AccountID{Id: "00000000-0000-0000-0000-000000000000"},
			want:  &emptypb.Empty{},
		},
		{
			name:  "Wrong id",
			stor:  ui,
			param: &pb.AccountID{Id: "000000000-0000-0000-000000000000"},
			want:  &emptypb.Empty{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewAccountGRPCServer(tc.stor, loggingService)
			got, err := u.DeleteAccount(context.Background(), tc.param)
			if err != nil {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAccount_GetAllUser(t *testing.T) {
	loggingService := loggingservice.New()
	ui := new(mockAccInt.AccInterface)
	uuidS := "00000000-0000-0000-0000-000000000000"
	id, _ := uuid.Parse(uuidS)
	m1 := model.Account{ID: id, UserID: id, Balance: 12}
	m2 := model.Account{ID: id, UserID: id, Balance: 13}
	m := []model.Account{
		m1,
		m2,
	}

	u1 := &pb.Account{Id: uuidS, UserID: uuidS, Balance: 12}
	u2 := &pb.Account{Id: uuidS, UserID: uuidS, Balance: 13}
	all := []*pb.Account{u1, u2}
	aa := pb.AllAccounts{
		Accounts: all,
	}
	ui.On("GetAll", context.Background()).Return(m, nil)

	ui2 := new(mockAccInt.AccInterface)
	ui2.On("GetAll", context.Background()).Return(nil, errors.New("err"))

	tests := []struct {
		name    string
		stor    *mockAccInt.AccInterface
		want    *pb.AllAccounts
		wantErr codes.Code
	}{
		{
			name: "Everything ok",
			stor: ui,
			want: &aa,
		},
		{
			name:    "Error",
			stor:    ui2,
			wantErr: codes.Internal,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewAccountGRPCServer(tc.stor, loggingService)
			got, err := u.GetAllUsers(context.Background(), &emptypb.Empty{})
			if err != nil && status.Code(err) != tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAccount_Update(t *testing.T) {
	loggingService := loggingservice.New()
	ui := new(mockAccInt.AccInterface)
	idd := "00000000-0000-0000-0000-000000000000"
	id, _ := uuid.Parse(idd)
	m := model.Account{ID: id, UserID: id, Balance: 100}
	ui.On("Update", context.Background(), m).Return(m, nil)
	ui2 := new(mockAccInt.AccInterface)
	ui2.On("Update", context.Background(), m).Return(model.Account{}, errors.New("err"))

	tests := []struct {
		name    string
		stor    *mockAccInt.AccInterface
		param   *pb.Account
		want    *pb.Account
		wantErr codes.Code
	}{
		{
			name:  "Everything ok",
			stor:  ui,
			param: &pb.Account{Id: idd, UserID: idd, Balance: 100},
			want:  &pb.Account{Id: idd, UserID: idd, Balance: 100},
		},
		{
			name:    "Error",
			stor:    ui2,
			param:   &pb.Account{Id: idd, UserID: idd, Balance: 100},
			wantErr: codes.Internal,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewAccountGRPCServer(tc.stor, loggingService)
			got, err := u.UpdateAccount(context.Background(), tc.param)
			if (err != nil) && status.Code(err) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAccount_GetUserAccount(t *testing.T) {
	loggingService := loggingservice.New()
	ui := new(mockAccInt.AccInterface)
	uuidS := "00000000-0000-0000-0000-000000000000"
	id, _ := uuid.Parse(uuidS)
	m1 := model.Account{ID: id, UserID: id, Balance: 12}
	m2 := model.Account{ID: id, UserID: id, Balance: 13}
	m := []model.Account{
		m1,
		m2,
	}

	u1 := &pb.Account{Id: uuidS, UserID: uuidS, Balance: 12}
	u2 := &pb.Account{Id: uuidS, UserID: uuidS, Balance: 13}
	all := []*pb.Account{u1, u2}
	aa := pb.AllAccounts{
		Accounts: all,
	}
	ui.On("GetUser", context.Background(), id).Return(m, nil)

	ui2 := new(mockAccInt.AccInterface)
	ui2.On("GetUser", context.Background(), id).Return(nil, errors.New("err"))

	tests := []struct {
		name    string
		stor    *mockAccInt.AccInterface
		param   *pb.UserID
		want    *pb.AllAccounts
		wantErr codes.Code
	}{
		{
			name:  "Everything ok",
			stor:  ui,
			param: &pb.UserID{UserID: uuidS},
			want:  &aa,
		},
		{
			name:    "Wrong uuid",
			param:   &pb.UserID{UserID: "000000-0000-0000-0000-000"},
			stor:    ui2,
			want:    nil,
			wantErr: codes.InvalidArgument,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewAccountGRPCServer(tc.stor, loggingService)
			got, err := u.GetUserAccounts(context.Background(), tc.param)
			if err != nil && status.Code(err) != tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
