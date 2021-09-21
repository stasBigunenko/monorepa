package mocks

import (
	context "context"

	pb "github.com/stasBigunenko/monorepa/pkg/userGRPC/proto"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type MockUserGrpcServiceClient struct {
	MockCreate      func(ctx context.Context, in *pb.Name, opts ...grpc.CallOption) (*pb.User, error)
	MockGet         func(ctx context.Context, in *pb.Id, opts ...grpc.CallOption) (*pb.User, error)
	MockDelete      func(ctx context.Context, in *pb.Id, opts ...grpc.CallOption) (*emptypb.Empty, error)
	MockGetAllUsers func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.AllUsers, error)
	MockUpdate      func(ctx context.Context, in *pb.User, opts ...grpc.CallOption) (*pb.User, error)
}

func (m MockUserGrpcServiceClient) Create(ctx context.Context, in *pb.Name, opts ...grpc.CallOption) (*pb.User, error) {
	return m.MockCreate(ctx, in, opts...)
}

func (m MockUserGrpcServiceClient) Get(ctx context.Context, in *pb.Id, opts ...grpc.CallOption) (*pb.User, error) {
	return m.MockGet(ctx, in, opts...)
}

func (m MockUserGrpcServiceClient) Delete(ctx context.Context, in *pb.Id, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return m.MockDelete(ctx, in, opts...)
}

func (m MockUserGrpcServiceClient) GetAllUsers(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.AllUsers, error) {
	return m.MockGetAllUsers(ctx, in, opts...)
}

func (m MockUserGrpcServiceClient) Update(ctx context.Context, in *pb.User, opts ...grpc.CallOption) (*pb.User, error) {
	return m.MockUpdate(ctx, in, opts...)
}
