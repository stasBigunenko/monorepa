package usergrpcserver

import (
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	userInt "github.com/stasBigunenko/monorepa/mocks/service/user"
	"github.com/stasBigunenko/monorepa/model"
	pb "github.com/stasBigunenko/monorepa/pkg/userGRPC/proto"
	loggingservice "github.com/stasBigunenko/monorepa/service/loggingService"
)

type MockLoggingService struct {
}

func (MockLoggingService) WriteLog(ctx context.Context, message string) {}

func Test_Create(t *testing.T) {
	loggingService := MockLoggingService{}
	ui := new(userInt.User)
	s := "Andrew"
	uuidS := "00000000-0000-0000-0000-000000000000"
	id, _ := uuid.Parse(uuidS)
	m := model.UserHTTP{ID: id, Name: "Andrew"}
	ui.On("Create", mock.Anything, s).Return(m, nil)

	tests := []struct {
		name    string
		stor    *userInt.User
		param   *pb.Name
		want    *pb.User
		wantErr codes.Code
	}{
		{
			name:  "Everything good",
			param: &pb.Name{Name: "Andrew"},
			stor:  ui,
			want:  &pb.User{Id: "00000000-0000-0000-0000-000000000000", Name: "Andrew"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewUsersGRPCServer(tc.stor, loggingService)
			got, err := u.Create(context.Background(), tc.param)
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
	ui := new(userInt.User)
	uuidS := "00000000-0000-0000-0000-000000000000"
	id, _ := uuid.Parse(uuidS)
	m := model.UserHTTP{ID: id, Name: "Andrew"}
	ui.On("Get", context.Background(), id).Return(m, nil)

	ui2 := new(userInt.User)
	m1 := model.UserHTTP{}
	ui.On("Get", mock.Anything, mock.Anything).Return(m1, errors.New("ss"))

	tests := []struct {
		name    string
		stor    *userInt.User
		param   *pb.Id
		want    *pb.User
		wantErr codes.Code
	}{
		{
			name:  "Everything ok",
			stor:  ui,
			param: &pb.Id{Id: "00000000-0000-0000-0000-000000000000"},
			want:  &pb.User{Id: "00000000-0000-0000-0000-000000000000", Name: "Andrew"},
		},
		{
			name:  "Wrong uuid",
			stor:  ui2,
			param: &pb.Id{Id: "00000000-0000-00-000000000000"},
			want:  &pb.User{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewUsersGRPCServer(tc.stor, loggingService)
			ctx := (context.Background())
			contextID := "test"
			c := metadata.AppendToOutgoingContext(ctx, "requestid", contextID)
			got, err := u.Get(c, tc.param)
			if err != nil && status.Code(err) != tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	loggingService := loggingservice.New()
	ui := new(userInt.User)
	uuidS := "00000000-0000-0000-0000-000000000000"
	id, _ := uuid.Parse(uuidS)
	ui.On("Delete", context.Background(), id).Return(nil)

	tests := []struct {
		name  string
		stor  *userInt.User
		param *pb.Id
		want  *emptypb.Empty
	}{
		{
			name:  "Everything ok",
			stor:  ui,
			param: &pb.Id{Id: "00000000-0000-0000-0000-000000000000"},
			want:  &emptypb.Empty{},
		},
		{
			name:  "Wrong id",
			stor:  ui,
			param: &pb.Id{Id: "000000000-0000-0000-000000000000"},
			want:  &emptypb.Empty{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewUsersGRPCServer(tc.stor, loggingService)
			got, err := u.Delete(context.Background(), tc.param)
			if err != nil {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

//
func TestUserService_GetAllUsers(t *testing.T) {
	loggingService := loggingservice.New()
	ui := new(userInt.User)
	id1 := uuid.New()
	id2 := uuid.New()
	m1 := model.UserHTTP{ID: id1, Name: "Andrew"}
	m2 := model.UserHTTP{ID: id2, Name: "Ivan"}
	m := []model.UserHTTP{
		m1,
		m2,
	}

	u1 := &pb.User{Id: id1.String(), Name: m1.Name}
	u2 := &pb.User{Id: id2.String(), Name: m2.Name}
	all := []*pb.User{u1, u2}
	au := pb.AllUsers{
		AllUsers: all,
	}
	ui.On("GetAll", context.Background()).Return(m, nil)

	ui2 := new(userInt.User)
	ui2.On("GetAll", context.Background()).Return(nil, errors.New("err"))

	tests := []struct {
		name    string
		stor    *userInt.User
		want    *pb.AllUsers
		wantErr codes.Code
	}{
		{
			name: "Everything ok",
			stor: ui,
			want: &au,
		},
		{
			name:    "Error",
			stor:    ui2,
			wantErr: codes.Internal,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewUsersGRPCServer(tc.stor, loggingService)
			got, err := u.GetAllUsers(context.Background(), &emptypb.Empty{})
			if err != nil && status.Code(err) != tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

//
func TestUserService_Update(t *testing.T) {
	loggingService := loggingservice.New()
	ui := new(userInt.User)
	idd := "00000000-0000-0000-0000-000000000000"
	id, _ := uuid.Parse(idd)
	m := model.UserHTTP{ID: id, Name: "Abdula"}
	ui.On("Update", context.Background(), m).Return(m, nil)
	ui2 := new(userInt.User)
	ui2.On("Update", context.Background(), m).Return(model.UserHTTP{}, errors.New("err"))

	tests := []struct {
		name    string
		stor    *userInt.User
		param   *pb.User
		want    *pb.User
		wantErr codes.Code
	}{
		{
			name:  "Everything ok",
			stor:  ui,
			param: &pb.User{Id: idd, Name: "Abdula"},
			want:  &pb.User{Id: idd, Name: "Abdula"},
		},
		{
			name:    "Error",
			stor:    ui2,
			param:   &pb.User{Id: idd, Name: "Abdula"},
			wantErr: codes.Internal,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewUsersGRPCServer(tc.stor, loggingService)
			got, err := u.Update(context.Background(), tc.param)
			if (err != nil) && status.Code(err) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
