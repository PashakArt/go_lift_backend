package service

import (
	"context"
	"fmt"

	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/google/uuid"
)

type ExerciseService interface {
	GetExercises(ctx context.Context, tenantIDStr string, muscleGroupCode string) ([]*domain.Exercise, error)
}

type exerciseService struct {
	exerciseRepo domain.ExerciseRepository
}

func NewExerciseService(exerciseRepo domain.ExerciseRepository) ExerciseService {
	return &exerciseService{
		exerciseRepo: exerciseRepo,
	}
}

func (s *exerciseService) GetExercises(ctx context.Context, tenantIDStr string, muscleGroupCode string) ([]*domain.Exercise, error) {
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant id format: %w", err)
	}

	exercises, err := s.exerciseRepo.List(ctx, tenantID, muscleGroupCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get exercises list from repo: %w", err)
	}

	return exercises, nil
}
