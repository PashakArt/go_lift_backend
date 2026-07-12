package handlers

import (
	"context"
	"time"

	workoutv1 "github.com/PashakArt/go_lift_backend/api/proto/workout/v1"
	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WorkoutHandler struct {
	workoutv1.UnimplementedWorkoutServiceServer
	authService        service.AuthService
	exerciseService    service.ExerciseService
	sessionService     service.WorkoutSessionService
	muscleGroupService service.MuscleGroupService
}

func NewWorkoutHandler(
	authService service.AuthService,
	exerciseService service.ExerciseService,
	sessionService service.WorkoutSessionService,
	muscleGroupService service.MuscleGroupService,
) *WorkoutHandler {
	return &WorkoutHandler{
		authService:        authService,
		exerciseService:    exerciseService,
		sessionService:     sessionService,
		muscleGroupService: muscleGroupService,
	}
}

func (h *WorkoutHandler) SignInOrSignUp(ctx context.Context, req *workoutv1.SignInOrSignUpRequest) (*workoutv1.SignInOrSignUpResponse, error) {
	response, err := h.authService.SignInOrSignUp(ctx, req.TenantId, req.TelegramId)
	if err != nil {
		return nil, err
	}

	var activeSessionInfo *workoutv1.ActiveSessionInfo
	if response.ActiveSession != nil {
		activeSessionInfo = &workoutv1.ActiveSessionInfo{
			SessionId: response.ActiveSession.SessionID.String(),
			StartedAt: response.ActiveSession.StartedAt.Format(time.RFC3339),
		}
	}

	return &workoutv1.SignInOrSignUpResponse{
		User: &workoutv1.UserInfo{
			UserId:     response.User.UserID.String(),
			TenantId:   response.User.TenantID.String(),
			TelegramId: response.User.TelegramID,
			Phone:      response.User.Phone,
			CreatedAt:  response.User.CreatedAt.Format(time.RFC3339),
		},
		ActiveSession: activeSessionInfo,
		IsNewUser:     response.IsNewUser,
	}, nil
}

func (h *WorkoutHandler) GetExercises(ctx context.Context, req *workoutv1.GetExercisesRequest) (*workoutv1.GetExercisesResponse, error) {
	domainExercises, err := h.exerciseService.GetExercises(ctx, req.GetMuscleGroupId(), req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch exercises: %v", err)
	}

	protoExercises := make([]*workoutv1.ExerciseInfo, 0, len(domainExercises))
	for _, ex := range domainExercises {
		protoExercises = append(protoExercises, &workoutv1.ExerciseInfo{
			ExerciseId:       ex.ExerciseID.String(),
			Name:             ex.Name,
			Type:             mapDomainTypeToProto(ex.Type),
			IsGlobal:         ex.IsGlobal,
			MuscleGroupCodes: ex.MuscleGroupCodes,
			CreatedAt:        timestamppb.New(ex.CreatedAt),
		})
	}

	return &workoutv1.GetExercisesResponse{
		Exercises: protoExercises,
	}, nil
}

func (h *WorkoutHandler) GetMuscleGroups(ctx context.Context, req *workoutv1.GetMuscleGroupsRequest) (*workoutv1.GetMuscleGroupsResponse, error) {
	domainMuscleGroups, err := h.muscleGroupService.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch exercises: %v", err)
	}

	protoMuscleGroups := make([]*workoutv1.MuscleGroup, 0, len(domainMuscleGroups))
	for _, mg := range domainMuscleGroups {
		protoMuscleGroups = append(protoMuscleGroups, &workoutv1.MuscleGroup{
			MuscleGroupId: mg.MuscleGroupId,
			Code:          mg.Code,
			Name:          mg.Name,
		})
	}

	return &workoutv1.GetMuscleGroupsResponse{
		MuscleGroups: protoMuscleGroups,
	}, nil
}

func (h *WorkoutHandler) StartWorkoutSession(
	ctx context.Context,
	req *workoutv1.StartWorkoutSessionRequest,
) (*workoutv1.StartWorkoutSessionResponse, error) {
	tenantId := req.GetTenantId()
	userId := req.GetUserId()

	if tenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant_id is required")
	}

	if userId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	domainType := mapProtoTypeToDomain(req.GetType())

	session, err := h.sessionService.StartSession(
		ctx,
		tenantId,
		userId,
		string(domainType),
		req.GetTemplateId(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start workout session: %v", err)
	}

	return &workoutv1.StartWorkoutSessionResponse{
		SessionId: session.SessionID.String(),
		TenantId:  tenantId,
		UserId:    userId,
		Type:      req.GetType(),
		StartedAt: timestamppb.New(session.StartedAt),
	}, nil
}

func mapDomainTypeToProto(t domain.ExerciseType) workoutv1.ExerciseType {
	switch t {
	case domain.ExerciseTypeDynamic:
		return workoutv1.ExerciseType_EXERCISE_TYPE_DYNAMIC
	case domain.ExerciseTypeStatic:
		return workoutv1.ExerciseType_EXERCISE_TYPE_STATIC
	case domain.ExerciseTypeBodyweight:
		return workoutv1.ExerciseType_EXERCISE_TYPE_BODYWEIGHT
	case domain.ExerciseTypeCardio:
		return workoutv1.ExerciseType_EXERCISE_TYPE_CARDIO
	default:
		return workoutv1.ExerciseType_EXERCISE_TYPE_UNSPECIFIED
	}
}

func mapProtoTypeToDomain(t workoutv1.SessionType) domain.SessionType {
	switch t {
	case workoutv1.SessionType_SESSION_TYPE_CLASSIC:
		return domain.SessionTypeClassic
	case workoutv1.SessionType_SESSION_TYPE_CIRCUIT:
		return domain.SessionTypeCircuit
	default:
		return domain.SessionTypeClassic
	}
}
