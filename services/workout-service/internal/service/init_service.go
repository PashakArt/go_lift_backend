package service

import (
	"context"
	"fmt"

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
		return nil, fmt.Errorf("failed to check tenantId: %w", err)
	}

	user := &domain.User{
		TenantID:    tenant.TenantID,
		TelegramID:  tgId,
		TgUsername:  username,
		TgFirstName: firstName,
		TgLastName:  lastName,
	}

	isNewUser, err := s.userRepo.Upsert(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to sync user profile: %w", err)
	}

	var activeSession *domain.WorkoutSession
	if !isNewUser {
		activeSession, err = s.sessionRepo.GetActiveByUserID(ctx, user.UserID.String())
		if err != nil {
			return nil, fmt.Errorf("failed to get active session by user_id: %w", err)
		}
	}

	return &InitResponse{
		User:           user,
		IsNewUser:      isNewUser,
		ActiveSession:  activeSession,
		TenantBranding: tenant.BrandingJSON,
	}, nil
}
