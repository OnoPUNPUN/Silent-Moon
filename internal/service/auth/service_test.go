package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
	"github.com/silentmoon/api/internal/platform/security"
	authservice "github.com/silentmoon/api/internal/service/auth"
)

type mockUserRepo struct {
	users map[uuid.UUID]*domain.User
	byEmail map[string]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:   make(map[uuid.UUID]*domain.User),
		byEmail: make(map[string]*domain.User),
	}
}

func (m *mockUserRepo) Create(ctx context.Context, email, passwordHash, displayName string) (*domain.User, error) {
	u := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		DisplayName:  displayName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	m.users[u.ID] = u
	m.byEmail[u.Email] = u
	return u, nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if u, ok := m.byEmail[email]; ok {
		return u, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepo) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	if u, ok := m.users[id]; ok {
		u.PasswordHash = passwordHash
		return nil
	}
	return domain.ErrNotFound
}

func (m *mockUserRepo) UpdateProfile(ctx context.Context, id uuid.UUID, displayName string) (*domain.User, error) {
	if u, ok := m.users[id]; ok {
		u.DisplayName = displayName
		return u, nil
	}
	return nil, domain.ErrNotFound
}

type mockDeviceRepo struct {
	devices map[uuid.UUID]*domain.Device
}

func newMockDeviceRepo() *mockDeviceRepo {
	return &mockDeviceRepo{devices: make(map[uuid.UUID]*domain.Device)}
}

func (m *mockDeviceRepo) Create(ctx context.Context, userID uuid.UUID, name, platform string) (*domain.Device, error) {
	d := &domain.Device{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		Platform:  platform,
		CreatedAt: time.Now(),
	}
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
	if d, ok := m.devices[id]; ok && d.UserID == userID {
		d.Name = name
		return d, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockDeviceRepo) UpdateSyncTime(ctx context.Context, id, userID uuid.UUID, syncTime time.Time) error {
	if d, ok := m.devices[id]; ok && d.UserID == userID {
		d.LastSyncedAt = &syncTime
		return nil
	}
	return domain.ErrNotFound
}

func (m *mockDeviceRepo) Revoke(ctx context.Context, id, userID uuid.UUID) error {
	if d, ok := m.devices[id]; ok && d.UserID == userID {
		now := time.Now()
		d.RevokedAt = &now
		return nil
	}
	return domain.ErrNotFound
}

type mockTokenRepo struct {
	tokens map[string]*domain.RefreshTokenRecord
}

func newMockTokenRepo() *mockTokenRepo {
	return &mockTokenRepo{tokens: make(map[string]*domain.RefreshTokenRecord)}
}

func (m *mockTokenRepo) Create(ctx context.Context, userID, deviceID uuid.UUID, tokenHash string, expiresAt time.Time) (*domain.RefreshTokenRecord, error) {
	r := &domain.RefreshTokenRecord{
		ID:        uuid.New(),
		UserID:    userID,
		DeviceID:  &deviceID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	m.tokens[tokenHash] = r
	return r, nil
}

func (m *mockTokenRepo) GetByHash(ctx context.Context, tokenHash string) (*domain.RefreshTokenRecord, error) {
	if r, ok := m.tokens[tokenHash]; ok {
		if r.RevokedAt != nil {
			return nil, domain.ErrNotFound
		}
		if time.Now().After(r.ExpiresAt) {
			return nil, domain.ErrNotFound
		}
		return r, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockTokenRepo) Revoke(ctx context.Context, tokenHash string) error {
	if r, ok := m.tokens[tokenHash]; ok {
		now := time.Now()
		r.RevokedAt = &now
		return nil
	}
	return nil
}

func (m *mockTokenRepo) RevokeAllDeviceTokens(ctx context.Context, userID, deviceID uuid.UUID) error {
	now := time.Now()
	for _, r := range m.tokens {
		if r.UserID == userID && r.DeviceID != nil && *r.DeviceID == deviceID {
			r.RevokedAt = &now
		}
	}
	return nil
}

func TestAuthService_Register_Login_Refresh(t *testing.T) {
	userRepo := newMockUserRepo()
	devRepo := newMockDeviceRepo()
	tokenRepo := newMockTokenRepo()
	tm := security.NewTokenManager("test-secret-12345678901234567890", 15*time.Minute, 24*time.Hour)

	svc := authservice.NewService(userRepo, devRepo, tokenRepo, tm)
	ctx := context.Background()

	// 1. Register
	regPayload, err := svc.Register(ctx, authservice.RegisterRequest{
		Email:       "test@silentmoon.app",
		Password:    "supersecret123",
		DisplayName: "Test Reader",
		DeviceName:  "Pixel 9 Pro",
		Platform:    "android",
	})
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	if regPayload.User.Email != "test@silentmoon.app" {
		t.Errorf("expected email 'test@silentmoon.app', got '%s'", regPayload.User.Email)
	}
	if regPayload.Tokens.AccessToken == "" || regPayload.Tokens.RefreshToken == "" {
		t.Fatal("expected tokens to be issued")
	}

	// 2. Login
	loginPayload, err := svc.Login(ctx, authservice.LoginRequest{
		Email:      "test@silentmoon.app",
		Password:   "supersecret123",
		DeviceName: "MacBook Pro",
		Platform:   "macos",
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if loginPayload.User.ID != regPayload.User.ID {
		t.Errorf("user ID mismatch on login")
	}

	// 3. Refresh Token Rotation
	newTokens, err := svc.RefreshToken(ctx, loginPayload.Tokens.RefreshToken)
	if err != nil {
		t.Fatalf("refresh token rotation failed: %v", err)
	}
	if newTokens.AccessToken == "" || newTokens.RefreshToken == "" {
		t.Fatal("expected rotated token pair")
	}

	// Old refresh token must be revoked now
	_, err = svc.RefreshToken(ctx, loginPayload.Tokens.RefreshToken)
	if err == nil {
		t.Fatal("expected reusing old refresh token to fail")
	}
}
