// Package content implements the Content Management function of the
// Epileptic application component. It owns the Content data object that is
// accessed through the Communication Service.
package content

import (
	"context"
	"errors"
	"time"
)

// ErrNotFound is returned when a Content record does not exist.
var ErrNotFound = errors.New("content not found")

// ValidationError describes an invalid Content input.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// Content is the Content data object managed by the Epileptic application
// component.
type Content struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate checks the mandatory fields of a Content record.
func (c Content) Validate() error {
	if c.Title == "" {
		return ValidationError{Field: "title", Message: "must not be empty"}
	}
	return nil
}

// Repository persists Content records.
type Repository interface {
	Create(ctx context.Context, content *Content) error
	Get(ctx context.Context, id string) (*Content, error)
	List(ctx context.Context) ([]Content, error)
	Update(ctx context.Context, content *Content) error
	Delete(ctx context.Context, id string) error
}
