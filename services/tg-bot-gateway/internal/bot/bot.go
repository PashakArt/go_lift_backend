package bot

import (
	"log"
	"net/http"

	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/clients/workout"
)

const (
	defaultTenantId = "00000000-0000-0000-0000-000000000000"
)

type HTTPServer struct {
	workoutClient *workout.Client
	botToken      string
}

func NewHTTPServer(botToken string, workoutClient *workout.Client) *HTTPServer {
	return &HTTPServer{
		workoutClient: workoutClient,
		botToken:      botToken,
	}
}

func (s *HTTPServer) Start(port string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/auth", s.HandleAuth)
	mux.HandleFunc("POST /api/v1/start", s.HandleStartTraining)
	mux.HandleFunc("GET /api/v1/muscle-groups", s.HandleGetMuscleGroups)
	mux.HandleFunc("GET /api/v1/{muscleGroupId}/exercises", s.HandleGetExercises)

	handler := corsMiddleware(mux)

	log.Printf("TG Gateway HTTP Server running on port :%s\n", port)

	return http.ListenAndServe(":"+port, handler)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Если это preflight-запрос, сразу отвечаем 200
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
