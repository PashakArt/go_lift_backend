package service

import (
	"context"
	"fmt"
	"time"

	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/google/uuid"
)

type InitResponse struct {
	User           *domain.User
	IsNewUser      bool
	ActiveSession  *domain.WorkoutSession
	TenantBranding []byte
}

type InitService interface {
	Init(ctx context.Context, tenantId, tgId, username, firstName, lastName string) (*InitResponse, error)
}

type initService struct {
	userRepo    domain.UserRepository
	sessionRepo domain.WorkoutSessionRepository
	tenantRepo  domain.TenantRepository
}

func NewInitService(
	ur domain.UserRepository,
	sr domain.WorkoutSessionRepository,
	tr domain.TenantRepository,
) InitService {
	return &initService{
		userRepo:    ur,
		sessionRepo: sr,
		tenantRepo:  tr,
	}
}

func (s *initService) Init(ctx context.Context, tenantId, tgId, username, firstName, lastName string) (*InitResponse, error) {
	parsedTenantId, err := uuid.Parse(tenantId)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant id format in service: %w", err)
	}

	tenant, err := s.tenantRepo.GetById(ctx, parsedTenantId)
	if err != nil {
		return nil, fmt.Errorf("failed to checking tenantId: %w", err)
	}

	existingUser, err := s.userRepo.GetByTenantAndTelegramID(ctx, parsedTenantId, tgId)
	if err != nil {
		return nil, fmt.Errorf("auth check failed: %w", err)
	}

	if existingUser == nil {
		newUser := &domain.User{
			UserID:      uuid.New(),
			TelegramID:  tgId,
			TenantID:    tenant.TenantID,
			CreatedAt:   time.Now(),
			TgUsername:  username,
			TgFirstName: firstName,
			TgLastName:  lastName,
		}

		err = s.userRepo.Create(ctx, newUser)
		if err != nil {
			return nil, fmt.Errorf("failed to register new user: %w", err)
		}

		return &InitResponse{
			User:           newUser,
			IsNewUser:      true,
			ActiveSession:  nil,
			TenantBranding: tenant.BrandingJSON,
		}, nil
	}

	activeSession, err := s.sessionRepo.GetActiveByUserID(ctx, existingUser.UserID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get active session by user_id: %w", err)
	}

	return &InitResponse{
		User:           existingUser,
		IsNewUser:      false,
		ActiveSession:  activeSession,
		TenantBranding: tenant.BrandingJSON,
	}, nil
}
