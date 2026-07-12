package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (s *HTTPServer) HandleGetExercises(w http.ResponseWriter, r *http.Request) {
	muscleGroupId := r.PathValue("muscleGroupId")
	if muscleGroupId == "" {
		http.Error(w, "muscleGroupId not pass", http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()
	userID := queryParams.Get("user_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcResponse, err := s.workoutClient.GetExercises(ctx, userID, muscleGroupId)
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
		w.WriteHeader(http.StatusInternalServerError)
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

func (s *HTTPServer) HandleAuth(w http.ResponseWriter, r *http.Request) {
	var req InitDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	params, err := s.ValidateAndParseInitData(req.InitData)
	if err != nil {
		log.Printf("Telegram InitData validation failed: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var tgUser struct {
		ID int64 `json:"id"`
	}

	err = json.Unmarshal([]byte(params.Get("user")), &tgUser)
	if err != nil {
		log.Printf("Failed to unmarshal telegram user data: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tgIDStr := fmt.Sprintf("%d", tgUser.ID)
	tenantID := params.Get("start_param")

	if tenantID == "" {
		tenantID = defaultTenantId
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.workoutClient.Auth(ctx, tenantID, tgIDStr)
	if err != nil {
		log.Printf("gRPC Auth failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
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
		// TODO: генерировать реальный токен
		Token: "mock-jwt-token",
	})
}
