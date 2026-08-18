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
	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		if err = godotenv.Load("../../.env"); err != nil {
			log.Println("Warning: .env file not found, using system env")
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Starting tg-bot-gateway...")

	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	appDomain := os.Getenv("APP_DOMAIN")
	webhookSecret := os.Getenv("TELEGRAM_WEBHOOK_SECRET")
	customEndpoint := os.Getenv("TELEGRAM_API_ENDPOINT")

	workoutServiceAddr := os.Getenv("WORKOUT_SERVICE_ADDR")
	if workoutServiceAddr == "" {
		workoutServiceAddr = "localhost:50051"
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	log.Printf("Connecting to workout-service at %s...\n", workoutServiceAddr)
	workoutClient, err := workout.NewClient(workoutServiceAddr)
	if err != nil {
		log.Fatalf("Failed to initialize workout gRPC client: %v", err)
	}

	var botAPI *tgbotapi.BotAPI

	if env == "local" {
		log.Println("[INFO] Running in LOCAL mode: Telegram API connection and Webhook setup skipped.")
	} else {
		if customEndpoint != "" {
			log.Printf("[INFO] Using custom Telegram API Endpoint: %s\n", customEndpoint)
			botAPI, err = tgbotapi.NewBotAPIWithAPIEndpoint(botToken, customEndpoint)
		} else {
			log.Println("[INFO] Using default Telegram API Endpoint")
			botAPI, err = tgbotapi.NewBotAPI(botToken)
		}
		if err != nil {
			log.Printf("[ERROR] Failed to connect to Telegram API (timeout/blocked): %v", err)
			log.Println("[WARNING] Gateway will work in REST mode only, Telegram commands/webhooks unavailable.")
		}

		if appDomain != "" {
			if err := telegram.SetupWebhook(botToken, appDomain, webhookSecret, customEndpoint); err != nil {
				log.Printf("[ERROR] Failed to setup webhook: %v", err)
			}
		} else {
			log.Println("[WARNING] APP_DOMAIN is empty, skipping setWebhook call")
		}
	}

	tgRouter := telegram.NewRouter(botAPI, workoutClient, webhookSecret)
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Fatalf("[ERROR] JWT_SECRET_KEY env is not set")
	}
	jwtManager := auth.NewJwtManager(jwtSecretKey, 24*time.Hour)

	server := server.NewHTTPServer(botToken, workoutClient, jwtManager, tgRouter)

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
