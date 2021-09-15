package grpccontroller

import (
	"context"
	"fmt"

	"monorepa/model"
	pb "github.com/stasBigunenko/monorepa/pkg/grpc/proto"

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

func (s GRPCСontroller) GetItems(username string) ([]model.Item, error) {
	resp, err := s.client.GetItems(context.Background(), &pb.Username{
		Username: username,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get items from GRPC-server: %w", err)
	}

	items := []model.Item{}
	for _, obj := range resp.Items {
		id, err := uuid.Parse(obj.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to parse uuid from string: %w", err)
		}

		items = append(items, model.Item{
			ID:          id,
			Title:       obj.Title,
			Description: obj.Description,
		})
	}

	return items, nil
}
