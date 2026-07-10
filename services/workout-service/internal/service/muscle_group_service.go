package service

import (
	"context"
	"fmt"

	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
)

type MuscleGroupService interface {
	List(ctx context.Context) ([]*domain.MuscleGroup, error)
}

type muscleGroupService struct {
	r domain.MuscleGroupRepository
}

func NewMuscleGroupService(r domain.MuscleGroupRepository) MuscleGroupService {
	return &muscleGroupService{
		r: r,
	}
}

func (s *muscleGroupService) List(ctx context.Context) ([]*domain.MuscleGroup, error) {
	muscleGroups, err := s.r.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get muscle group list from repo: %w", err)
	}

	return muscleGroups, nil
}
