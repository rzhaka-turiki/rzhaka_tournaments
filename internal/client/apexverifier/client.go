package apexverifier

import (
	"context"

	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/client/apexverifier/proto/apexpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	VerifyAccount(ctx context.Context, player string, platform string, level int32) (*VerifyAccountResult, error)
	Close() error
}

type VerifyAccountResult struct {
	UID      string
	Player   string
	Platform string
	Level    int32
}

type client struct {
	conn   *grpc.ClientConn
	client apexpb.ApexVerifierClient
}

func NewClient(addr string) (Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &client{
		conn:   conn,
		client: apexpb.NewApexVerifierClient(conn),
	}, nil
}

func (c *client) VerifyAccount(ctx context.Context, player string, platform string, level int32) (*VerifyAccountResult, error) {
	resp, err := c.client.VerifyAccount(ctx, &apexpb.VerifyAccountRequest{
		Player:   player,
		Platform: platform,
		Level:    level,
	})
	if err != nil {
		return nil, err
	}
	return &VerifyAccountResult{
		UID:      resp.Uid,
		Player:   resp.Player,
		Platform: resp.Platform,
		Level:    resp.Level,
	}, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}
