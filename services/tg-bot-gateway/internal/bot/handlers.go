package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/auth"
	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/types"
)

func (s *HTTPServer) HandleStartTraining(w http.ResponseWriter, r *http.Request) {
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

	res, err := s.workoutClient.StartTraining(ctx, tenantId, userId)
	if err != nil {
		log.Printf("gRPC StartTraining failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	RespondWithJSON(w, http.StatusCreated, types.StartTrainingResponse{SessionId: res.SessionId})
}

func (s *HTTPServer) HandleFinishTraining(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)
	if userId == "" {
		RespondWithError(w, http.StatusBadRequest, "JWT has not user_id")
		return
	}

	err := s.workoutClient.FinishTraining(ctx, userId)
	if err != nil {
		log.Printf("gRPC FinishTraining failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	RespondWithJSON[any](w, http.StatusOK, nil)
}

func (s *HTTPServer) HandleGetTrainingDays(w http.ResponseWriter, r *http.Request) {
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

	grpcResponse, err := s.workoutClient.GetTrainingDays(ctx, userId, year, month)
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

func (s *HTTPServer) HandleGetExercises(w http.ResponseWriter, r *http.Request) {
	muscleGroupId := r.PathValue("muscleGroupId")
	if muscleGroupId == "" {
		RespondWithError(w, http.StatusBadRequest, "muscleGroupId not pass")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)

	grpcResponse, err := s.workoutClient.GetExercises(ctx, userId, muscleGroupId)
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

func (s *HTTPServer) HandleGetWorkoutsForDay(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	dateStr := query.Get("date")

	if dateStr == "" {
		RespondWithError(w, http.StatusBadRequest, "Query parameters 'date' are required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)
	res, err := s.workoutClient.GetWorkoutsForDay(ctx, userId, dateStr)
	if err != nil {
		log.Printf("gRPC GetWorkoutsForDay failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response := MapGetWorkoutsForDayToHTTP(res)
	RespondWithJSON(w, http.StatusOK, response)
}

func (s *HTTPServer) HandleGetCompletedExercises(w http.ResponseWriter, r *http.Request) {
	exerciseId := r.PathValue("exerciseId")
	if exerciseId == "" {
		RespondWithError(w, http.StatusBadRequest, "exerciseId not pass")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)
	grpcRes, err := s.workoutClient.GetCompletedExercise(ctx, userId, exerciseId)
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

func (s *HTTPServer) HandleGetMuscleGroups(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcResponse, err := s.workoutClient.GetMuscleGroups(ctx)
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

func (s *HTTPServer) HandleLogWorkoutSet(w http.ResponseWriter, r *http.Request) {
	var req types.LogSetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "handle workout set data is incorrect")
		return
	}
	defer r.Body.Close()

	if err := s.validate.Struct(req); err != nil {
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

	grpcResponse, err := s.workoutClient.LogWorkoutSet(ctx, req, tenantId)
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

func (s *HTTPServer) HandleInit(w http.ResponseWriter, r *http.Request) {
	var req types.InitDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "init data is incorrect")
		return
	}
	defer r.Body.Close()

	var tgIDStr string

	// 🛠 CHECKS FOR LOCAL DEV MOCK DATA
	// Если мы в dev-режиме и пришли тестовые данные, пропускаем валидацию HMAC
	isDev := os.Getenv("ENV") == "local" || os.Getenv("ENV") == "development"
	if isDev && strings.Contains(req.InitData, "user=%7B%22id%22%3A") {
		log.Println("⚠️ [DEV MODE] Skipping Telegram InitData HMAC validation for mock data")

		tgIDStr = "77777"
	} else {
		// PRODUCTION LOGIC
		params, err := s.ValidateAndParseInitData(req.InitData)
		if err != nil {
			log.Printf("Telegram InitData validation failed: %v\n", err)
			RespondWithError(w, http.StatusUnauthorized, "Telegram InitData validation failed")
			return
		}

		var tgUser struct {
			ID int64 `json:"id"`
		}

		err = json.Unmarshal([]byte(params.Get("user")), &tgUser)
		if err != nil {
			log.Printf("Failed to unmarshal telegram user data: %v\n", err)
			RespondWithError(w, http.StatusBadRequest, "Failed to unmarshal telegram user data")
			return
		}

		tgIDStr = fmt.Sprintf("%d", tgUser.ID)
	}

	tenantID := req.TenantId
	if tenantID == "" {
		tenantID = defaultTenantId
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := s.workoutClient.Init(ctx, tenantID, tgIDStr)
	if err != nil {
		log.Printf("gRPC init failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	token, err := s.jwtManager.Generate(res.User.UserId, res.User.TenantId)
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	var sessionId string
	var hasActiveSession bool
	if res.ActiveSession != nil {
		sessionId = res.ActiveSession.SessionId
		hasActiveSession = true
	}

	var statusCode int
	if res.IsNewUser {
		statusCode = http.StatusCreated
	} else {
		statusCode = http.StatusOK
	}

	RespondWithJSON(w, statusCode, types.InitResponse{
		UserID:           res.User.GetUserId(),
		HasActiveSession: hasActiveSession,
		SessionId:        sessionId,
		Token:            token,
		IsNewUser:        res.IsNewUser,
		Branding:         res.TenantBranding,
	})
}

func (s *HTTPServer) HandleCreateTemplate(w http.ResponseWriter, r *http.Request) {
	var req types.CreateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	if err := s.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %v", err))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)
	if userId == "" {
		RespondWithError(w, http.StatusBadRequest, "JWT missing user_id")
		return
	}

	res, err := s.workoutClient.CreateTemplate(ctx, req, userId)
	if err != nil {
		log.Printf("[ERROR] CreateTemplate gRPC call failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	RespondWithJSON(w, http.StatusCreated, types.CreateTemplateResponse{
		TemplateID: res.TemplateId,
	})
}

func (s *HTTPServer) HandleGetTemplates(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)
	if userId == "" {
		RespondWithError(w, http.StatusBadRequest, "JWT missing user_id")
		return
	}

	grpcRes, err := s.workoutClient.GetTemplates(ctx, userId)
	if err != nil {
		log.Printf("[ERROR] GetTemplates gRPC call failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	templates := make([]types.TemplateSummaryResponse, 0, len(grpcRes.GetTemplates()))
	for _, t := range grpcRes.GetTemplates() {
		createdAt := ""
		if t.CreatedAt != nil {
			createdAt = t.CreatedAt.AsTime().Format(time.RFC3339)
		}

		templates = append(templates, types.TemplateSummaryResponse{
			TemplateID:     t.GetTemplateId(),
			Name:           t.GetName(),
			ExercisesCount: int(t.GetExercisesCount()),
			CreatedAt:      createdAt,
		})
	}

	RespondWithJSON(w, http.StatusOK, templates)
}

func (s *HTTPServer) HandleGetTemplate(w http.ResponseWriter, r *http.Request) {
	templateId := r.PathValue("templateId")
	if templateId == "" {
		RespondWithError(w, http.StatusBadRequest, "templateId is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)
	if userId == "" {
		RespondWithError(w, http.StatusBadRequest, "JWT missing user_id")
		return
	}

	grpcRes, err := s.workoutClient.GetTemplate(ctx, templateId, userId)
	if err != nil {
		log.Printf("[ERROR] GetTemplate gRPC call failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	items := make([]types.TemplateDetailItem, 0, len(grpcRes.GetItems()))
	for _, item := range grpcRes.GetItems() {
		targetSets := make([]types.TargetSet, 0, len(item.GetTargetSets()))
		for _, set := range item.GetTargetSets() {
			targetSets = append(targetSets, types.TargetSet{
				SetNum:          int32(set.GetSetNum()),
				Weight:          set.Weight,
				Reps:            set.Reps,
				DurationSeconds: set.DurationSec,
				DistanceMeters:  set.DistanceM,
			})
		}

		items = append(items, types.TemplateDetailItem{
			ExerciseID: item.GetExerciseId(),
			OrderIndex: item.GetOrderIndex(),
			TargetSets: targetSets,
			Name:       item.Name,
			Type:       item.Type.String(),
		})
	}

	createdAt := ""
	if grpcRes.CreatedAt != nil {
		createdAt = grpcRes.CreatedAt.AsTime().Format(time.RFC3339)
	}

	response := types.TemplateDetailResponse{
		TemplateID: grpcRes.GetTemplateId(),
		Name:       grpcRes.GetName(),
		CreatedAt:  createdAt,
		Items:      items,
	}

	RespondWithJSON(w, http.StatusOK, response)
}

func (s *HTTPServer) HandleDeleteTemplate(w http.ResponseWriter, r *http.Request) {
	templateId := r.PathValue("templateId")
	if templateId == "" {
		RespondWithError(w, http.StatusBadRequest, "templateId is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)
	if userId == "" {
		RespondWithError(w, http.StatusBadRequest, "JWT missing user_id")
		return
	}

	err := s.workoutClient.DeleteTemplate(ctx, templateId, userId)
	if err != nil {
		log.Printf("[ERROR] DeleteTemplate gRPC call failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	RespondWithJSON[any](w, http.StatusOK, nil)
}

func (s *HTTPServer) HandleUpdateTemplate(w http.ResponseWriter, r *http.Request) {
	templateId := r.PathValue("templateId")
	if templateId == "" {
		RespondWithError(w, http.StatusBadRequest, "templateId is required")
		return
	}

	var req types.UpdateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	if err := s.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %v", err))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)
	if userId == "" {
		RespondWithError(w, http.StatusBadRequest, "JWT missing user_id")
		return
	}

	err := s.workoutClient.UpdateTemplate(ctx, templateId, userId, req)
	if err != nil {
		log.Printf("[ERROR] UpdateTemplate gRPC call failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	RespondWithJSON[any](w, http.StatusOK, nil)
}
