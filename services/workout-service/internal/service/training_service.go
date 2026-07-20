package service

import (
	"context"
	"fmt"

	workoutv1 "github.com/PashakArt/go_lift_backend/api/proto/workout/v1"
	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/google/uuid"
)

type TrainingService interface {
	StartSession(ctx context.Context, tenantIDStr, userIDStr, sessionType, templateIDStr string) (*domain.WorkoutSession, error)
	LogWorkoutSet(ctx context.Context, req *workoutv1.LogSetRequest) (*domain.WorkoutSet, error)
}

type trainingService struct {
	sessionRepo    domain.WorkoutSessionRepository
	workoutSetRepo domain.WorkoutSetRepository
}

func NewTrainingService(sessionRepo domain.WorkoutSessionRepository) TrainingService {
	return &trainingService{
		sessionRepo: sessionRepo,
	}
}

func (s *trainingService) StartSession(ctx context.Context, tenantIDStr, userIDStr, sessionType, templateIDStr string) (*domain.WorkoutSession, error) {
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant id format: %w", err)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id format: %w", err)
	}

	var templateID *uuid.UUID
	if templateIDStr != "" {
		parsedTemplateID, err := uuid.Parse(templateIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid template id format: %w", err)
		}
		templateID = &parsedTemplateID
	}

	session := &domain.WorkoutSession{
		TenantID:   tenantID,
		UserID:     userID,
		TemplateID: templateID,
		Type:       domain.SessionType(sessionType),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create workout session in repo: %w", err)
	}

	return session, nil
}

func (s *trainingService) LogWorkoutSet(
	ctx context.Context,
	protoReq *workoutv1.LogSetRequest,
) (*domain.WorkoutSet, error) {
	sessionID, err := uuid.Parse(protoReq.GetSessionId())
	if err != nil {
		return nil, fmt.Errorf("invalid session id: %w", err)
	}

	exerciseID, err := uuid.Parse(protoReq.GetExerciseId())
	if err != nil {
		return nil, fmt.Errorf("invalid exercise id: %w", err)
	}

	var setID uuid.UUID
	if protoReq.SetId != nil && *protoReq.SetId != "" {
		parsedSetID, err := uuid.Parse(*protoReq.SetId)
		if err != nil {
			return nil, fmt.Errorf("invalid set id format: %w", err)
		}
		setID = parsedSetID
	}

	var weight *float64
	if protoReq.Weight != nil {
		w := float64(*protoReq.Weight)
		weight = &w
	}

	var reps *int
	if protoReq.Reps != nil {
		r := int(*protoReq.Reps)
		reps = &r
	}

	var durationSec *int
	if protoReq.DurationSec != nil { // Если в proto было duration_sec
		d := int(*protoReq.DurationSec)
		durationSec = &d
	}

	var distanceMet *int
	if protoReq.DistanceMet != nil { // Если в proto было distance_met
		dm := int(*protoReq.DistanceMet)
		distanceMet = &dm
	}

	set := &domain.WorkoutSet{
		SetID:           setID,
		SessionID:       sessionID,
		ExerciseID:      exerciseID,
		SetNumber:       int(protoReq.SetNumber),
		Weight:          weight,
		Reps:            reps,
		DurationSeconds: durationSec,
		DistanceMeters:  distanceMet,
	}

	res, err := s.workoutSetRepo.LogWorkoutSet(ctx, set)
	if err != nil {
		return nil, fmt.Errorf("failed to log workout set: %w", err)
	}

	return res, nil
}
