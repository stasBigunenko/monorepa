package items

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "monorepa/pkg/grpc/proto"
	"monorepa/pkg/storage"
)

type gRPCClient struct {
	client pb.GrpcServiceClient
}

func New(client pb.GrpcServiceClient) gRPCClient {
	return gRPCClient{
		client: client,
	}
}

func (gc gRPCClient) GetItems(un string) ([]storage.StorageInterface, error) {
	data, err := gc.client.GetItems(context.Background(), &pb.Username{
		Username: un,
	})
	if err != nil {
		return []storage.StorageInterface{}, status.Error(codes.Internal, "internal problem")
	}
	itemsAll := []storage.StorageInterface{}

	for _, val := range data.Items {
		itemsAll = append(itemsAll, &storage.Storage{
			Id:          val.Id,
			Title:       val.Title,
			Description: val.Description,
		})
	}
	return itemsAll, nil
}
