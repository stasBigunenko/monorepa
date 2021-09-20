package mocks

import (
	context "context"

	pb "github.com/stasBigunenko/monorepa/pkg/accountGRPC/proto"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type MockAccountGrpcServiceClient struct {
	MockGetAccount      func(ctx context.Context, in *pb.AccountID, opts ...grpc.CallOption) (*pb.Account, error)
	MockGetUserAccounts func(ctx context.Context, in *pb.UserID, opts ...grpc.CallOption) (*pb.AllAccounts, error)
	MockGetAllUsers     func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.AllAccounts, error)
	MockCreateAccount   func(ctx context.Context, in *pb.UserID, opts ...grpc.CallOption) (*pb.Account, error)
	MockUpdateAccount   func(ctx context.Context, in *pb.Account, opts ...grpc.CallOption) (*pb.Account, error)
	MockDeleteAccount   func(ctx context.Context, in *pb.AccountID, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

func (m MockAccountGrpcServiceClient) GetAccount(ctx context.Context, in *pb.AccountID, opts ...grpc.CallOption) (*pb.Account, error) {
	return m.MockGetAccount(ctx, in, opts...)
}

func (m MockAccountGrpcServiceClient) GetUserAccounts(ctx context.Context, in *pb.UserID, opts ...grpc.CallOption) (*pb.AllAccounts, error) {
	return m.MockGetUserAccounts(ctx, in, opts...)
}

func (m MockAccountGrpcServiceClient) GetAllUsers(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.AllAccounts, error) {
	return m.MockGetAllUsers(ctx, in, opts...)
}

func (m MockAccountGrpcServiceClient) CreateAccount(ctx context.Context, in *pb.UserID, opts ...grpc.CallOption) (*pb.Account, error) {
	return m.MockCreateAccount(ctx, in, opts...)
}

func (m MockAccountGrpcServiceClient) UpdateAccount(ctx context.Context, in *pb.Account, opts ...grpc.CallOption) (*pb.Account, error) {
	return m.MockUpdateAccount(ctx, in, opts...)
}

func (m MockAccountGrpcServiceClient) DeleteAccount(ctx context.Context, in *pb.AccountID, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return m.MockDeleteAccount(ctx, in, opts...)
}
