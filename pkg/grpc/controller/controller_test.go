package grpccontroller

import (
	"context"
	"fmt"
	pb "github.com/stasBigunenko/monorepa/pkg/grpc/proto"
	"reflect"
	"testing"

	"google.golang.org/grpc"
)

type MockGrpcServiceClient struct {
	MockGetItems func(ctx context.Context, in *pb.Username, opts ...grpc.CallOption) (*pb.Items, error)
}

func (m MockGrpcServiceClient) GetItems(ctx context.Context, in *pb.Username, opts ...grpc.CallOption) (*pb.Items, error) {
	return m.MockGetItems(ctx, in, opts...)
}

func TestGRPCСontroller_GetItems(t *testing.T) {
	type fields struct {
		client pb.GrpcServiceClient
	}
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Item
		wantErr bool
	}{
		{
			name: "GetItems OK",
			fields: fields{
				client: MockGrpcServiceClient{
					MockGetItems: func(ctx context.Context, in *pb.Username, opts ...grpc.CallOption) (*pb.Items, error) {
						it := pb.Items{
							Items: []*pb.Obj{
								{
									Id:    "00000000-0000-0000-0000-000000000000",
									Title: "Vasya",
								},
							},
						}
						return &it, nil
					},
				},
			},
			args: args{
				username: "Boris",
			},
			want: []Item{
				{
					Name: "Vasya",
				},
			},
			wantErr: false,
		},
		{
			name: "GetItems wrong uuid",
			fields: fields{
				client: MockGrpcServiceClient{
					MockGetItems: func(ctx context.Context, in *pb.Username, opts ...grpc.CallOption) (*pb.Items, error) {
						it := pb.Items{
							Items: []*pb.Obj{
								{
									Id:    "sdfsdfsdf",
									Title: "Vasya",
								},
							},
						}
						return &it, nil
					},
				},
			},
			args: args{
				username: "Boris",
			},
			wantErr: true,
		},
		{
			name: "GetItems internal server error",
			fields: fields{
				client: MockGrpcServiceClient{
					MockGetItems: func(ctx context.Context, in *pb.Username, opts ...grpc.CallOption) (*pb.Items, error) {
						return nil, fmt.Errorf("Internal server error")
					},
				},
			},
			args: args{
				username: "Boris",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := GRPCСontroller{
				client: tt.fields.client,
			}
			got, err := s.GetItems(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCСontroller.GetItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCСontroller.GetItems() = %v, want %v", got, tt.want)
			}
		})
	}
}
