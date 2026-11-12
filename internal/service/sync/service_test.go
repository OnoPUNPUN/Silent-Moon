package sync_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
	syncservice "github.com/silentmoon/api/internal/service/sync"
)

// Mock Repositories
type mockSyncRepo struct {
	changes    []domain.SyncChange
	operations map[uuid.UUID]domain.OperationResultStatus
	nextSeq    int64
}

func newMockSyncRepo() *mockSyncRepo {
	return &mockSyncRepo{
		operations: make(map[uuid.UUID]domain.OperationResultStatus),
		nextSeq:    1,
	}
}

func (m *mockSyncRepo) CreateChange(ctx context.Context, userID uuid.UUID, deviceID *uuid.UUID, entityType domain.SyncEntityType, entityID uuid.UUID, op domain.SyncOperationType) (*domain.SyncChange, error) {
	c := domain.SyncChange{
		ID:         m.nextSeq,
		UserID:     userID,
		DeviceID:   deviceID,
		EntityType: entityType,
		EntityID:   entityID,
		Operation:  op,
		CreatedAt:  time.Now(),
	}
	m.nextSeq++
	m.changes = append(m.changes, c)
	return &c, nil
}

func (m *mockSyncRepo) ListChangesAfterCursor(ctx context.Context, userID uuid.UUID, cursor int64, limit int32) ([]domain.SyncChange, error) {
	var res []domain.SyncChange
	for _, c := range m.changes {
		if c.UserID == userID && c.ID > cursor {
			res = append(res, c)
			if int32(len(res)) == limit {
				break
			}
		}
	}
	return res, nil
}

func (m *mockSyncRepo) GetLatestSequence(ctx context.Context, userID uuid.UUID) (int64, error) {
	var max int64
	for _, c := range m.changes {
		if c.UserID == userID && c.ID > max {
			max = c.ID
		}
	}
	return max, nil
}

func (m *mockSyncRepo) RecordOperation(ctx context.Context, opID, userID, deviceID uuid.UUID, entityType domain.SyncEntityType, entityID uuid.UUID, status domain.OperationResultStatus) error {
	m.operations[opID] = status
	return nil
}

func (m *mockSyncRepo) GetOperation(ctx context.Context, opID, userID uuid.UUID) (domain.OperationResultStatus, error) {
	if status, ok := m.operations[opID]; ok {
		return status, nil
	}
	return "", domain.ErrNotFound
}

type mockDeviceRepo struct {
	devices map[uuid.UUID]*domain.Device
}

func (m *mockDeviceRepo) Create(ctx context.Context, userID uuid.UUID, name, platform string) (*domain.Device, error) {
	d := &domain.Device{ID: uuid.New(), UserID: userID, Name: name, Platform: platform, CreatedAt: time.Now()}
	m.devices[d.ID] = d
	return d, nil
}
func (m *mockDeviceRepo) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Device, error) {
	if d, ok := m.devices[id]; ok && d.UserID == userID {
		return d, nil
	}
	return nil, domain.ErrNotFound
}
func (m *mockDeviceRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Device, error) {
	var list []domain.Device
	for _, d := range m.devices {
		if d.UserID == userID {
			list = append(list, *d)
		}
	}
	return list, nil
}
func (m *mockDeviceRepo) Update(ctx context.Context, id, userID uuid.UUID, name string) (*domain.Device, error) {
	return nil, nil
}
func (m *mockDeviceRepo) UpdateSyncTime(ctx context.Context, id, userID uuid.UUID, syncTime time.Time) error {
	if d, ok := m.devices[id]; ok && d.UserID == userID {
		d.LastSyncedAt = &syncTime
	}
	return nil
}
func (m *mockDeviceRepo) Revoke(ctx context.Context, id, userID uuid.UUID) error {
	if d, ok := m.devices[id]; ok && d.UserID == userID {
		now := time.Now()
		d.RevokedAt = &now
	}
	return nil
}

type mockAnnRepo struct {
	annotations map[uuid.UUID]*domain.Annotation
}

func (m *mockAnnRepo) Create(ctx context.Context, ann *domain.Annotation) (*domain.Annotation, error) {
	a := *ann
	a.Version = 1
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	m.annotations[a.ID] = &a
	return &a, nil
}
func (m *mockAnnRepo) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Annotation, error) {
	if a, ok := m.annotations[id]; ok && a.UserID == userID && a.DeletedAt == nil {
		return a, nil
	}
	return nil, domain.ErrNotFound
}
func (m *mockAnnRepo) ListByBookID(ctx context.Context, bookID, userID uuid.UUID) ([]domain.Annotation, error) {
	return nil, nil
}
func (m *mockAnnRepo) Update(ctx context.Context, id, userID uuid.UUID, annType *domain.AnnotationType, selectedText, color, note *string, location json.RawMessage) (*domain.Annotation, error) {
	if a, ok := m.annotations[id]; ok && a.UserID == userID && a.DeletedAt == nil {
		if annType != nil {
			a.Type = *annType
		}
		if selectedText != nil {
			a.SelectedText = selectedText
		}
		if color != nil {
			a.Color = color
		}
		if note != nil {
			a.Note = note
		}
		if len(location) > 0 {
			a.Location = location
		}
		a.Version++
		a.UpdatedAt = time.Now()
		return a, nil
	}
	return nil, domain.ErrNotFound
}
func (m *mockAnnRepo) SoftDelete(ctx context.Context, id, userID uuid.UUID) error {
	if a, ok := m.annotations[id]; ok && a.UserID == userID && a.DeletedAt == nil {
		now := time.Now()
		a.DeletedAt = &now
		a.Version++
		return nil
	}
	return domain.ErrNotFound
}

func TestMultiDeviceSyncFlow(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	bookID := uuid.New()

	devA := &domain.Device{ID: uuid.New(), UserID: userID, Name: "Android Tablet"}
	devB := &domain.Device{ID: uuid.New(), UserID: userID, Name: "Desktop Mac"}

	syncRepo := newMockSyncRepo()
	deviceRepo := &mockDeviceRepo{devices: map[uuid.UUID]*domain.Device{devA.ID: devA, devB.ID: devB}}
	annRepo := &mockAnnRepo{annotations: make(map[uuid.UUID]*domain.Annotation)}

	svc := syncservice.NewService(syncRepo, deviceRepo, annRepo, nil, nil, nil, nil)

	annotationID := uuid.New()
	op1ID := uuid.New()

	// 1. Device A pushes an annotation creation
	loc := json.RawMessage(`{"format":"epub","cfi":"epubcfi(/6/4[chap1]!/4/2/10:5)"}`)
	txt := "First synchronized highlight"
	col := "#FFFF00"

	annData, _ := json.Marshal(map[string]any{
		"bookId":       bookID,
		"type":         "highlight",
		"selectedText": txt,
		"color":        col,
		"location":     loc,
	})

	pushA, err := svc.Push(ctx, userID, domain.SyncPushRequest{
		DeviceID: devA.ID,
		Operations: []domain.SyncOperationItem{
			{
				OperationID: op1ID,
				EntityType:  domain.SyncEntityAnnotation,
				EntityID:    annotationID,
				Operation:   domain.SyncOpUpsert,
				BaseVersion: 0,
				Data:        annData,
			},
		},
	})
	if err != nil {
		t.Fatalf("Device A push failed: %v", err)
	}
	if len(pushA.Processed) != 1 || pushA.Processed[0].Status != domain.OpStatusApplied {
		t.Fatalf("expected applied operation, got: %+v", pushA.Processed)
	}
	if pushA.Processed[0].CurrentVersion != 1 {
		t.Fatalf("expected version 1, got %d", pushA.Processed[0].CurrentVersion)
	}

	// 2. Retry of op1 must be detected as duplicate (Idempotency)
	pushRetry, err := svc.Push(ctx, userID, domain.SyncPushRequest{
		DeviceID: devA.ID,
		Operations: []domain.SyncOperationItem{
			{
				OperationID: op1ID,
				EntityType:  domain.SyncEntityAnnotation,
				EntityID:    annotationID,
				Operation:   domain.SyncOpUpsert,
				BaseVersion: 0,
				Data:        annData,
			},
		},
	})
	if err != nil {
		t.Fatalf("Device A retry push failed: %v", err)
	}
	if pushRetry.Processed[0].Status != domain.OpStatusDuplicate {
		t.Fatalf("expected duplicate status on retry, got: %s", pushRetry.Processed[0].Status)
	}

	// 3. Device B pulls remote changes from cursor 0
	pullB, err := svc.Pull(ctx, userID, devB.ID, 0, 100)
	if err != nil {
		t.Fatalf("Device B pull failed: %v", err)
	}
	if len(pullB.Changes) != 1 {
		t.Fatalf("Device B expected 1 change, got %d", len(pullB.Changes))
	}
	if pullB.Changes[0].EntityID != annotationID {
		t.Fatalf("unexpected change entity ID: %s", pullB.Changes[0].EntityID)
	}
	devBCursor := pullB.NextCursor

	// 4. Device B modifies the annotation (now at version 1 -> 2)
	op2ID := uuid.New()
	newNote := "Device B added a note here"
	updatedAnnData, _ := json.Marshal(map[string]any{
		"bookId":       bookID,
		"type":         "highlight",
		"selectedText": txt,
		"color":        col,
		"note":         newNote,
		"location":     loc,
	})

	pushB, err := svc.Push(ctx, userID, domain.SyncPushRequest{
		DeviceID: devB.ID,
		Operations: []domain.SyncOperationItem{
			{
				OperationID: op2ID,
				EntityType:  domain.SyncEntityAnnotation,
				EntityID:    annotationID,
				Operation:   domain.SyncOpUpsert,
				BaseVersion: 1, // Device B observes version 1
				Data:        updatedAnnData,
			},
		},
	})
	if err != nil {
		t.Fatalf("Device B update failed: %v", err)
	}
	if pushB.Processed[0].Status != domain.OpStatusApplied || pushB.Processed[0].CurrentVersion != 2 {
		t.Fatalf("expected version 2, got: %+v", pushB.Processed[0])
	}

	// 5. Device A attempts to modify outdated version (BaseVersion 1 while server is at 2) -> Conflict
	op3ID := uuid.New()
	conflictingNote := "Device A conflicting edit"
	conflictData, _ := json.Marshal(map[string]any{
		"bookId":       bookID,
		"type":         "highlight",
		"note":         conflictingNote,
		"location":     loc,
	})

	pushConflict, err := svc.Push(ctx, userID, domain.SyncPushRequest{
		DeviceID: devA.ID,
		Operations: []domain.SyncOperationItem{
			{
				OperationID: op3ID,
				EntityType:  domain.SyncEntityAnnotation,
				EntityID:    annotationID,
				Operation:   domain.SyncOpUpsert,
				BaseVersion: 1, // Outdated! Server is now at 2
				Data:        conflictData,
			},
		},
	})
	if err != nil {
		t.Fatalf("push conflict request failed: %v", err)
	}
	if pushConflict.Processed[0].Status != domain.OpStatusConflict {
		t.Fatalf("expected conflict status, got %s", pushConflict.Processed[0].Status)
	}

	// 6. Device A pulls and catches up
	pullA, err := svc.Pull(ctx, userID, devA.ID, 1, 100)
	if err != nil {
		t.Fatalf("Device A pull failed: %v", err)
	}
	if len(pullA.Changes) != 1 {
		t.Fatalf("Device A expected 1 new change, got %d", len(pullA.Changes))
	}

	// 7. Device A deletes the annotation -> tombstone
	op4ID := uuid.New()
	pushDelete, err := svc.Push(ctx, userID, domain.SyncPushRequest{
		DeviceID: devA.ID,
		Operations: []domain.SyncOperationItem{
			{
				OperationID: op4ID,
				EntityType:  domain.SyncEntityAnnotation,
				EntityID:    annotationID,
				Operation:   domain.SyncOpDelete,
			},
		},
	})
	if err != nil {
		t.Fatalf("delete push failed: %v", err)
	}
	if pushDelete.Processed[0].Status != domain.OpStatusApplied {
		t.Fatalf("expected delete to be applied, got %s", pushDelete.Processed[0].Status)
	}

	// 8. Device B pulls after devBCursor and receives delete tombstone
	pullBTombstone, err := svc.Pull(ctx, userID, devB.ID, devBCursor, 100)
	if err != nil {
		t.Fatalf("Device B pull tombstone failed: %v", err)
	}
	// Changes since devBCursor: 1 update from B (seq 2), 1 delete from A (seq 3)
	foundTombstone := false
	for _, ch := range pullBTombstone.Changes {
		if ch.EntityID == annotationID && ch.Operation == domain.SyncOpDelete {
			foundTombstone = true
			break
		}
	}
	if !foundTombstone {
		t.Fatalf("Device B failed to receive delete tombstone")
	}
}
