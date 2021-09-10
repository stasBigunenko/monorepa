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

func New() *GRPCClient {
	gc := &GRPCClient{}

	conn, _ := grpc.Dial(":8080", grpc.WithInsecure())

	gc.Client = pb.NewGrpcServiceClient(conn)

	return gc
}

func (gc *GRPCClient) GetItems(ctx context.Context, un string) ([]storage.Storage, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	data, err := gc.Client.GetItems(context.Background(), &pb.Username{
		Username: un,
	})
	if err != nil {
		return []storage.Storage{}, status.Error(codes.Internal, "internal problem")
	}

	itemsAll := []storage.Storage{}

	for _, val := range data.Items {
		itemsAll = append(itemsAll, storage.Storage{
			Id:          val.Id,
			Title:       val.Title,
			Description: val.Description,
		})
	}
	return itemsAll, nil
}
