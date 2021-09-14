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

func New(addr string) *GRPCClient {
	gc := &GRPCClient{}

	conn, _ := grpc.Dial(addr, grpc.WithInsecure())

	gc.Client = pb.NewGrpcServiceClient(conn)

	return gc
}

func (gc *GRPCClient) GetItems(_ context.Context, un string) ([]storage.Item, error) {

	data, err := gc.Client.GetItems(context.Background(), &pb.Username{
		Username: un,
	})
	if err != nil {
		return []storage.Item{}, status.Error(codes.Internal, "internal problem")
	}

	itemsAll := []storage.Item{}

	for _, val := range data.Items {
		itemsAll = append(itemsAll, storage.Item{
			ID:          val.Id,
			Title:       val.Title,
			Description: val.Description,
		})
	}
	return itemsAll, nil
}
