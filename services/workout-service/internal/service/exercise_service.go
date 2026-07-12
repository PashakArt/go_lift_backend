package service

import (
	"context"
	"fmt"

	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/google/uuid"
)

type ExerciseService interface {
	GetExercises(ctx context.Context, muscleGroupId, userId string) ([]*domain.Exercise, error)
}

type exerciseService struct {
	exerciseRepo domain.ExerciseRepository
}

func NewExerciseService(exerciseRepo domain.ExerciseRepository) ExerciseService {
	return &exerciseService{
		exerciseRepo: exerciseRepo,
	}
}

func (s *exerciseService) GetExercises(ctx context.Context, muscleGroupId, userId string) ([]*domain.Exercise, error) {
	parsedMuscleGroupIdId, err := uuid.Parse(muscleGroupId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse muscleGroupId: %w", err)
	}

	if userId != "" {
		parsedUserId, err := uuid.Parse(userId)
		if err != nil {
			return nil, fmt.Errorf("failed to parse user_id: %w", err)
		}

		exercises, err := s.exerciseRepo.UserFavoriteList(ctx, parsedMuscleGroupIdId, parsedUserId)
		if err != nil {
			return nil, fmt.Errorf("failed to get user favorite exercises list from repo: %w", err)
		}

		return exercises, nil
	}

	exercises, err := s.exerciseRepo.List(ctx, parsedMuscleGroupIdId)
	if err != nil {
		return nil, fmt.Errorf("failed to get exercises list from repo: %w", err)
	}

	return exercises, nil
}
