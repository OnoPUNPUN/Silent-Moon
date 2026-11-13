package security_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/platform/security"
)

func TestPasswordHashing(t *testing.T) {
	password := "my-secret-password-123"

	hash, err := security.HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	if hash == password {
		t.Fatal("hash should not equal plaintext password")
	}

	if !security.CheckPassword(password, hash) {
		t.Fatal("password verification failed for correct password")
	}

	if security.CheckPassword("wrong-password", hash) {
		t.Fatal("password verification succeeded for incorrect password")
	}
}

func TestTokenManager(t *testing.T) {
	tm := security.NewTokenManager("test-secret-key-1234567890123456", 5*time.Minute, 1*time.Hour)

	userID := uuid.New()
	deviceID := uuid.New()
	email := "reader@example.com"

	// Access token
	tokenStr, expiresAt, err := tm.GenerateAccessToken(userID, deviceID, email)
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	if tokenStr == "" {
		t.Fatal("token string should not be empty")
	}

	if !expiresAt.After(time.Now()) {
		t.Fatal("expiration should be in the future")
	}

	claims, err := tm.ValidateAccessToken(tokenStr)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected user ID %s, got %s", userID, claims.UserID)
	}
	if claims.DeviceID != deviceID {
		t.Errorf("expected device ID %s, got %s", deviceID, claims.DeviceID)
	}
	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}

	// Refresh token
	rawRefresh, hash, refreshExp, err := tm.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}
	if rawRefresh == "" || hash == "" {
		t.Fatal("refresh token and hash must not be empty")
	}
	if !refreshExp.After(time.Now()) {
		t.Fatal("refresh expiration should be in the future")
	}

	computedHash := security.HashToken(rawRefresh)
	if computedHash != hash {
		t.Errorf("expected hash %s, got %s", hash, computedHash)
	}
}
