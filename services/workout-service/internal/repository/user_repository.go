package repository

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	//go:embed queries/create_user.sql
	createUserQuery string

	//go:embed queries/get_user_by_tg_id.sql
	getUserByTgId string
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) domain.UserRepository {
	return &userRepository{pool}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	if user.UserID == uuid.Nil {
		user.UserID = uuid.New()
	}

	err := r.pool.QueryRow(
		ctx,
		createUserQuery,
		user.UserID,
		user.TenantID,
		user.TelegramID,
		user.Phone,
		user.CreatedAt,
	).Scan(&user.UserID, &user.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to execute create user query: %w", err)
	}

	return nil
}

func (r *userRepository) GetByTenantAndTelegramID(ctx context.Context, tenantId uuid.UUID, tgId string) (*domain.User, error) {
	var user domain.User

	err := r.pool.QueryRow(ctx, getUserByTgId, tenantId, tgId).Scan(
		&user.UserID,
		&user.TenantID,
		&user.TelegramID,
		&user.Phone,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to execute get user by tg id query: %w", err)
	}

	return &user, nil
}
