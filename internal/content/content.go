// Package content implements the Content Management function of the
// Epileptic application component. It owns access to the Content data object.
package content

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrNotFound is returned when a Content record does not exist.
var ErrNotFound = errors.New("content: not found")

// Content is the data object managed by the Content Management function.
type Content struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Repository persists Content records.
type Repository interface {
	Create(ctx context.Context, c *Content) error
	List(ctx context.Context) ([]Content, error)
	Get(ctx context.Context, id uuid.UUID) (*Content, error)
	Update(ctx context.Context, c *Content) error
	Delete(ctx context.Context, id uuid.UUID) error
}
