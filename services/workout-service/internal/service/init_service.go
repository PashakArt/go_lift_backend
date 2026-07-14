package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/google/uuid"
)

type SignInOrSignUpResponse struct {
	User          *domain.User
	IsNewUser     bool
	ActiveSession *domain.WorkoutSession
}

type InitService interface {
	Init(ctx context.Context, tenantId, tgId string) (*SignInOrSignUpResponse, error)
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

func (s *initService) Init(ctx context.Context, tenantId, tgId string) (*SignInOrSignUpResponse, error) {
	parsedTenantId, err := uuid.Parse(tenantId)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant id format in service: %w", err)
	}

	existingUser, err := s.userRepo.GetByTenantAndTelegramID(ctx, parsedTenantId, tgId)
	if err != nil {
		return nil, fmt.Errorf("auth check failed: %w", err)
	}

	if existingUser == nil {
		log.Printf("1111")
		tenant, err := s.tenantRepo.GetById(ctx, parsedTenantId)
		log.Printf("2222")
		if err != nil {
			return nil, fmt.Errorf("failed to checking tenantId: %w", err)
		}

		targetTenantID := parsedTenantId
		if tenant == nil {
			targetTenantID = uuid.Nil
		}

		newUser := &domain.User{
			UserID:     uuid.New(),
			TelegramID: tgId,
			TenantID:   targetTenantID,
			CreatedAt:  time.Now(),
		}

		err = s.userRepo.Create(ctx, newUser)
		if err != nil {
			return nil, fmt.Errorf("failed to register new user: %w", err)
		}

		return &SignInOrSignUpResponse{
			User:          newUser,
			IsNewUser:     true,
			ActiveSession: nil,
		}, nil
	}

	activeSession, err := s.sessionRepo.GetActiveByUserID(ctx, existingUser.UserID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get active session by user_id: %w", err)
	}

	return &SignInOrSignUpResponse{
		User:          existingUser,
		IsNewUser:     false,
		ActiveSession: activeSession,
	}, nil
}
