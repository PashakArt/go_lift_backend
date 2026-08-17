package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/auth"
	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/clients/workout"
	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Warning: .env file not found, using system env")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Starting tg-bot-gateway...")

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if botToken == "" {
		// TODO Пока оставим предупреждение. Для локальных тестов без валидации хэша
		// можно будет передавать заглушку, но для честной проверки токен нужен.
		log.Println("[WARNING] TELEGRAM_BOT_TOKEN is not set")
	}

	workoutServiceAddr := os.Getenv("WORKOUT_SERVICE_ADDR")
	if workoutServiceAddr == "" {
		workoutServiceAddr = "localhost:50051" // Дефолтный локальный порт workout-service
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080" // Порт, который будет слушать наш шлюз для фронтенда
	}

	log.Printf("Connecting to workout-service at %s...\n", workoutServiceAddr)
	workoutClient, err := workout.NewClient(workoutServiceAddr)
	if err != nil {
		log.Fatalf("Failed to initialize workout gRPC client: %v", err)
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Fatalf("[ERROR] JWT_SECRET_KEY env is not set")
	}
	jwtManager := auth.NewJwtManager(jwtSecretKey, 24*time.Hour)

	server := server.NewHTTPServer(botToken, workoutClient, jwtManager)

	go func() {
		if err := server.Start(httpPort); err != nil {
			log.Fatalf("HTTP server failed to start: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down tg-bot-gateway...")

	_, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	log.Println("Tg-bot-gateway stopped completely.")
}
