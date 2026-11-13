package domain

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"userId"`
	Name         string     `json:"name"`
	Platform     string     `json:"platform"`
	LastSyncedAt *time.Time `json:"lastSyncedAt,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	RevokedAt    *time.Time `json:"revokedAt,omitempty"`
}
