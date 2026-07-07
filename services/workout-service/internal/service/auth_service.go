package service

import (
	"context"
	"fmt"
	"time"

	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/google/uuid"
)

type AuthService interface {
	SignInOrSignUp(ctx context.Context, tenantId, tgId string) (*domain.User, error)
}

type authService struct {
	r domain.UserRepository
}

func NewAuthService(r domain.UserRepository) AuthService {
	return &authService{r: r}
}

func (s *authService) SignInOrSignUp(ctx context.Context, tenantId, tgId string) (*domain.User, error) {
	parsedTenantId, err := uuid.Parse(tenantId)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant id format in service: %w", err)
	}

	existingUser, err := s.r.GetByTenantAndTelegramID(ctx, parsedTenantId, tgId)
	if err != nil {
		return nil, fmt.Errorf("auth check failed: %w", err)
	}

	if existingUser != nil {
		return existingUser, nil
	}

	newUser := &domain.User{
		UserID:     uuid.New(),
		TelegramID: tgId,
		TenantID:   &parsedTenantId,
		CreatedAt:  time.Now(),
	}

	err = s.r.Create(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to register new user: %w", err)
	}

	return newUser, nil
}
