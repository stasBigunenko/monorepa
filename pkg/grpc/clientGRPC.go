package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "github.com/stasBigunenko/monorepa/pkg/grpc/proto"
	"github.com/stasBigunenko/monorepa/pkg/storage"
)

// Client GRPC

type gRPCClient struct {
	client pb.GrpcServiceClient
}

func (gc *gRPCClient) getItems(_ context.Context, un string) ([]storage.Item, error) {

	data, err := gc.client.GetItems(context.Background(), &pb.Username{
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
