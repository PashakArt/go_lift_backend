package service

import "github.com/PashakArt/go_lift_backend/services/workout-service/internal/repository"

type Services struct {
	Auth        AuthService
	Exercise    ExerciseService
	Session     WorkoutSessionService
	MuscleGroup MuscleGroupService
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Auth:        NewAuthService(repos.User, repos.Session),
		Exercise:    NewExerciseService(repos.Exercise),
		Session:     NewWorkoutSessionService(repos.Session),
		MuscleGroup: NewMuscleGroupService(repos.MuscleGroup),
	}
}
