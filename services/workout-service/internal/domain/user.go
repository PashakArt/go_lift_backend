package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID     uuid.UUID `json:"user_id"`
	TenantID   uuid.UUID `json:"tenant_id,omitempty"`
	TelegramID string    `json:"telegram_id"`
	Phone      string    `json:"phone,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
