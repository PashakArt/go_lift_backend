package workout

import (
	"context"
	"fmt"

	workoutv1 "github.com/PashakArt/go_lift_backend/api/proto/workout/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client workoutv1.WorkoutServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to workout service: %w", err)
	}

	return &Client{
		client: workoutv1.NewWorkoutServiceClient(conn),
	}, nil
}

func (c *Client) Auth(ctx context.Context, tenantId, tgId string) (*workoutv1.SignInOrSignUpResponse, error) {
	req := &workoutv1.SignInOrSignUpRequest{
		TenantId:   tenantId,
		TelegramId: tgId,
	}

	res, err := c.client.SignInOrSignup(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("gRPC auth call failed: %w", err)
	}

	return res, nil
}
