package service

import (
	"context"
	"fmt"

	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/google/uuid"
)

type WorkoutSessionService interface {
	StartSession(ctx context.Context, tenantIDStr, userIDStr, sessionType string, templateIDStr string) (*domain.WorkoutSession, error)
}

type workoutSessionService struct {
	sessionRepo domain.WorkoutSessionRepository
}

func NewWorkoutSessionService(sessionRepo domain.WorkoutSessionRepository) WorkoutSessionService {
	return &workoutSessionService{
		sessionRepo: sessionRepo,
	}
}

func (s *workoutSessionService) StartSession(ctx context.Context, tenantIDStr, userIDStr, sessionType string, templateIDStr string) (*domain.WorkoutSession, error) {
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
