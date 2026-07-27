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
	GetTemplates(ctx context.Context, userID string) ([]*domain.WorkoutTemplate, error)
	GetTemplate(ctx context.Context, templateID, userID string) (*domain.WorkoutTemplate, error)
}

type templateService struct {
	templateRepo domain.TemplateRepository
	exerciseRepo domain.ExerciseRepository
}

func NewTemplateService(
	templateRepo domain.TemplateRepository,
	exerciseRepo domain.ExerciseRepository,
) TemplateService {
	return &templateService{
		templateRepo: templateRepo,
		exerciseRepo: exerciseRepo,
	}
}

func (s *templateService) CreateTemplate(ctx context.Context, req *workoutv1.CreateTemplateRequest) (string, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return "", fmt.Errorf("invalid user_id: %w", err)
	}

	exIDMap := make(map[uuid.UUID]struct{})
	var uniqueExIDs []uuid.UUID

	items := make([]domain.TemplateItem, 0, len(req.GetItems()))
	for _, itemReq := range req.GetItems() {
		exerciseID, err := uuid.Parse(itemReq.GetExerciseId())
		if err != nil {
			return "", fmt.Errorf("invalid exercise_id %s: %w", itemReq.GetExerciseId(), err)
		}

		if _, exists := exIDMap[exerciseID]; !exists {
			exIDMap[exerciseID] = struct{}{}
			uniqueExIDs = append(uniqueExIDs, exerciseID)
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

	existingExercises, err := s.exerciseRepo.GetByIDs(ctx, uniqueExIDs)
	if err != nil {
		return "", fmt.Errorf("failed to validate exercise ids: %w", err)
	}

	if len(existingExercises) != len(uniqueExIDs) {
		return "", fmt.Errorf("validation failed: one or more exercise_ids do not exist")
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

func (s *templateService) GetTemplates(ctx context.Context, rawUserID string) ([]*domain.WorkoutTemplate, error) {
	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	templates, err := s.templateRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get templates from repository: %w", err)
	}

	return templates, nil
}

func (s *templateService) GetTemplate(ctx context.Context, rawTemplateID, rawUserID string) (*domain.WorkoutTemplate, error) {
	templateID, err := uuid.Parse(rawTemplateID)
	if err != nil {
		return nil, fmt.Errorf("invalid template_id: %w", err)
	}

	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	template, err := s.templateRepo.GetByID(ctx, templateID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get template by id from repository: %w", err)
	}

	if len(template.Items) == 0 {
		return template, nil
	}

	exIDMap := make(map[uuid.UUID]struct{})
	var uniqueExIDs []uuid.UUID
	for _, item := range template.Items {
		if _, exists := exIDMap[item.ExerciseID]; !exists {
			exIDMap[item.ExerciseID] = struct{}{}
			uniqueExIDs = append(uniqueExIDs, item.ExerciseID)
		}
	}

	exercises, err := s.exerciseRepo.GetByIDs(ctx, uniqueExIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exercises for template enrichment: %w", err)
	}

	exerciseLookup := make(map[uuid.UUID]*domain.Exercise, len(exercises))
	for _, ex := range exercises {
		exerciseLookup[ex.ExerciseID] = ex
	}

	for i := range template.Items {
		if ex, ok := exerciseLookup[template.Items[i].ExerciseID]; ok {
			template.Items[i].ExerciseName = ex.Name
			template.Items[i].ExerciseType = ex.Type
		}
	}

	return template, nil
}
