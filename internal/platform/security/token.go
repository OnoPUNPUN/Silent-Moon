package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidTokenClaims = errors.New("invalid token claims")
	ErrExpiredToken       = errors.New("token has expired")
)

type UserClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	DeviceID uuid.UUID `json:"device_id"`
	Email    string    `json:"email"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secretKey  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	issuer     string
}

func NewTokenManager(secret string, accessTTL, refreshTTL time.Duration) *TokenManager {
	return &TokenManager{
		secretKey:  []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		issuer:     "silentmoon-api",
	}
}

func (tm *TokenManager) GenerateAccessToken(userID, deviceID uuid.UUID, email string) (string, time.Time, error) {
	expiresAt := time.Now().Add(tm.accessTTL)
	claims := UserClaims{
		UserID:   userID,
		DeviceID: deviceID,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    tm.issuer,
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tm.secretKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, expiresAt, nil
}

func (tm *TokenManager) ValidateAccessToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return tm.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidTokenClaims
	}

	return claims, nil
}

func (tm *TokenManager) GenerateRefreshToken() (rawToken string, tokenHash string, expiresAt time.Time, err error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to generate random token: %w", err)
	}

	rawToken = hex.EncodeToString(bytes)
	tokenHash = HashToken(rawToken)
	expiresAt = time.Now().Add(tm.refreshTTL)

	return rawToken, tokenHash, expiresAt, nil
}

func HashToken(rawToken string) string {
	h := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(h[:])
}
