package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ReadingProgress struct {
	ID                 uuid.UUID       `json:"id"`
	UserID             uuid.UUID       `json:"userId"`
	BookID             uuid.UUID       `json:"bookId"`
	DeviceID           uuid.UUID       `json:"deviceId"`
	Location           json.RawMessage `json:"location"`
	ProgressPercentage float64         `json:"progressPercentage"`
	LastReadAt         time.Time       `json:"lastReadAt"`
	UpdatedAt          time.Time       `json:"updatedAt"`
}
