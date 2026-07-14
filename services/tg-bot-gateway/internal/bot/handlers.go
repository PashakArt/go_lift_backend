package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/auth"
)

func (s *HTTPServer) HandleStartTraining(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.workoutClient.StartTraining(ctx)
	if err != nil {
		log.Printf("gRPC StartTraining failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	RespondWithJSON(w, http.StatusCreated, nil)
}

func (s *HTTPServer) HandleGetExercises(w http.ResponseWriter, r *http.Request) {
	muscleGroupId := r.PathValue("muscleGroupId")
	if muscleGroupId == "" {
		RespondWithError(w, http.StatusBadRequest, "muscleGroupId not pass")
		return
	}

	userId := auth.UserIDFromContext(r.Context())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcResponse, err := s.workoutClient.GetExercises(ctx, userId, muscleGroupId)
	if err != nil {
		log.Printf("gRPC GetExercises failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	exercises := []ExerciseResponse{}

	for _, exercise := range grpcResponse.Exercises {
		exercises = append(exercises, ExerciseResponse{
			ExerciseId: exercise.ExerciseId,
			Type:       exercise.Type.String(),
			Name:       exercise.Name,
		})
	}

	RespondWithJSON(w, http.StatusOK, exercises)
}

func (s *HTTPServer) HandleGetMuscleGroups(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcResponse, err := s.workoutClient.GetMuscleGroups(ctx)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	muscleGroups := []MuscleGroupsResponse{}

	for _, muscleGroup := range grpcResponse.MuscleGroups {
		muscleGroups = append(muscleGroups, MuscleGroupsResponse{
			MuscleGroupId: muscleGroup.MuscleGroupId,
			Code:          muscleGroup.Code,
			Name:          muscleGroup.Name,
		})
	}

	RespondWithJSON(w, http.StatusOK, muscleGroups)
}

func (s *HTTPServer) HandleInit(w http.ResponseWriter, r *http.Request) {
	log.Printf("123")
	var req InitDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "init data is incorrect")
		return
	}
	defer r.Body.Close()

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

	tgIDStr := fmt.Sprintf("%d", tgUser.ID)
	tenantID := req.TenantId

	if tenantID == "" {
		tenantID = defaultTenantId
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	RespondWithJSON(w, statusCode, AuthResponse{
		UserID:           res.User.GetUserId(),
		HasActiveSession: hasActiveSession,
		SessionId:        sessionId,
		Token:            token,
	})
}
