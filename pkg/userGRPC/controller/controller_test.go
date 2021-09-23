package usergrpccontroller

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	mocks "github.com/stasBigunenko/monorepa/mocks/pkg/userGRPC/proto"
	"github.com/stasBigunenko/monorepa/model"
	pb "github.com/stasBigunenko/monorepa/pkg/userGRPC/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MockLoggingService struct {
}

func (MockLoggingService) WriteLog(ctx context.Context, message string) {}

func TestUserGRPCСontroller_CreateUser(t *testing.T) {
	type fields struct {
		client         pb.UserGRPCServiceClient
		loggingService LoggingService
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uuid.UUID
		wantErr bool
	}{
		{
			name: "CreateUser OK",
			fields: fields{
				client: mocks.MockUserGrpcServiceClient{
					MockCreate: func(ctx context.Context, in *pb.Name, opts ...grpc.CallOption) (*pb.User, error) {
						it := pb.User{
							Id:   "00000000-0000-0000-0000-000000000000",
							Name: in.Name,
						}
						return &it, nil
					},
				},
			},
			args: args{
				name: "Boris",
			},

			want:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			wantErr: false,
		},
		{
			name: "CreateUser !OK",
			fields: fields{
				client: mocks.MockUserGrpcServiceClient{
					MockCreate: func(ctx context.Context, in *pb.Name, opts ...grpc.CallOption) (*pb.User, error) {
						return nil, errors.New("err")
					},
				},
			},
			args: args{
				name: "Boris",
			},

			want:    uuid.Nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := UserGRPCСontroller{
				client:         tt.fields.client,
				loggingService: MockLoggingService{},
			}
			got, err := s.CreateUser(context.Background(), tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserGRPCСontroller.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserGRPCСontroller.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserGRPCСontroller_GetUser(t *testing.T) {
	type fields struct {
		client         pb.UserGRPCServiceClient
		loggingService LoggingService
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.UserHTTP
		wantErr bool
	}{
		{
			name: "GetUser OK",
			fields: fields{
				client: mocks.MockUserGrpcServiceClient{
					MockGet: func(ctx context.Context, in *pb.Id, opts ...grpc.CallOption) (*pb.User, error) {
						it := pb.User{
							Id:   "00000000-0000-0000-0000-000000000000",
							Name: "Boris",
						}
						return &it, nil
					},
				},
			},
			args: args{
				id: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			},

			want: model.UserHTTP{
				ID:   uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				Name: "Boris",
			},
			wantErr: false,
		},
		{
			name: "GetUser !OK",
			fields: fields{
				client: mocks.MockUserGrpcServiceClient{
					MockGet: func(ctx context.Context, in *pb.Id, opts ...grpc.CallOption) (*pb.User, error) {
						return &pb.User{}, errors.New("err")
					},
				},
			},
			args: args{
				id: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			},

			want:    model.UserHTTP{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := UserGRPCСontroller{
				client:         tt.fields.client,
				loggingService: MockLoggingService{},
			}
			got, err := s.GetUser(context.Background(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserGRPCСontroller.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserGRPCСontroller.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserGRPCСontroller_GetAllUsers(t *testing.T) {
	type fields struct {
		client         pb.UserGRPCServiceClient
		loggingService LoggingService
	}
	tests := []struct {
		name    string
		fields  fields
		want    []model.UserHTTP
		wantErr bool
	}{
		{
			name: "GetAllUsers OK",
			fields: fields{
				client: mocks.MockUserGrpcServiceClient{
					MockGetAllUsers: func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.AllUsers, error) {
						it := pb.AllUsers{
							AllUsers: []*pb.User{
								{
									Id:   "00000000-0000-0000-0000-000000000000",
									Name: "Boris",
								},
							},
						}
						return &it, nil
					},
				},
			},
			want: []model.UserHTTP{
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					Name: "Boris",
				},
			},
			wantErr: false,
		},
		{
			name: "GetAllUSers !OK",
			fields: fields{
				client: mocks.MockUserGrpcServiceClient{
					MockGetAllUsers: func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.AllUsers, error) {
						return &pb.AllUsers{}, errors.New("err")
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := UserGRPCСontroller{
				client:         tt.fields.client,
				loggingService: MockLoggingService{},
			}
			got, err := s.GetAllUsers(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("UserGRPCСontroller.GetAllUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserGRPCСontroller.GetAllUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserGRPCСontroller_UpdateUser(t *testing.T) {
	type fields struct {
		client         pb.UserGRPCServiceClient
		loggingService LoggingService
	}
	type args struct {
		user model.UserHTTP
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "UpdateUser OK",
			fields: fields{
				client: mocks.MockUserGrpcServiceClient{
					MockUpdate: func(ctx context.Context, in *pb.User, opts ...grpc.CallOption) (*pb.User, error) {
						it := pb.User{
							Id:   "00000000-0000-0000-0000-000000000000",
							Name: "Boris",
						}
						return &it, nil
					},
				},
			},
			args: args{
				user: model.UserHTTP{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					Name: "Jake",
				},
			},
			wantErr: false,
		},
		{
			name: "UpdateUser !OK",
			fields: fields{
				client: mocks.MockUserGrpcServiceClient{
					MockUpdate: func(ctx context.Context, in *pb.User, opts ...grpc.CallOption) (*pb.User, error) {
						return nil, errors.New("err")
					},
				},
			},
			args: args{
				user: model.UserHTTP{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := UserGRPCСontroller{
				client:         tt.fields.client,
				loggingService: MockLoggingService{},
			}
			if err := s.UpdateUser(context.Background(), tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("UserGRPCСontroller.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserGRPCСontroller_DeleteUser(t *testing.T) {
	type fields struct {
		client         pb.UserGRPCServiceClient
		loggingService LoggingService
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "DeleteUser OK",
			fields: fields{
				client: mocks.MockUserGrpcServiceClient{
					MockDelete: func(ctx context.Context, in *pb.Id, opts ...grpc.CallOption) (*emptypb.Empty, error) {
						return &emptypb.Empty{}, nil
					},
				},
			},
			args: args{
				id: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			},
			wantErr: false,
		},
		{
			name: "DeleteUser !OK",
			fields: fields{
				client: mocks.MockUserGrpcServiceClient{
					MockDelete: func(ctx context.Context, in *pb.Id, opts ...grpc.CallOption) (*emptypb.Empty, error) {
						return &emptypb.Empty{}, errors.New("err")
					},
				},
			},
			args: args{
				id: uuid.Nil,
			},

			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := UserGRPCСontroller{
				client:         tt.fields.client,
				loggingService: MockLoggingService{},
			}
			if err := s.DeleteUser(context.Background(), tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("UserGRPCСontroller.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
