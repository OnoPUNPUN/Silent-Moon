package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type AnnotationType string

const (
	AnnotationTypeHighlight     AnnotationType = "highlight"
	AnnotationTypeUnderline     AnnotationType = "underline"
	AnnotationTypeStrikethrough AnnotationType = "strikethrough"
	AnnotationTypeNote          AnnotationType = "note"
)

type PDFRectangle struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type PDFLocation struct {
	Format       string         `json:"format"`
	Page         int            `json:"page"`
	SelectedText string         `json:"selectedText,omitempty"`
	Rectangles   []PDFRectangle `json:"rectangles,omitempty"`
	PageWidth    float64        `json:"pageWidth,omitempty"`
	PageHeight   float64        `json:"pageHeight,omitempty"`
}

type EPUBLocation struct {
	Format       string `json:"format"`
	CFI          string `json:"cfi"`
	SelectedText string `json:"selectedText,omitempty"`
	Prefix       string `json:"prefix,omitempty"`
	Suffix       string `json:"suffix,omitempty"`
}

type BaseLocation struct {
	Format string `json:"format"`
}

func ValidateLocation(raw json.RawMessage) error {
	if len(raw) == 0 {
		return errors.New("location cannot be empty")
	}

	var base BaseLocation
	if err := json.Unmarshal(raw, &base); err != nil {
		return fmt.Errorf("invalid location json: %w", err)
	}

	switch strings.ToLower(base.Format) {
	case "pdf":
		var pdfLoc PDFLocation
		if err := json.Unmarshal(raw, &pdfLoc); err != nil {
			return fmt.Errorf("invalid pdf location: %w", err)
		}
		if pdfLoc.Page <= 0 {
			return errors.New("pdf page must be greater than 0")
		}
	case "epub":
		var epubLoc EPUBLocation
		if err := json.Unmarshal(raw, &epubLoc); err != nil {
			return fmt.Errorf("invalid epub location: %w", err)
		}
		if strings.TrimSpace(epubLoc.CFI) == "" {
			return errors.New("epub cfi is required")
		}
	default:
		return fmt.Errorf("unsupported location format: %s (must be 'pdf' or 'epub')", base.Format)
	}

	return nil
}

type Annotation struct {
	ID           uuid.UUID       `json:"id"`
	UserID       uuid.UUID       `json:"userId"`
	BookID       uuid.UUID       `json:"bookId"`
	DeviceID     *uuid.UUID      `json:"deviceId,omitempty"`
	Type         AnnotationType  `json:"type"`
	SelectedText *string         `json:"selectedText,omitempty"`
	Color        *string         `json:"color,omitempty"`
	Note         *string         `json:"note,omitempty"`
	Location     json.RawMessage `json:"location"`
	Version      int64           `json:"version"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
	DeletedAt    *time.Time      `json:"deletedAt,omitempty"`
}
