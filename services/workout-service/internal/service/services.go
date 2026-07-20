package service

import "github.com/PashakArt/go_lift_backend/services/workout-service/internal/repository"

type Services struct {
	Auth        InitService
	Exercise    ExerciseService
	Session     TrainingService
	MuscleGroup MuscleGroupService
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Auth:        NewInitService(repos.User, repos.Session, repos.Tenant),
		Exercise:    NewExerciseService(repos.Exercise),
		Session:     NewTrainingService(repos.Session, repos.WorkoutSet),
		MuscleGroup: NewMuscleGroupService(repos.MuscleGroup),
	}
}
