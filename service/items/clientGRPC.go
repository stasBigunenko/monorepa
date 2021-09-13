package items

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "monorepa/pkg/grpc/proto"
	"monorepa/pkg/storage"
)

// Client GRPC

type GRPCClient struct {
	Client pb.GrpcServiceClient
}

<<<<<<< HEAD
func New() *GRPCClient {
	gc := &GRPCClient{}

	conn, _ := grpc.Dial(":8080", grpc.WithInsecure())
=======
func New(addr string) *GRPCClient {
	gc := &GRPCClient{}

	conn, _ := grpc.Dial(addr, grpc.WithInsecure())
>>>>>>> gRPC_protofile

	gc.Client = pb.NewGrpcServiceClient(conn)

	return gc
}

<<<<<<< HEAD
func (gc *GRPCClient) GetItems(ctx context.Context, un string) ([]storage.Storage, error) {
=======
func (gc *GRPCClient) GetItems(ctx context.Context, un string) ([]storage.StorageItem, error) {
>>>>>>> gRPC_protofile
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	data, err := gc.Client.GetItems(context.Background(), &pb.Username{
		Username: un,
	})
	if err != nil {
<<<<<<< HEAD
		return []storage.Storage{}, status.Error(codes.Internal, "internal problem")
	}

	itemsAll := []storage.Storage{}

	for _, val := range data.Items {
		itemsAll = append(itemsAll, storage.Storage{
=======
		return []storage.StorageItem{}, status.Error(codes.Internal, "internal problem")
	}

	itemsAll := []storage.StorageItem{}

	for _, val := range data.Items {
		itemsAll = append(itemsAll, storage.StorageItem{
>>>>>>> gRPC_protofile
			Id:          val.Id,
			Title:       val.Title,
			Description: val.Description,
		})
	}
	return itemsAll, nil
}
