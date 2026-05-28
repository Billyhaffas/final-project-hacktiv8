package repository

import (
	"context"
	"fmt"
	"strconv"

	"notification-service/internal/domain"
	pbemission "notification-service/proto/emission"

	"google.golang.org/grpc/metadata"
)

type emissionGRPCClient struct {
	client pbemission.EmissionClient
}

func NewEmissionGRPCClient(client pbemission.EmissionClient) domain.EmissionClient {
	return &emissionGRPCClient{client: client}
}

func (c *emissionGRPCClient) GetDailyTotal(ctx context.Context, userID int32, date string) (float64, error) {
	resp, err := c.client.GetDailyTotal(ctx, &pbemission.DailyTotalRequest{
		UserId: userID,
		Date:   date,
	})
	if err != nil {
		return 0, fmt.Errorf("GetDailyTotal: %w", err)
	}
	return float64(resp.DailyTotalKg), nil
}

func (c *emissionGRPCClient) GetUserPreferences(ctx context.Context, userID int32) (string, float64, error) {
	md := metadata.Pairs("user-id", strconv.Itoa(int(userID)))
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := c.client.GetUserPreferences(ctx, &pbemission.Empty{})
	if err != nil {
		return "", 0, fmt.Errorf("GetUserPreferences: %w", err)
	}
	return resp.CountryCode, resp.CustomDailyLimitKgCo2, nil
}
