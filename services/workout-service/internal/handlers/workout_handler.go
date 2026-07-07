package handlers

import (
	"context"
	"time"

	desc "github.com/PashakArt/go_lift_backend/api/proto/workout/v1"
	workoutv1 "github.com/PashakArt/go_lift_backend/api/proto/workout/v1"
	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WorkoutHandler struct {
	workoutv1.UnimplementedWorkoutServiceServer
	authService     service.AuthService
	exerciseService service.ExerciseService // Вторая зависимость
}

func NewWorkoutHandler(authService service.AuthService, exerciseService service.ExerciseService) *WorkoutHandler {
	return &WorkoutHandler{
		authService:     authService,
		exerciseService: exerciseService,
	}
}

func (h *WorkoutHandler) SignInOrSignUp(ctx context.Context, req *workoutv1.SignInOrSignUpRequest) (*workoutv1.SignInOrSignUpResponse, error) {
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

func (h *WorkoutHandler) GetExercises(ctx context.Context, req *desc.GetExercisesRequest) (*desc.GetExercisesResponse, error) {
	if req.GetTenantId() == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant_id is required")
	}

	domainExercises, err := h.exerciseService.GetExercises(ctx, req.GetTenantId(), req.GetMuscleGroupCode())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch exercises: %v", err)
	}

	protoExercises := make([]*desc.ExerciseInfo, 0, len(domainExercises))
	for _, ex := range domainExercises {
		protoExercises = append(protoExercises, &desc.ExerciseInfo{
			ExerciseId:       ex.ExerciseID.String(),
			Name:             ex.Name,
			Type:             mapDomainTypeToProto(ex.Type),
			IsGlobal:         ex.IsGlobal,
			MuscleGroupCodes: ex.MuscleGroupCodes,
			CreatedAt:        timestamppb.New(ex.CreatedAt),
		})
	}

	return &desc.GetExercisesResponse{
		Exercises: protoExercises,
	}, nil
}

func mapDomainTypeToProto(t domain.ExerciseType) desc.ExerciseType {
	switch t {
	case domain.ExerciseTypeDynamic:
		return desc.ExerciseType_EXERCISE_TYPE_DYNAMIC
	case domain.ExerciseTypeStatic:
		return desc.ExerciseType_EXERCISE_TYPE_STATIC
	case domain.ExerciseTypeBodyweight:
		return desc.ExerciseType_EXERCISE_TYPE_BODYWEIGHT
	case domain.ExerciseTypeCardio:
		return desc.ExerciseType_EXERCISE_TYPE_CARDIO
	default:
		return desc.ExerciseType_EXERCISE_TYPE_UNSPECIFIED
	}
}
