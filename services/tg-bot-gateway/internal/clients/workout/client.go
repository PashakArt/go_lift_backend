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

	res, err := c.client.SignInOrSignUp(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("gRPC auth call failed: %w", err)
	}

	return res, nil
}

func (c *Client) GetMuscleGroups(ctx context.Context) (*workoutv1.GetMuscleGroupsResponse, error) {
	res, err := c.client.GetMuscleGroups(ctx, &workoutv1.GetMuscleGroupsRequest{})
	if err != nil {
		return nil, fmt.Errorf("gRPC get muscle group failed: %w", err)
	}
	return res, nil
}

func (c *Client) GetExercises(ctx context.Context, userId, muscleGroupId string) (*workoutv1.GetExercisesResponse, error) {
	res, err := c.client.GetExercises(ctx, &workoutv1.GetExercisesRequest{
		UserId:        &userId,
		MuscleGroupId: muscleGroupId,
	})
	if err != nil {
		return nil, fmt.Errorf("gRPC get exercises failed: %w", err)
	}
	return res, nil
}
