package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/silentmoon/api/internal/domain"
	"github.com/silentmoon/api/internal/middleware"
	"github.com/silentmoon/api/internal/platform/response"
	authservice "github.com/silentmoon/api/internal/service/auth"
)

type Handler struct {
	service *authservice.Service
}

func NewHandler(service *authservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req authservice.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid request body")
		return
	}

	res, err := h.service.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			response.ValidationError(w, err.Error())
			return
		}
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			response.Conflict(w, "a user with this email already exists")
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, res)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req authservice.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid request body")
		return
	}

	res, err := h.service.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.Unauthorized(w, "invalid email or password")
			return
		}
		if errors.Is(err, domain.ErrDeviceRevoked) {
			response.Forbidden(w, "device has been revoked")
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, res)
}

type RefreshTokenBody struct {
	RefreshToken string `json:"refreshToken"`
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body RefreshTokenBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RefreshToken == "" {
		response.ValidationError(w, "valid refreshToken is required")
		return
	}

	tokens, err := h.service.RefreshToken(r.Context(), body.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrDeviceRevoked) {
			response.Forbidden(w, "device has been revoked")
			return
		}
		response.Unauthorized(w, "invalid or expired refresh token")
		return
	}

	response.JSON(w, http.StatusOK, tokens)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var body RefreshTokenBody
	_ = json.NewDecoder(r.Body).Decode(&body)

	_ = h.service.Logout(r.Context(), body.RefreshToken)
	response.JSON(w, http.StatusOK, map[string]string{"message": "successfully logged out"})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "unauthenticated")
		return
	}

	user, err := h.service.GetMe(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "user not found")
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, user)
}
