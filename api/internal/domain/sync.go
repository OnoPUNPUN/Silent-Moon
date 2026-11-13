package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SyncOperationType string

const (
	SyncOpUpsert SyncOperationType = "upsert"
	SyncOpDelete SyncOperationType = "delete"
)

type SyncEntityType string

const (
	SyncEntityAnnotation      SyncEntityType = "annotation"
	SyncEntityBookmark        SyncEntityType = "bookmark"
	SyncEntityReadingProgress SyncEntityType = "reading_progress"
	SyncEntityBook            SyncEntityType = "book"
	SyncEntityCollection      SyncEntityType = "collection"
	SyncEntityTag             SyncEntityType = "tag"
)

type SyncOperationItem struct {
	OperationID uuid.UUID         `json:"operationId"`
	EntityType  SyncEntityType    `json:"entityType"`
	EntityID    uuid.UUID         `json:"entityId"`
	Operation   SyncOperationType `json:"operation"`
	BaseVersion int64             `json:"baseVersion"`
	Data        json.RawMessage   `json:"data"`
}

type SyncPushRequest struct {
	DeviceID   uuid.UUID           `json:"deviceId"`
	Operations []SyncOperationItem `json:"operations"`
}

type OperationResultStatus string

const (
	OpStatusApplied   OperationResultStatus = "applied"
	OpStatusConflict  OperationResultStatus = "conflict"
	OpStatusDuplicate OperationResultStatus = "duplicate"
	OpStatusError     OperationResultStatus = "error"
)

type OperationResult struct {
	OperationID    uuid.UUID             `json:"operationId"`
	Status         OperationResultStatus `json:"status"`
	EntityType     SyncEntityType        `json:"entityType"`
	EntityID       uuid.UUID             `json:"entityId"`
	CurrentVersion int64                 `json:"currentVersion,omitempty"`
	Message        string                `json:"message,omitempty"`
}

type SyncPushResponse struct {
	Processed []OperationResult `json:"processed"`
}

type SyncChangeItem struct {
	Sequence   int64             `json:"sequence"`
	EntityType SyncEntityType    `json:"entityType"`
	EntityID   uuid.UUID         `json:"entityId"`
	Operation  SyncOperationType `json:"operation"`
	Data       any               `json:"data"`
	CreatedAt  time.Time         `json:"createdAt"`
}

type SyncPullResponse struct {
	Changes    []SyncChangeItem `json:"changes"`
	NextCursor int64            `json:"nextCursor"`
	HasMore    bool             `json:"hasMore"`
}

type SyncStatusResponse struct {
	LatestSequence    int64     `json:"latestSequence"`
	LastSyncedAt      *time.Time `json:"lastSyncedAt,omitempty"`
	RegisteredDevices int       `json:"registeredDevices"`
}

type SyncChange struct {
	ID         int64             `json:"id"`
	UserID     uuid.UUID         `json:"userId"`
	DeviceID   *uuid.UUID        `json:"deviceId,omitempty"`
	EntityType SyncEntityType    `json:"entityType"`
	EntityID   uuid.UUID         `json:"entityId"`
	Operation  SyncOperationType `json:"operation"`
	CreatedAt  time.Time         `json:"createdAt"`
}
