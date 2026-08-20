package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/auth"
	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/clients/workout"
	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/types"
	"github.com/go-playground/validator/v10"
)

type WorkoutHandler struct {
	workoutClient *workout.Client
	validate      *validator.Validate
}

func NewWorkoutHandler(workoutClient *workout.Client, validate *validator.Validate) *WorkoutHandler {
	return &WorkoutHandler{
		workoutClient: workoutClient,
		validate:      validate,
	}
}

func (h *WorkoutHandler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler) {
	mux.Handle("POST /api/v1/start", authMiddleware(http.HandlerFunc(h.HandleStartTraining)))
	mux.Handle("POST /api/v1/finish", authMiddleware(http.HandlerFunc(h.HandleFinishTraining)))
	mux.Handle("POST /api/v1/workout/sets", authMiddleware(http.HandlerFunc(h.HandleLogWorkoutSet)))

	mux.Handle("GET /api/v1/muscle-groups", authMiddleware(http.HandlerFunc(h.HandleGetMuscleGroups)))
	mux.Handle("GET /api/v1/{muscleGroupId}/exercises", authMiddleware(http.HandlerFunc(h.HandleGetExercises)))
	mux.Handle("GET /api/v1/exercises/{exerciseId}/completed", authMiddleware(http.HandlerFunc(h.HandleGetCompletedExercises)))

	mux.Handle("GET /api/v1/workouts/calendar", authMiddleware(http.HandlerFunc(h.HandleGetTrainingDays)))
	mux.Handle("GET /api/v1/workouts/day", authMiddleware(http.HandlerFunc(h.HandleGetWorkoutsForDay)))
	mux.Handle("GET /api/v1/sessions/{sessionId}/exercises", authMiddleware(http.HandlerFunc(h.HandleGetSessionExercises)))

	// mux.Handle("POST /api/v1/templates", authMiddleware(http.HandlerFunc(h.HandleCreateTemplate)))
	// mux.Handle("GET /api/v1/templates", authMiddleware(http.HandlerFunc(h.HandleGetTemplates)))
	// mux.Handle("GET /api/v1/templates/detail/{templateId}", authMiddleware(http.HandlerFunc(h.HandleGetTemplate)))
	// mux.Handle("DELETE /api/v1/templates/{templateId}", authMiddleware(http.HandlerFunc(h.HandleDeleteTemplate)))
	// mux.Handle("PUT /api/v1/templates/{templateId}", authMiddleware(http.HandlerFunc(h.HandleUpdateTemplate)))
}

func (h *WorkoutHandler) HandleStartTraining(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)
	if userId == "" {
		RespondWithError(w, http.StatusBadRequest, "JWT has not user_id")
		return
	}

	tenantId := auth.TenantIDFromContext(r.Context())
	if tenantId == "" {
		RespondWithError(w, http.StatusBadRequest, "JWT has not tenant_id")
		return
	}

	var req types.StartTrainingRequest
	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		defer r.Body.Close()
	}

	res, err := h.workoutClient.StartTraining(ctx, tenantId, req.TemplateID, userId)
	if err != nil {
		log.Printf("gRPC StartTraining failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	RespondWithJSON(w, http.StatusCreated, types.StartTrainingResponse{SessionId: res.SessionId})
}

func (h *WorkoutHandler) HandleFinishTraining(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)
	if userId == "" {
		RespondWithError(w, http.StatusBadRequest, "JWT has not user_id")
		return
	}

	err := h.workoutClient.FinishTraining(ctx, userId)
	if err != nil {
		log.Printf("gRPC FinishTraining failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	RespondWithJSON[any](w, http.StatusOK, nil)
}

func (h *WorkoutHandler) HandleGetTrainingDays(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	yearStr := query.Get("year")
	monthStr := query.Get("month")

	if yearStr == "" || monthStr == "" {
		RespondWithError(w, http.StatusBadRequest, "Query parameters 'year' and 'month' are required")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 || year > 2100 {
		RespondWithError(w, http.StatusBadRequest, "Invalid 'year' parameter")
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		RespondWithError(w, http.StatusBadRequest, "Invalid 'month' parameter (must be between 1 and 12)")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)
	if userId == "" {
		RespondWithError(w, http.StatusBadRequest, "JWT missing user_id")
		return
	}

	grpcResponse, err := h.workoutClient.GetTrainingDays(ctx, userId, year, month)
	if err != nil {
		log.Printf("gRPC GetTrainingDays failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	daysResp := types.TrainingDaysResponse{
		Days: grpcResponse.GetTrainingDays(),
	}

	if daysResp.Days == nil {
		daysResp.Days = make([]string, 0)
	}

	RespondWithJSON(w, http.StatusOK, daysResp)
}

func (h *WorkoutHandler) HandleGetExercises(w http.ResponseWriter, r *http.Request) {
	muscleGroupId := r.PathValue("muscleGroupId")
	if muscleGroupId == "" {
		RespondWithError(w, http.StatusBadRequest, "muscleGroupId not pass")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)

	grpcResponse, err := h.workoutClient.GetExercises(ctx, userId, muscleGroupId)
	if err != nil {
		log.Printf("gRPC GetExercises failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	exercises := []types.ExerciseResponse{}

	for _, exercise := range grpcResponse.Exercises {
		exercises = append(exercises, types.ExerciseResponse{
			ExerciseId: exercise.ExerciseId,
			Type:       exercise.Type.String(),
			Name:       exercise.Name,
		})
	}

	RespondWithJSON(w, http.StatusOK, exercises)
}

func (h *WorkoutHandler) HandleGetWorkoutsForDay(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	dateStr := query.Get("date")

	if dateStr == "" {
		RespondWithError(w, http.StatusBadRequest, "Query parameters 'date' are required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)
	res, err := h.workoutClient.GetWorkoutsForDay(ctx, userId, dateStr)
	if err != nil {
		log.Printf("gRPC GetWorkoutsForDay failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response := MapGetWorkoutsForDayToHTTP(res)
	RespondWithJSON(w, http.StatusOK, response)
}

func (h *WorkoutHandler) HandleGetCompletedExercises(w http.ResponseWriter, r *http.Request) {
	exerciseId := r.PathValue("exerciseId")
	if exerciseId == "" {
		RespondWithError(w, http.StatusBadRequest, "exerciseId not pass")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)
	grpcRes, err := h.workoutClient.GetCompletedExercise(ctx, userId, exerciseId)
	if err != nil {
		log.Printf("gRPC GetCompletedExercises failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	completedExercises := make([]types.CompletedExerciseResponse, 0, len(grpcRes.GetSets()))
	for _, set := range grpcRes.GetSets() {
		item := types.CompletedExerciseResponse{
			SetNumber: int(set.GetSetNumber()),
			SetId:     set.SetId,
		}

		if set.Weight != nil {
			item.Weight = set.Weight
		}
		if set.Reps != nil {
			item.Reps = set.Reps
		}
		if set.DurationSec != nil {
			item.DurationSeconds = set.DurationSec
		}
		if set.DistanceM != nil {
			item.DistanceM = set.DistanceM
		}

		completedExercises = append(completedExercises, item)
	}

	RespondWithJSON(w, http.StatusOK, completedExercises)
}

func (h *WorkoutHandler) HandleGetSessionExercises(w http.ResponseWriter, r *http.Request) {
	sessionId := r.PathValue("sessionId")
	if sessionId == "" {
		RespondWithError(w, http.StatusBadRequest, "sessionId not pass")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcRes, err := h.workoutClient.GetSessionExercises(ctx, sessionId)
	if err != nil {
		log.Printf("gRPC GetSessionExercises failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response := MapGetSessionExercisesToHTTP(grpcRes)
	RespondWithJSON(w, http.StatusOK, response)
}

func (h *WorkoutHandler) HandleGetMuscleGroups(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcResponse, err := h.workoutClient.GetMuscleGroups(ctx)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	muscleGroups := []types.MuscleGroupsResponse{}

	for _, muscleGroup := range grpcResponse.MuscleGroups {
		muscleGroups = append(muscleGroups, types.MuscleGroupsResponse{
			MuscleGroupId: muscleGroup.MuscleGroupId,
			Code:          muscleGroup.Code,
			Name:          muscleGroup.Name,
		})
	}

	RespondWithJSON(w, http.StatusOK, muscleGroups)
}

func (h *WorkoutHandler) HandleLogWorkoutSet(w http.ResponseWriter, r *http.Request) {
	var req types.LogSetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "handle workout set data is incorrect")
		return
	}
	defer r.Body.Close()

	if err := h.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %v", err))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	tenantId := auth.TenantIDFromContext(ctx)
	if tenantId == "" {
		RespondWithError(w, http.StatusBadRequest, "JWT has not tenant_id")
		return
	}

	grpcResponse, err := h.workoutClient.LogWorkoutSet(ctx, req, tenantId)
	if err != nil {
		log.Printf("[ERROR] LogWorkoutSet gRPC call failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	RespondWithJSON(
		w,
		http.StatusOK,
		types.LogSetResponse{
			SetID:     grpcResponse.SetId,
			SetNumber: grpcResponse.SetNumber,
		},
	)
}
