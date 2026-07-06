package handlers

import (
	"context"
	"time"

	workoutv1 "github.com/PashakArt/go_lift_backend/api/proto/workout/v1"

	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
)

type AuthService interface {
	SignInOrSignUp(ctx context.Context, tenantId, tgId string) (*domain.User, error)
}

type AuthHandler struct {
	workoutv1.UnimplementedWorkoutServiceServer
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) SignInOrSignUp(ctx context.Context, req *workoutv1.SignInOrSignUpRequest) (*workoutv1.SignInOrSignUpResponse, error) {
	user, err := h.authService.SignInOrSignUp(ctx, req.TenantId, req.TelegramId)
	if err != nil {
		return nil, err
	}

	return &workoutv1.SignInOrSignUpResponse{
		UserId:     user.UserID.String(),
		TenantId:   user.TenantID.String(),
		TelegramId: user.TelegramID,
		Phone:      user.Phone,
		CreatedAt:  user.CreatedAt.Format(time.RFC3339),
	}, nil
}
