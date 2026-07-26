package service

import (
	"context"
	"fmt"

	workoutv1 "github.com/PashakArt/go_lift_backend/api/proto/workout/v1"
	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/google/uuid"
)

type TemplateService interface {
	CreateTemplate(ctx context.Context, req *workoutv1.CreateTemplateRequest) (string, error)
}

type templateService struct {
	templateRepo domain.TemplateRepository
}

func NewTemplateService(templateRepo domain.TemplateRepository) TemplateService {
	return &templateService{
		templateRepo: templateRepo,
	}
}

func (s *templateService) CreateTemplate(ctx context.Context, req *workoutv1.CreateTemplateRequest) (string, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return "", fmt.Errorf("invalid user_id: %w", err)
	}

	items := make([]domain.TemplateItem, 0, len(req.GetItems()))
	for _, itemReq := range req.GetItems() {
		exerciseID, err := uuid.Parse(itemReq.GetExerciseId())
		if err != nil {
			return "", fmt.Errorf("invalid exercise_id: %w", err)
		}

		targetSets := make([]domain.TargetSet, 0, len(itemReq.GetTargetSets()))
		for _, setReq := range itemReq.GetTargetSets() {
			targetSets = append(targetSets, domain.TargetSet{
				SetNum:          setReq.GetSetNum(),
				Weight:          setReq.Weight,
				Reps:            setReq.Reps,
				DurationSeconds: setReq.DurationSec,
				DistanceMeters:  setReq.DistanceM,
			})
		}

		items = append(items, domain.TemplateItem{
			ExerciseID: exerciseID,
			OrderIndex: itemReq.GetOrderIndex(),
			TargetSets: targetSets,
		})
	}

	template := &domain.WorkoutTemplate{
		UserID: userID,
		Name:   req.GetName(),
		Items:  items,
	}

	templateID, err := s.templateRepo.Create(ctx, template)
	if err != nil {
		return "", fmt.Errorf("failed to create template in repository: %w", err)
	}

	return templateID.String(), nil
}
