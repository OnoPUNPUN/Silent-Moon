package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
	"github.com/silentmoon/api/internal/platform/security"
)

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
	DeviceName  string `json:"deviceName"`
	Platform    string `json:"platform"`
}

type LoginRequest struct {
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	DeviceName string     `json:"deviceName"`
	Platform   string     `json:"platform"`
	DeviceID   *uuid.UUID `json:"deviceId,omitempty"`
}

type Service struct {
	userRepo   domain.UserRepository
	deviceRepo domain.DeviceRepository
	tokenRepo  domain.RefreshTokenRepository
	tokens     *security.TokenManager
}

func NewService(
	userRepo domain.UserRepository,
	deviceRepo domain.DeviceRepository,
	tokenRepo domain.RefreshTokenRepository,
	tokens *security.TokenManager,
) *Service {
	return &Service{
		userRepo:   userRepo,
		deviceRepo: deviceRepo,
		tokenRepo:  tokenRepo,
		tokens:     tokens,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*domain.AuthPayload, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.DisplayName = strings.TrimSpace(req.DisplayName)
	req.DeviceName = strings.TrimSpace(req.DeviceName)
	req.Platform = strings.TrimSpace(req.Platform)

	if req.Email == "" || !strings.Contains(req.Email, "@") {
		return nil, fmt.Errorf("%w: valid email is required", domain.ErrValidation)
	}
	if len(req.Password) < 8 {
		return nil, fmt.Errorf("%w: password must be at least 8 characters", domain.ErrValidation)
	}
	if req.DisplayName == "" {
		req.DisplayName = strings.Split(req.Email, "@")[0]
	}
	if req.DeviceName == "" {
		req.DeviceName = "Primary Device"
	}
	if req.Platform == "" {
		req.Platform = "unknown"
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	hash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.userRepo.Create(ctx, req.Email, hash, req.DisplayName)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	device, err := s.deviceRepo.Create(ctx, user.ID, req.DeviceName, req.Platform)
	if err != nil {
		return nil, fmt.Errorf("failed to register device: %w", err)
	}

	tokenPair, err := s.generateAndStoreTokens(ctx, user.ID, device.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &domain.AuthPayload{
		User:   *user,
		Device: *device,
		Tokens: *tokenPair,
	}, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*domain.AuthPayload, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		return nil, domain.ErrInvalidCredentials
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if !security.CheckPassword(req.Password, user.PasswordHash) {
		return nil, domain.ErrInvalidCredentials
	}

	var device *domain.Device
	if req.DeviceID != nil && *req.DeviceID != uuid.Nil {
		d, err := s.deviceRepo.GetByID(ctx, *req.DeviceID, user.ID)
		if err == nil && d != nil {
			if d.RevokedAt != nil {
				return nil, domain.ErrDeviceRevoked
			}
			device = d
		}
	}

	if device == nil {
		devName := strings.TrimSpace(req.DeviceName)
		if devName == "" {
			devName = "Registered Device"
		}
		platform := strings.TrimSpace(req.Platform)
		if platform == "" {
			platform = "unknown"
		}

		newDevice, err := s.deviceRepo.Create(ctx, user.ID, devName, platform)
		if err != nil {
			return nil, fmt.Errorf("failed to register device: %w", err)
		}
		device = newDevice
	}

	tokenPair, err := s.generateAndStoreTokens(ctx, user.ID, device.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &domain.AuthPayload{
		User:   *user,
		Device: *device,
		Tokens: *tokenPair,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, rawRefreshToken string) (*domain.TokenPair, error) {
	rawRefreshToken = strings.TrimSpace(rawRefreshToken)
	if rawRefreshToken == "" {
		return nil, domain.ErrInvalidToken
	}

	tokenHash := security.HashToken(rawRefreshToken)
	record, err := s.tokenRepo.GetByHash(ctx, tokenHash)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	if record.DeviceID == nil {
		return nil, domain.ErrInvalidToken
	}

	// Verify device is still active
	device, err := s.deviceRepo.GetByID(ctx, *record.DeviceID, record.UserID)
	if err != nil || device.RevokedAt != nil {
		return nil, domain.ErrDeviceRevoked
	}

	user, err := s.userRepo.GetByID(ctx, record.UserID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	// Revoke old refresh token (rotation)
	_ = s.tokenRepo.Revoke(ctx, tokenHash)

	// Issue new token pair
	return s.generateAndStoreTokens(ctx, user.ID, device.ID, user.Email)
}

func (s *Service) Logout(ctx context.Context, rawRefreshToken string) error {
	rawRefreshToken = strings.TrimSpace(rawRefreshToken)
	if rawRefreshToken == "" {
		return nil
	}

	tokenHash := security.HashToken(rawRefreshToken)
	return s.tokenRepo.Revoke(ctx, tokenHash)
}

func (s *Service) GetMe(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *Service) generateAndStoreTokens(ctx context.Context, userID, deviceID uuid.UUID, email string) (*domain.TokenPair, error) {
	accessToken, expiresAt, err := s.tokens.GenerateAccessToken(userID, deviceID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	rawRefreshToken, tokenHash, refreshExpiresAt, err := s.tokens.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	_, err = s.tokenRepo.Create(ctx, userID, deviceID, tokenHash, refreshExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to persist refresh token: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		ExpiresAt:    expiresAt,
		TokenType:    "Bearer",
	}, nil
}
