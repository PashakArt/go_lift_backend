package server

import (
	"log"
	"net/http"

	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/server/handlers"
	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/telegram"
)

type HTTPServer struct {
	workoutHandler  *handlers.WorkoutHandler
	authHandler     *handlers.AuthHandler
	templateHandler *handlers.TemplateHandler
	authMiddleware  func(http.Handler) http.Handler
	telegramRouter  *telegram.Router
}

func NewHTTPServer(
	workoutHandler *handlers.WorkoutHandler,
	authHandler *handlers.AuthHandler,
	telegramRouter *telegram.Router,
	authMiddleware func(http.Handler) http.Handler,
) *HTTPServer {
	return &HTTPServer{
		workoutHandler: workoutHandler,
		authHandler:    authHandler,
		telegramRouter: telegramRouter,
		authMiddleware: authMiddleware,
	}
}

func (s *HTTPServer) Start(port string) error {
	mux := http.NewServeMux()

	s.authHandler.RegisterRoutes(mux)
	s.workoutHandler.RegisterRoutes(mux, s.authMiddleware)
	s.templateHandler.RegisterRoutes(mux, s.authMiddleware)

	mux.HandleFunc("POST /api/v1/telegram/webhook", s.telegramRouter.HandleWebhook)

	handler := CorsMiddleware(mux)

	log.Printf("TG Gateway HTTP Server running on port :%s\n", port)

	return http.ListenAndServe(":"+port, handler)
}
