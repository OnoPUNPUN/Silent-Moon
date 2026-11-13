package sync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
)

type Service struct {
	syncRepo     domain.SyncRepository
	deviceRepo   domain.DeviceRepository
	annRepo      domain.AnnotationRepository
	bookmarkRepo domain.BookmarkRepository
	progressRepo domain.ReadingProgressRepository
	bookRepo     domain.BookRepository
	collRepo     domain.CollectionRepository
}

func NewService(
	syncRepo domain.SyncRepository,
	deviceRepo domain.DeviceRepository,
	annRepo domain.AnnotationRepository,
	bookmarkRepo domain.BookmarkRepository,
	progressRepo domain.ReadingProgressRepository,
	bookRepo domain.BookRepository,
	collRepo domain.CollectionRepository,
) *Service {
	return &Service{
		syncRepo:     syncRepo,
		deviceRepo:   deviceRepo,
		annRepo:      annRepo,
		bookmarkRepo: bookmarkRepo,
		progressRepo: progressRepo,
		bookRepo:     bookRepo,
		collRepo:     collRepo,
	}
}

func (s *Service) Push(ctx context.Context, userID uuid.UUID, req domain.SyncPushRequest) (*domain.SyncPushResponse, error) {
	// 1. Verify device exists and is not revoked
	device, err := s.deviceRepo.GetByID(ctx, req.DeviceID, userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if device.RevokedAt != nil {
		return nil, domain.ErrDeviceRevoked
	}

	// Update device sync timestamp
	_ = s.deviceRepo.UpdateSyncTime(ctx, device.ID, userID, time.Now().UTC())

	results := make([]domain.OperationResult, 0, len(req.Operations))

	for _, op := range req.Operations {
		res := s.processOperation(ctx, userID, device.ID, op)
		results = append(results, res)
	}

	return &domain.SyncPushResponse{Processed: results}, nil
}

func (s *Service) processOperation(ctx context.Context, userID, deviceID uuid.UUID, op domain.SyncOperationItem) domain.OperationResult {
	// Idempotency check: if operation was already recorded, return cached duplicate result
	prevStatus, err := s.syncRepo.GetOperation(ctx, op.OperationID, userID)
	if err == nil && prevStatus != "" {
		return domain.OperationResult{
			OperationID: op.OperationID,
			Status:      domain.OpStatusDuplicate,
			EntityType:  op.EntityType,
			EntityID:    op.EntityID,
			Message:     "operation previously processed",
		}
	}

	var result domain.OperationResult

	switch op.EntityType {
	case domain.SyncEntityAnnotation:
		result = s.handleAnnotationSync(ctx, userID, deviceID, op)
	case domain.SyncEntityBookmark:
		result = s.handleBookmarkSync(ctx, userID, deviceID, op)
	case domain.SyncEntityReadingProgress:
		result = s.handleProgressSync(ctx, userID, deviceID, op)
	default:
		result = domain.OperationResult{
			OperationID: op.OperationID,
			Status:      domain.OpStatusError,
			EntityType:  op.EntityType,
			EntityID:    op.EntityID,
			Message:     fmt.Sprintf("unsupported sync entity type: %s", op.EntityType),
		}
	}

	// Record operation outcome for idempotency
	_ = s.syncRepo.RecordOperation(ctx, op.OperationID, userID, deviceID, op.EntityType, op.EntityID, result.Status)

	return result
}

func (s *Service) handleAnnotationSync(ctx context.Context, userID, deviceID uuid.UUID, op domain.SyncOperationItem) domain.OperationResult {
	if op.Operation == domain.SyncOpDelete {
		_ = s.annRepo.SoftDelete(ctx, op.EntityID, userID)
		_, _ = s.syncRepo.CreateChange(ctx, userID, &deviceID, op.EntityType, op.EntityID, domain.SyncOpDelete)
		return domain.OperationResult{
			OperationID: op.OperationID,
			Status:      domain.OpStatusApplied,
			EntityType:  op.EntityType,
			EntityID:    op.EntityID,
		}
	}

	// Upsert
	type AnnPayload struct {
		BookID       uuid.UUID             `json:"bookId"`
		Type         domain.AnnotationType `json:"type"`
		SelectedText *string               `json:"selectedText,omitempty"`
		Color        *string               `json:"color,omitempty"`
		Note         *string               `json:"note,omitempty"`
		Location     json.RawMessage       `json:"location"`
	}

	var payload AnnPayload
	if err := json.Unmarshal(op.Data, &payload); err != nil {
		return domain.OperationResult{
			OperationID: op.OperationID,
			Status:      domain.OpStatusError,
			EntityType:  op.EntityType,
			EntityID:    op.EntityID,
			Message:     "invalid annotation data json",
		}
	}

	existing, err := s.annRepo.GetByID(ctx, op.EntityID, userID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return domain.OperationResult{
			OperationID: op.OperationID,
			Status:      domain.OpStatusError,
			EntityType:  op.EntityType,
			EntityID:    op.EntityID,
			Message:     err.Error(),
		}
	}

	if existing == nil {
		// Create new annotation with client-supplied EntityID
		ann := &domain.Annotation{
			ID:           op.EntityID,
			UserID:       userID,
			BookID:       payload.BookID,
			DeviceID:     &deviceID,
			Type:         payload.Type,
			SelectedText: payload.SelectedText,
			Color:        payload.Color,
			Note:         payload.Note,
			Location:     payload.Location,
		}

		created, err := s.annRepo.Create(ctx, ann)
		if err != nil {
			return domain.OperationResult{
				OperationID: op.OperationID,
				Status:      domain.OpStatusError,
				EntityType:  op.EntityType,
				EntityID:    op.EntityID,
				Message:     err.Error(),
			}
		}

		_, _ = s.syncRepo.CreateChange(ctx, userID, &deviceID, op.EntityType, op.EntityID, domain.SyncOpUpsert)

		return domain.OperationResult{
			OperationID:    op.OperationID,
			Status:         domain.OpStatusApplied,
			EntityType:     op.EntityType,
			EntityID:       op.EntityID,
			CurrentVersion: created.Version,
		}
	}

	// Existing entity - version check
	if op.BaseVersion < existing.Version {
		return domain.OperationResult{
			OperationID:    op.OperationID,
			Status:         domain.OpStatusConflict,
			EntityType:     op.EntityType,
			EntityID:       op.EntityID,
			CurrentVersion: existing.Version,
			Message:        "newer version exists on server",
		}
	}

	updated, err := s.annRepo.Update(ctx, op.EntityID, userID, &payload.Type, payload.SelectedText, payload.Color, payload.Note, payload.Location)
	if err != nil {
		return domain.OperationResult{
			OperationID: op.OperationID,
			Status:      domain.OpStatusError,
			EntityType:  op.EntityType,
			EntityID:    op.EntityID,
			Message:     err.Error(),
		}
	}

	_, _ = s.syncRepo.CreateChange(ctx, userID, &deviceID, op.EntityType, op.EntityID, domain.SyncOpUpsert)

	return domain.OperationResult{
		OperationID:    op.OperationID,
		Status:         domain.OpStatusApplied,
		EntityType:     op.EntityType,
		EntityID:       op.EntityID,
		CurrentVersion: updated.Version,
	}
}

func (s *Service) handleBookmarkSync(ctx context.Context, userID, deviceID uuid.UUID, op domain.SyncOperationItem) domain.OperationResult {
	if op.Operation == domain.SyncOpDelete {
		_ = s.bookmarkRepo.SoftDelete(ctx, op.EntityID, userID)
		_, _ = s.syncRepo.CreateChange(ctx, userID, &deviceID, op.EntityType, op.EntityID, domain.SyncOpDelete)
		return domain.OperationResult{
			OperationID: op.OperationID,
			Status:      domain.OpStatusApplied,
			EntityType:  op.EntityType,
			EntityID:    op.EntityID,
		}
	}

	type BmPayload struct {
		BookID   uuid.UUID       `json:"bookId"`
		Title    *string         `json:"title,omitempty"`
		Location json.RawMessage `json:"location"`
	}

	var payload BmPayload
	if err := json.Unmarshal(op.Data, &payload); err != nil {
		return domain.OperationResult{
			OperationID: op.OperationID,
			Status:      domain.OpStatusError,
			EntityType:  op.EntityType,
			EntityID:    op.EntityID,
			Message:     "invalid bookmark data json",
		}
	}

	existing, err := s.bookmarkRepo.GetByID(ctx, op.EntityID, userID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return domain.OperationResult{
			OperationID: op.OperationID,
			Status:      domain.OpStatusError,
			EntityType:  op.EntityType,
			EntityID:    op.EntityID,
			Message:     err.Error(),
		}
	}

	if existing == nil {
		bm := &domain.Bookmark{
			ID:       op.EntityID,
			UserID:   userID,
			BookID:   payload.BookID,
			DeviceID: &deviceID,
			Title:    payload.Title,
			Location: payload.Location,
		}

		created, err := s.bookmarkRepo.Create(ctx, bm)
		if err != nil {
			return domain.OperationResult{
				OperationID: op.OperationID,
				Status:      domain.OpStatusError,
				EntityType:  op.EntityType,
				EntityID:    op.EntityID,
				Message:     err.Error(),
			}
		}

		_, _ = s.syncRepo.CreateChange(ctx, userID, &deviceID, op.EntityType, op.EntityID, domain.SyncOpUpsert)

		return domain.OperationResult{
			OperationID:    op.OperationID,
			Status:         domain.OpStatusApplied,
			EntityType:     op.EntityType,
			EntityID:       op.EntityID,
			CurrentVersion: created.Version,
		}
	}

	if op.BaseVersion < existing.Version {
		return domain.OperationResult{
			OperationID:    op.OperationID,
			Status:         domain.OpStatusConflict,
			EntityType:     op.EntityType,
			EntityID:       op.EntityID,
			CurrentVersion: existing.Version,
			Message:        "newer version exists on server",
		}
	}

	updated, err := s.bookmarkRepo.Update(ctx, op.EntityID, userID, payload.Title, payload.Location)
	if err != nil {
		return domain.OperationResult{
			OperationID: op.OperationID,
			Status:      domain.OpStatusError,
			EntityType:  op.EntityType,
			EntityID:    op.EntityID,
			Message:     err.Error(),
		}
	}

	_, _ = s.syncRepo.CreateChange(ctx, userID, &deviceID, op.EntityType, op.EntityID, domain.SyncOpUpsert)

	return domain.OperationResult{
		OperationID:    op.OperationID,
		Status:         domain.OpStatusApplied,
		EntityType:     op.EntityType,
		EntityID:       op.EntityID,
		CurrentVersion: updated.Version,
	}
}

func (s *Service) handleProgressSync(ctx context.Context, userID, deviceID uuid.UUID, op domain.SyncOperationItem) domain.OperationResult {
	type ProgressPayload struct {
		BookID             uuid.UUID       `json:"bookId"`
		Location           json.RawMessage `json:"location"`
		ProgressPercentage float64         `json:"progressPercentage"`
		LastReadAt         time.Time       `json:"lastReadAt"`
	}

	var payload ProgressPayload
	if err := json.Unmarshal(op.Data, &payload); err != nil {
		return domain.OperationResult{
			OperationID: op.OperationID,
			Status:      domain.OpStatusError,
			EntityType:  op.EntityType,
			EntityID:    op.EntityID,
			Message:     "invalid reading progress data json",
		}
	}

	prog := &domain.ReadingProgress{
		UserID:             userID,
		BookID:             payload.BookID,
		DeviceID:           deviceID,
		Location:           payload.Location,
		ProgressPercentage: payload.ProgressPercentage,
		LastReadAt:         payload.LastReadAt,
	}

	_, err := s.progressRepo.Upsert(ctx, prog)
	if err != nil {
		return domain.OperationResult{
			OperationID: op.OperationID,
			Status:      domain.OpStatusError,
			EntityType:  op.EntityType,
			EntityID:    op.EntityID,
			Message:     err.Error(),
		}
	}

	_, _ = s.syncRepo.CreateChange(ctx, userID, &deviceID, op.EntityType, payload.BookID, domain.SyncOpUpsert)

	return domain.OperationResult{
		OperationID: op.OperationID,
		Status:      domain.OpStatusApplied,
		EntityType:  op.EntityType,
		EntityID:    payload.BookID,
	}
}

func (s *Service) Pull(ctx context.Context, userID, deviceID uuid.UUID, cursor int64, limit int32) (*domain.SyncPullResponse, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}

	changes, err := s.syncRepo.ListChangesAfterCursor(ctx, userID, cursor, limit+1)
	if err != nil {
		return nil, err
	}

	hasMore := false
	if len(changes) > int(limit) {
		hasMore = true
		changes = changes[:limit]
	}

	nextCursor := cursor
	changeItems := make([]domain.SyncChangeItem, 0, len(changes))

	for _, c := range changes {
		if c.ID > nextCursor {
			nextCursor = c.ID
		}

		var data any
		if c.Operation == domain.SyncOpDelete {
			data = map[string]any{
				"id":      c.EntityID,
				"deleted": true,
			}
		} else {
			// Load current entity snapshot
			switch c.EntityType {
			case domain.SyncEntityAnnotation:
				if ann, err := s.annRepo.GetByID(ctx, c.EntityID, userID); err == nil {
					data = ann
				}
			case domain.SyncEntityBookmark:
				if bm, err := s.bookmarkRepo.GetByID(ctx, c.EntityID, userID); err == nil {
					data = bm
				}
			case domain.SyncEntityReadingProgress:
				if prog, err := s.progressRepo.GetLatestForBook(ctx, c.EntityID, userID); err == nil {
					data = prog
				}
			case domain.SyncEntityBook:
				if bk, err := s.bookRepo.GetByID(ctx, c.EntityID, userID); err == nil {
					data = bk
				}
			case domain.SyncEntityCollection:
				if col, err := s.collRepo.GetByID(ctx, c.EntityID, userID); err == nil {
					data = col
				}
			}
		}

		if data != nil {
			changeItems = append(changeItems, domain.SyncChangeItem{
				Sequence:   c.ID,
				EntityType: c.EntityType,
				EntityID:   c.EntityID,
				Operation:  c.Operation,
				Data:       data,
				CreatedAt:  c.CreatedAt,
			})
		}
	}

	_ = s.deviceRepo.UpdateSyncTime(ctx, deviceID, userID, time.Now().UTC())

	return &domain.SyncPullResponse{
		Changes:    changeItems,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (s *Service) Status(ctx context.Context, userID, deviceID uuid.UUID) (*domain.SyncStatusResponse, error) {
	latestSeq, err := s.syncRepo.GetLatestSequence(ctx, userID)
	if err != nil {
		return nil, err
	}

	devices, err := s.deviceRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var lastSynced *time.Time
	for _, d := range devices {
		if d.ID == deviceID {
			lastSynced = d.LastSyncedAt
			break
		}
	}

	return &domain.SyncStatusResponse{
		LatestSequence:    latestSeq,
		LastSyncedAt:      lastSynced,
		RegisteredDevices: len(devices),
	}, nil
}
