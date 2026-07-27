package handlers

import (
	"context"
	"time"

	workoutv1 "github.com/PashakArt/go_lift_backend/api/proto/workout/v1"
	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WorkoutHandler struct {
	workoutv1.UnimplementedWorkoutServiceServer
	authService        service.InitService
	exerciseService    service.ExerciseService
	trainingService    service.TrainingService
	muscleGroupService service.MuscleGroupService
	templateService    service.TemplateService
}

func NewWorkoutHandler(
	authService service.InitService,
	exerciseService service.ExerciseService,
	sessionService service.TrainingService,
	muscleGroupService service.MuscleGroupService,
	templateService service.TemplateService,
) *WorkoutHandler {
	return &WorkoutHandler{
		authService:        authService,
		exerciseService:    exerciseService,
		trainingService:    sessionService,
		muscleGroupService: muscleGroupService,
		templateService:    templateService,
	}
}

func (h *WorkoutHandler) Init(ctx context.Context, req *workoutv1.InitRequest) (*workoutv1.InitResponse, error) {
	response, err := h.authService.Init(ctx, req.TenantId, req.TelegramId)
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

	brandingProto := &workoutv1.TenantBranding{}
	if len(response.TenantBranding) > 0 {
		err = protojson.Unmarshal(response.TenantBranding, brandingProto)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to unmarshal tenant branding json: %v", err)
		}
	}

	return &workoutv1.InitResponse{
		User: &workoutv1.UserInfo{
			UserId:     response.User.UserID.String(),
			TenantId:   response.User.TenantID.String(),
			TelegramId: response.User.TelegramID,
			Phone:      response.User.Phone,
			CreatedAt:  response.User.CreatedAt.Format(time.RFC3339),
		},
		ActiveSession:  activeSessionInfo,
		IsNewUser:      response.IsNewUser,
		TenantBranding: brandingProto}, nil
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
			Type:             MapDomainTypeToProto(ex.Type),
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

	domainType := MapProtoTypeToDomain(req.GetType())

	session, err := h.trainingService.StartSession(
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

func (h *WorkoutHandler) FinishWorkoutSession(
	ctx context.Context,
	req *workoutv1.FinishWorkoutSessionRequest,
) (*workoutv1.FinishWorkoutSessionResponse, error) {
	err := h.trainingService.FinishSession(ctx, req.GetUserId())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to finish workout session: %v", err)
	}

	return nil, nil
}

func (h *WorkoutHandler) LogWorkoutSet(
	ctx context.Context,
	protoReq *workoutv1.LogSetRequest,
) (*workoutv1.LogSetResponse, error) {
	workoutSet, err := h.trainingService.LogWorkoutSet(ctx, protoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to log workout set: %v", err)
	}

	return &workoutv1.LogSetResponse{
		SetId:     workoutSet.SetID.String(),
		SetNumber: int32(workoutSet.SetNumber),
	}, nil
}

func (h *WorkoutHandler) GetCompletedExercises(
	ctx context.Context,
	req *workoutv1.GetCompletedExercisesRequest,
) (*workoutv1.GetCompletedExercisesResponse, error) {
	completedExercises, err := h.trainingService.GetCompletedExercises(ctx, req.GetUserId(), req.GetExerciseId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get completed exercises: %v", err)
	}

	sets := make([]*workoutv1.CompletedSet, 0, len(completedExercises))
	for _, set := range completedExercises {
		item := &workoutv1.CompletedSet{
			SetNumber: int32(set.SetNumber),
			SetId:     set.SetID.String(),
		}

		if set.Weight != nil {
			w := float32(*set.Weight)
			item.Weight = &w
		}
		if set.Reps != nil {
			r := int32(*set.Reps)
			item.Reps = &r
		}
		if set.DurationSeconds != nil {
			d := int32(*set.DurationSeconds)
			item.DurationSec = &d
		}
		if set.DistanceMeters != nil {
			dist := int32(*set.DistanceMeters)
			item.DistanceM = &dist
		}

		sets = append(sets, item)
	}

	return &workoutv1.GetCompletedExercisesResponse{
		Sets: sets,
	}, nil
}

func (h *WorkoutHandler) GetTrainingDays(ctx context.Context, req *workoutv1.GetTrainingDaysRequest) (*workoutv1.GetTrainingDaysResponse, error) {
	days, err := h.trainingService.GetTrainingDays(ctx, req.GetUserId(), int(req.GetYear()), int(req.GetMonths()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get training days: %v", err)
	}

	return &workoutv1.GetTrainingDaysResponse{
		TrainingDays: days,
	}, nil
}

func (h *WorkoutHandler) GetWorkoutsForDay(ctx context.Context, req *workoutv1.GetWorkoutsForDayRequest) (*workoutv1.GetWorkoutsForDayResponse, error) {
	workoutForDay, err := h.trainingService.GetWorkoutsForDay(ctx, req.GetUserId(), req.GetDate())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get workouts for day: %v", err)
	}

	pbSessions := make([]*workoutv1.WorkoutSessionDTO, 0, len(workoutForDay.Sessions))
	for _, session := range workoutForDay.Sessions {

		pbExercises := make([]*workoutv1.ExerciseDTO, 0, len(session.Exercises))
		for _, ex := range session.Exercises {

			pbSets := make([]*workoutv1.CompletedSet, 0, len(ex.Sets))
			for _, st := range ex.Sets {
				pbSet := &workoutv1.CompletedSet{
					SetId:     st.SetId,
					SetNumber: int32(st.SetNumber),
				}
				if st.Weight != nil {
					w := float32(*st.Weight)
					pbSet.Weight = &w
				}
				if st.Reps != nil {
					r := *st.Reps
					pbSet.Reps = &r
				}
				if st.DurationSec != nil {
					d := *st.DurationSec
					pbSet.DurationSec = &d
				}
				if st.DistanceM != nil {
					d := *st.DistanceM
					pbSet.DistanceM = &d
				}
				pbSets = append(pbSets, pbSet)
			}

			pbExercises = append(pbExercises, &workoutv1.ExerciseDTO{
				ExerciseId: ex.ExerciseID,
				Name:       ex.Name,
				Type:       ex.Type,
				Sets:       pbSets,
			})
		}

		var endedAtStr string
		if session.EndedAt != nil {
			endedAtStr = session.EndedAt.Format(time.RFC3339)
		}

		pbSessions = append(pbSessions, &workoutv1.WorkoutSessionDTO{
			SessionId:       session.SessionID,
			StartedAt:       session.StartedAt.Format(time.RFC3339),
			EndedAt:         endedAtStr,
			DurationSeconds: session.DurationSeconds,
			Exercises:       pbExercises,
		})
	}

	return &workoutv1.GetWorkoutsForDayResponse{
		Date:     workoutForDay.Date,
		Sessions: pbSessions,
	}, nil
}

func (h *WorkoutHandler) CreateTemplate(
	ctx context.Context,
	req *workoutv1.CreateTemplateRequest,
) (*workoutv1.CreateTemplateResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if len(req.GetItems()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "template items cannot be empty")
	}

	templateID, err := h.templateService.CreateTemplate(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create template: %v", err)
	}

	return &workoutv1.CreateTemplateResponse{
		TemplateId: templateID,
	}, nil
}
func (h *WorkoutHandler) GetTemplates(
	ctx context.Context,
	req *workoutv1.GetTemplatesRequest,
) (*workoutv1.GetTemplatesResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	domainTemplates, err := h.templateService.GetTemplates(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch templates: %v", err)
	}

	protoTemplates := make([]*workoutv1.TemplateSummary, 0, len(domainTemplates))
	for _, t := range domainTemplates {
		protoTemplates = append(protoTemplates, &workoutv1.TemplateSummary{
			TemplateId:     t.TemplateID.String(),
			Name:           t.Name,
			ExercisesCount: int32(len(t.Items)),
			CreatedAt:      timestamppb.New(t.CreatedAt),
		})
	}

	return &workoutv1.GetTemplatesResponse{
		Templates: protoTemplates,
	}, nil
}

func (h *WorkoutHandler) GetTemplate(
	ctx context.Context,
	req *workoutv1.GetTemplateRequest,
) (*workoutv1.GetTemplateResponse, error) {
	if req.GetTemplateId() == "" {
		return nil, status.Error(codes.InvalidArgument, "template_id is required")
	}

	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	tpl, err := h.templateService.GetTemplate(ctx, req.GetTemplateId(), req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch template: %v", err)
	}

	protoItems := make([]*workoutv1.TemplateItem, 0, len(tpl.Items))
	for _, item := range tpl.Items {
		targetSets := make([]*workoutv1.TargetSet, 0, len(item.TargetSets))
		for _, set := range item.TargetSets {
			targetSets = append(targetSets, &workoutv1.TargetSet{
				SetNum:      set.SetNum,
				Weight:      set.Weight,
				Reps:        set.Reps,
				DurationSec: set.DurationSeconds,
				DistanceM:   set.DistanceMeters,
			})
		}

		protoItems = append(protoItems, &workoutv1.TemplateItem{
			ExerciseId: item.ExerciseID.String(),
			OrderIndex: item.OrderIndex,
			TargetSets: targetSets,
			Name:       item.ExerciseName,
			Type:       MapDomainTypeToProto(item.ExerciseType),
		})
	}

	return &workoutv1.GetTemplateResponse{
		TemplateId: tpl.TemplateID.String(),
		Name:       tpl.Name,
		Items:      protoItems,
		CreatedAt:  timestamppb.New(tpl.CreatedAt),
	}, nil
}

func (h *WorkoutHandler) DeleteTemplate(
	ctx context.Context,
	req *workoutv1.DeleteTemplateRequest,
) (*workoutv1.DeleteTemplateResponse, error) {
	if req.GetTemplateId() == "" {
		return nil, status.Error(codes.InvalidArgument, "template_id is required")
	}

	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	err := h.templateService.DeleteTemplate(ctx, req.GetTemplateId(), req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete template: %v", err)
	}

	return &workoutv1.DeleteTemplateResponse{}, nil
}

func (h *WorkoutHandler) UpdateTemplate(
	ctx context.Context,
	req *workoutv1.UpdateTemplateRequest,
) (*workoutv1.UpdateTemplateResponse, error) {
	if req.GetTemplateId() == "" {
		return nil, status.Error(codes.InvalidArgument, "template_id is required")
	}

	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if len(req.GetItems()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "template items cannot be empty")
	}

	err := h.templateService.UpdateTemplate(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update template: %v", err)
	}

	return &workoutv1.UpdateTemplateResponse{}, nil
}
