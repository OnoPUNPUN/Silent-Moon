package device

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
)

type RegisterDeviceRequest struct {
	Name     string `json:"name"`
	Platform string `json:"platform"`
}

type UpdateDeviceRequest struct {
	Name string `json:"name"`
}

type Service struct {
	deviceRepo domain.DeviceRepository
	tokenRepo  domain.RefreshTokenRepository
}

func NewService(deviceRepo domain.DeviceRepository, tokenRepo domain.RefreshTokenRepository) *Service {
	return &Service{
		deviceRepo: deviceRepo,
		tokenRepo:  tokenRepo,
	}
}

func (s *Service) Register(ctx context.Context, userID uuid.UUID, req RegisterDeviceRequest) (*domain.Device, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("%w: device name is required", domain.ErrValidation)
	}

	platform := strings.TrimSpace(req.Platform)
	if platform == "" {
		platform = "unknown"
	}

	return s.deviceRepo.Create(ctx, userID, name, platform)
}

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]domain.Device, error) {
	return s.deviceRepo.ListByUserID(ctx, userID)
}

func (s *Service) Update(ctx context.Context, userID, deviceID uuid.UUID, req UpdateDeviceRequest) (*domain.Device, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("%w: device name is required", domain.ErrValidation)
	}

	return s.deviceRepo.Update(ctx, deviceID, userID, name)
}

func (s *Service) Revoke(ctx context.Context, userID, deviceID uuid.UUID) error {
	if err := s.deviceRepo.Revoke(ctx, deviceID, userID); err != nil {
		return err
	}
	// Revoke all refresh tokens for this device
	return s.tokenRepo.RevokeAllDeviceTokens(ctx, userID, deviceID)
}
