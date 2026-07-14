package domain

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	TenantID     uuid.UUID `json:"tenant_id"`
	Name         string    `json:"name"`
	BrandingJSON []byte    `json:"branding"`
	CreatedAt    time.Time `json:"created_at"`
}
