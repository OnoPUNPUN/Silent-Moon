package response

import (
	"encoding/json"
	"net/http"
)

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

const (
	CodeValidationError    = "VALIDATION_ERROR"
	CodeUnauthorized       = "UNAUTHORIZED"
	CodeForbidden          = "FORBIDDEN"
	CodeNotFound           = "RESOURCE_NOT_FOUND"
	CodeConflict           = "CONFLICT"
	CodeDuplicateOperation = "DUPLICATE_OPERATION"
	CodeStorageError       = "STORAGE_ERROR"
	CodeInternalError      = "INTERNAL_ERROR"
)

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func Error(w http.ResponseWriter, status int, code, message string) {
	JSON(w, status, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

func ValidationError(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, CodeValidationError, message)
}

func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, CodeUnauthorized, message)
}

func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, CodeForbidden, message)
}

func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, CodeNotFound, message)
}

func Conflict(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, CodeConflict, message)
}

func DuplicateOperation(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, CodeDuplicateOperation, message)
}

func InternalError(w http.ResponseWriter, message string) {
	if message == "" {
		message = "An unexpected server error occurred."
	}
	Error(w, http.StatusInternalServerError, CodeInternalError, message)
}
