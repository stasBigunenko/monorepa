package grpccontroller

import (
	"context"
	"fmt"

	pb "monorepa/pkg/grpc/proto"

	"github.com/google/uuid"
)

type GRPCСontroller struct {
	client pb.GrpcServiceClient
}

func New(cli pb.GrpcServiceClient) GRPCСontroller {
	return GRPCСontroller{
		client: cli,
	}
}

type Item struct {
	ID       uuid.UUID
	Name     string
	Comments string
}

func (s GRPCСontroller) GetItems(username string) ([]Item, error) {
	resp, err := s.client.GetItems(context.Background(), &pb.Username{
		Username: username,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get items from GRPC-server: %w", err)
	}

	items := []Item{}
	for _, obj := range resp.Items {
		id, err := uuid.Parse(obj.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to parse uuid from string: %w", err)
		}

		items = append(items, Item{
			ID:       id,
			Name:     obj.Title,
			Comments: obj.Description,
		})
	}

	return items, nil
}
