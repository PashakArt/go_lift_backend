package workout

import (
	"context"
	"fmt"

	workoutv1 "github.com/PashakArt/go_lift_backend/api/proto/workout/v1"
	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/types"
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

func (c *Client) Init(ctx context.Context, tenantId, tgId string) (*workoutv1.InitResponse, error) {
	req := &workoutv1.InitRequest{
		TenantId:   tenantId,
		TelegramId: tgId,
	}

	res, err := c.client.Init(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("gRPC init call failed: %w", err)
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

func (c *Client) StartTraining(ctx context.Context, tenantId, userId string) (*workoutv1.StartWorkoutSessionResponse, error) {
	res, err := c.client.StartWorkoutSession(ctx, &workoutv1.StartWorkoutSessionRequest{
		UserId:   userId,
		TenantId: tenantId,
		// TODO расширить до кроссфит
		Type: 1,
	})

	if err != nil {
		return nil, fmt.Errorf("gRPC start session failed: %w", err)
	}

	return res, nil
}

func (c *Client) FinishTraining(ctx context.Context, userId string) error {
	_, err := c.client.FinishWorkoutSession(ctx, &workoutv1.FinishWorkoutSessionRequest{
		UserId: userId,
	})

	if err != nil {
		return fmt.Errorf("gRPC start session failed: %w", err)
	}

	return nil
}

func (c *Client) LogWorkoutSet(ctx context.Context, req types.LogSetRequest, tenantId string) (*workoutv1.LogSetResponse, error) {
	fmt.Println("222")
	res, err := c.client.LogWorkoutSet(ctx, &workoutv1.LogSetRequest{
		SessionId:  req.SessionId,
		ExerciseId: req.ExerciseId,
		TenantId:   tenantId,
		SetId:      req.SetID,
		SetNumber:  req.SetNumber,

		Weight:      req.Weight,
		Reps:        req.Reps,
		DurationSec: req.DurationSeconds,
		DistanceMet: req.DistanceMeters,
	})
	fmt.Println("333")

	if err != nil {
		return nil, fmt.Errorf("gRPC log workout set failed: %w", err)
	}
	fmt.Println("444")

	return res, nil
}
