package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	workoutv1 "github.com/PashakArt/go_lift_backend/api/proto/workout/v1"
	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/handlers"
	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/repository"
	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Warning: .env file not found, using system env")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Starting workout-service...")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("DATABASE_URL not found")
	}

	connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(connectCtx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	if err := pool.Ping(connectCtx); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	log.Println("Successfully connected to DB")

	userRepo := repository.NewUserRepository(pool)
	exerciseRepo := repository.NewExerciseRepository(pool)
	sessionRepo := repository.NewWorkoutSessionRepository(pool)
	muscleGroupRepo := repository.NewMuscleGroupRepository(pool)

	exerciseService := service.NewExerciseService(exerciseRepo)
	authService := service.NewAuthService(userRepo, sessionRepo)
	sessionService := service.NewWorkoutSessionService(sessionRepo)
	muscleGroupService := service.NewMuscleGroupService(muscleGroupRepo)

	workoutHandler := handlers.NewWorkoutHandler(
		authService,
		exerciseService,
		sessionService,
		muscleGroupService,
	)

	grpcPort := os.Getenv("WORKOUT_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()

	workoutv1.RegisterWorkoutServiceServer(grpcServer, workoutHandler)

	go func() {
		log.Printf("gRPC server is listening on :%s", grpcPort)
		if err := grpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gRPC server gracefully.")

	grpcServer.GracefulStop()
	log.Println("Workout-service stopped completely.")
}
