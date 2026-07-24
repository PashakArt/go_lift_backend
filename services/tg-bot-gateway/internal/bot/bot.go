package bot

import (
	"log"
	"net/http"

	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/auth"
	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/clients/workout"
	"github.com/go-playground/validator/v10"
)

const (
	defaultTenantId = "00000000-0000-0000-0000-000000000000"
)

type HTTPServer struct {
	workoutClient *workout.Client
	botToken      string
	jwtManager    *auth.JwtManager
	validate      *validator.Validate
}

func NewHTTPServer(
	botToken string,
	workoutClient *workout.Client,
	jwtManager *auth.JwtManager,
) *HTTPServer {
	return &HTTPServer{
		workoutClient: workoutClient,
		botToken:      botToken,
		jwtManager:    jwtManager,
		validate:      validator.New(),
	}
}

func (s *HTTPServer) Start(port string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/init", s.HandleInit)

	mux.Handle("POST /api/v1/workout/sets", s.AuthMiddleware(http.HandlerFunc(s.HandleLogWorkoutSet)))

	mux.Handle("POST /api/v1/start", s.AuthMiddleware(http.HandlerFunc(s.HandleStartTraining)))
	mux.Handle("POST /api/v1/finish", s.AuthMiddleware(http.HandlerFunc(s.HandleFinishTraining)))

	mux.Handle("GET /api/v1/muscle-groups", s.AuthMiddleware(http.HandlerFunc(s.HandleGetMuscleGroups)))
	mux.Handle("GET /api/v1/{muscleGroupId}/exercises", s.AuthMiddleware(http.HandlerFunc(s.HandleGetExercises)))
	mux.Handle(
		"GET /api/v1/exercises/{exerciseId}/completed",
		s.AuthMiddleware(http.HandlerFunc(s.HandleGetCompletedExercises)),
	)

	mux.Handle("GET /api/v1/workouts/calendar", s.AuthMiddleware(http.HandlerFunc(s.HandleGetTrainingDays)))

	handler := CorsMiddleware(mux)

	log.Printf("TG Gateway HTTP Server running on port :%s\n", port)

	return http.ListenAndServe(":"+port, handler)
}
