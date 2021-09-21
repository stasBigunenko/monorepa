package accountgrpccontroller

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	mocks "github.com/stasBigunenko/monorepa/mocks/pkg/accountGRPC/proto"
	"github.com/stasBigunenko/monorepa/model"
	pb "github.com/stasBigunenko/monorepa/pkg/accountGRPC/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestAccountGRPCСontroller_CreateAccount(t *testing.T) {
	type fields struct {
		client pb.AccountGRPCServiceClient
	}
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uuid.UUID
		wantErr bool
	}{
		{
			name: "CreateAccount OK",
			fields: fields{
				client: mocks.MockAccountGrpcServiceClient{
					MockCreateAccount: func(ctx context.Context, in *pb.UserID, opts ...grpc.CallOption) (*pb.Account, error) {
						it := pb.Account{
							Id:      "00000000-0000-0000-0000-000000000000",
							UserID:  "00000000-0000-0000-0000-000000000000",
							Balance: 0,
						}
						return &it, nil
					},
				},
			},
			args: args{
				userID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			},

			want:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			wantErr: false,
		},
		{
			name: "CreateAccount !OK",
			fields: fields{
				client: mocks.MockAccountGrpcServiceClient{
					MockCreateAccount: func(ctx context.Context, in *pb.UserID, opts ...grpc.CallOption) (*pb.Account, error) {
						return nil, errors.New("err")
					},
				},
			},
			args: args{
				userID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			},

			want:    uuid.Nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := AccountGRPCСontroller{
				client: tt.fields.client,
			}
			got, err := s.CreateAccount(tt.args.userID, context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountGRPCСontroller.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AccountGRPCСontroller.CreateAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountGRPCСontroller_GetAccount(t *testing.T) {
	type fields struct {
		client pb.AccountGRPCServiceClient
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Account
		wantErr bool
	}{
		{
			name: "GetAccount OK",
			fields: fields{
				client: mocks.MockAccountGrpcServiceClient{
					MockGetAccount: func(ctx context.Context, in *pb.AccountID, opts ...grpc.CallOption) (*pb.Account, error) {
						it := pb.Account{
							Id:      "00000000-0000-0000-0000-000000000000",
							UserID:  "00000000-0000-0000-0000-000000000000",
							Balance: 0,
						}
						return &it, nil
					},
				},
			},
			args: args{
				id: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			},

			want: model.Account{
				ID:      uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				UserID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				Balance: 0,
			},
			wantErr: false,
		},
		{
			name: "GetAccount !OK",
			fields: fields{
				client: mocks.MockAccountGrpcServiceClient{
					MockGetAccount: func(ctx context.Context, in *pb.AccountID, opts ...grpc.CallOption) (*pb.Account, error) {
						return &pb.Account{}, errors.New("err")
					},
				},
			},
			args: args{
				id: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			},

			want:    model.Account{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := AccountGRPCСontroller{
				client: tt.fields.client,
			}
			got, err := s.GetAccount(tt.args.id, context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountGRPCСontroller.GetAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AccountGRPCСontroller.GetAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountGRPCСontroller_GetUserAccounts(t *testing.T) {
	type fields struct {
		client pb.AccountGRPCServiceClient
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Account
		wantErr bool
	}{
		{
			name: "GetUserAccounts OK",
			args: args{
				id: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			},
			fields: fields{
				client: mocks.MockAccountGrpcServiceClient{
					MockGetUserAccounts: func(ctx context.Context, in *pb.UserID, opts ...grpc.CallOption) (*pb.AllAccounts, error) {
						it := pb.AllAccounts{
							Accounts: []*pb.Account{
								{
									Id:      "00000000-0000-0000-0000-000000000000",
									UserID:  "00000000-0000-0000-0000-000000000000",
									Balance: 0,
								},
							},
						}
						return &it, nil
					},
				},
			},
			want: []model.Account{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					UserID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					Balance: 0},
			},
			wantErr: false,
		},
		{
			name: "GetUserAccounts !OK",
			args: args{
				id: uuid.Nil,
			},
			fields: fields{
				client: mocks.MockAccountGrpcServiceClient{
					MockGetUserAccounts: func(ctx context.Context, in *pb.UserID, opts ...grpc.CallOption) (*pb.AllAccounts, error) {
						return &pb.AllAccounts{}, errors.New("err")
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := AccountGRPCСontroller{
				client: tt.fields.client,
			}
			got, err := s.GetUserAccounts(tt.args.id, context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountGRPCСontroller.GetAllAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AccountGRPCСontroller.GetAllAccounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountGRPCСontroller_GetAllAccounts(t *testing.T) {
	type fields struct {
		client pb.AccountGRPCServiceClient
	}
	tests := []struct {
		name    string
		fields  fields
		want    []model.Account
		wantErr bool
	}{
		{
			name: "GetAllAccounts OK",
			fields: fields{
				client: mocks.MockAccountGrpcServiceClient{
					MockGetAllUsers: func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.AllAccounts, error) {
						it := pb.AllAccounts{
							Accounts: []*pb.Account{
								{
									Id:      "00000000-0000-0000-0000-000000000000",
									UserID:  "00000000-0000-0000-0000-000000000000",
									Balance: 0,
								},
							},
						}
						return &it, nil
					},
				},
			},
			want: []model.Account{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					UserID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					Balance: 0},
			},
			wantErr: false,
		},
		{
			name: "GetAllAccounts !OK",
			fields: fields{
				client: mocks.MockAccountGrpcServiceClient{
					MockGetAllUsers: func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.AllAccounts, error) {
						return &pb.AllAccounts{}, errors.New("err")
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := AccountGRPCСontroller{
				client: tt.fields.client,
			}
			got, err := s.GetAllAccounts(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountGRPCСontroller.GetAllAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AccountGRPCСontroller.GetAllAccounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountGRPCСontroller_UpdateAccount(t *testing.T) {
	type fields struct {
		client pb.AccountGRPCServiceClient
	}
	type args struct {
		Account model.Account
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "UpdateAccount OK",
			fields: fields{
				client: mocks.MockAccountGrpcServiceClient{
					MockUpdateAccount: func(ctx context.Context, in *pb.Account, opts ...grpc.CallOption) (*pb.Account, error) {
						it := pb.Account{
							Id:      "00000000-0000-0000-0000-000000000000",
							UserID:  "00000000-0000-0000-0000-000000000000",
							Balance: 100,
						}
						return &it, nil
					},
				},
			},
			args: args{
				Account: model.Account{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					UserID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					Balance: 100,
				},
			},
			wantErr: false,
		},
		{
			name: "UpdateAccount !OK",
			fields: fields{
				client: mocks.MockAccountGrpcServiceClient{
					MockUpdateAccount: func(ctx context.Context, in *pb.Account, opts ...grpc.CallOption) (*pb.Account, error) {
						return nil, errors.New("err")
					},
				},
			},
			args: args{
				Account: model.Account{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := AccountGRPCСontroller{
				client: tt.fields.client,
			}
			if err := s.UpdateAccount(tt.args.Account, context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("AccountGRPCСontroller.UpdateAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAccountGRPCСontroller_DeleteAccount(t *testing.T) {
	type fields struct {
		client pb.AccountGRPCServiceClient
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
			name: "DeleteAccount OK",
			fields: fields{
				client: mocks.MockAccountGrpcServiceClient{
					MockDeleteAccount: func(ctx context.Context, in *pb.AccountID, opts ...grpc.CallOption) (*emptypb.Empty, error) {
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
			name: "DeleteAccount !OK",
			fields: fields{
				client: mocks.MockAccountGrpcServiceClient{
					MockDeleteAccount: func(ctx context.Context, in *pb.AccountID, opts ...grpc.CallOption) (*emptypb.Empty, error) {
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
			s := AccountGRPCСontroller{
				client: tt.fields.client,
			}
			if err := s.DeleteAccount(tt.args.id, context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("AccountGRPCСontroller.DeleteAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
