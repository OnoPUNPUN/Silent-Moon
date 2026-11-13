package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Bookmark struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"userId"`
	BookID    uuid.UUID       `json:"bookId"`
	DeviceID  *uuid.UUID      `json:"deviceId,omitempty"`
	Title     *string         `json:"title,omitempty"`
	Location  json.RawMessage `json:"location"`
	Version   int64           `json:"version"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	DeletedAt *time.Time      `json:"deletedAt,omitempty"`
}
