package grpc

import (
	"context"
	"github.com/stretchr/testify/require"
	"monorepa/service/items"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	monorepa "monorepa/pkg/grpc/proto"
	"monorepa/pkg/storage"
	"net"
	"testing"
)

func TestGrpcServiceServer(t *testing.T) {
	ln := bufconn.Listen(1024)
	dd := storage.NewStorage()
	d := storage.ItemService(dd)
	go serveBufconn(ln, d)
	client := makeBufconnClient(ln)

	x := []storage.Item{
		{Title: "XVLBZGBAICMRAJWW", Description: "  fdzdgrxomvt ler"},
	}

	getItems, _ := client.GetItems(context.Background(), "I")

	require.Equal(t, x[0].Title, getItems[0].Title, "need to fix test program")
	require.Equal(t, x[0].Description, getItems[0].Description, "need to fix test program")

	_, err := client.GetItems(context.Background(), "")
	require.NotNil(t, err, "grpc should return err invalid username")
}

func serveBufconn(ln *bufconn.Listener, data storage.ItemService) {

	s := grpc.NewServer()
	monorepa.RegisterGrpcServiceServer(s, NewGRPC(data))

	_ = s.Serve(ln)
}

func makeBufconnClient(ln *bufconn.Listener) *items.GRPCClient {

	makeBufDialer := func(ln *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
		return func(context.Context, string) (net.Conn, error) {
			return ln.Dial()
		}
	}

	client := &items.GRPCClient{}

	conn, _ := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(makeBufDialer(ln)),
		grpc.WithInsecure(),
	)

	client.Client = monorepa.NewGrpcServiceClient(conn)

	return client
}
