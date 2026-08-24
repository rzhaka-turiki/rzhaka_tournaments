package matchapi

import (
	"context"

	matchapipb "github.com/rzhaka-turiki/rzhaka_tournaments/internal/client/matchapi/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client interface {
	ListMatches(ctx context.Context, req *matchapipb.ListMatchesRequest) (*matchapipb.ListMatchesResponse, error)
	GetMatchDetail(ctx context.Context, req *matchapipb.MatchDetailRequest) (*matchapipb.MatchDetailResponse, error)
	AddToken(ctx context.Context, req *matchapipb.AddTokenRequest) (*matchapipb.AddTokenResponse, error)
	ListTokens(ctx context.Context, req *matchapipb.ListTokensRequest) (*matchapipb.ListTokensResponse, error)
	GetMatchStats(ctx context.Context, req *matchapipb.MatchStatsRequest) (*matchapipb.MatchStatsResponse, error)
	GetPlayerStats(ctx context.Context, req *matchapipb.PlayerStatsRequest) (*matchapipb.PlayerStatsResponse, error)
}

type client struct {
	conn   *grpc.ClientConn
	client matchapipb.MatchServiceClient
	apiKey string
}

func NewClient(addr, apiKey string) (Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &client{
		conn:   conn,
		client: matchapipb.NewMatchServiceClient(conn),
		apiKey: apiKey,
	}, nil
}

func (c *client) contextWithAuth(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(
		ctx,
		"authorization",
		"Bearer "+c.apiKey,
	)
}

func (c *client) AddToken(ctx context.Context, req *matchapipb.AddTokenRequest) (*matchapipb.AddTokenResponse, error) {
	return c.client.AddToken(ctx, req)
}

func (c *client) ListTokens(ctx context.Context, req *matchapipb.ListTokensRequest) (*matchapipb.ListTokensResponse, error) {
	return c.client.ListTokens(ctx, req)
}

func (c *client) ListMatches(ctx context.Context, req *matchapipb.ListMatchesRequest) (*matchapipb.ListMatchesResponse, error) {
	return c.client.ListMatches(ctx, req)
}

func (c *client) GetMatchStats(ctx context.Context, req *matchapipb.MatchStatsRequest) (*matchapipb.MatchStatsResponse, error) {
	return c.client.GetMatchStats(ctx, req)
}

func (c *client) GetMatchDetail(ctx context.Context, req *matchapipb.MatchDetailRequest) (*matchapipb.MatchDetailResponse, error) {
	return c.client.GetMatchDetail(ctx, req)
}

func (c *client) GetPlayerStats(ctx context.Context, req *matchapipb.PlayerStatsRequest) (*matchapipb.PlayerStatsResponse, error) {
	return c.client.GetPlayerStats(ctx, req)
}
