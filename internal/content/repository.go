package content

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository persists Content records in Postgres.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates a Postgres-backed Content repository.
func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// Create inserts a new Content record.
func (r *PostgresRepository) Create(ctx context.Context, c *Content) error {
	const query = `
		INSERT INTO content (id, title, body, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`
	if _, err := r.pool.Exec(ctx, query, c.ID, c.Title, c.Body, c.CreatedAt, c.UpdatedAt); err != nil {
		return fmt.Errorf("content: insert: %w", err)
	}
	return nil
}

// List returns every Content record, newest first.
func (r *PostgresRepository) List(ctx context.Context) ([]Content, error) {
	const query = `
		SELECT id, title, body, created_at, updated_at
		FROM content
		ORDER BY created_at DESC, id`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("content: list: %w", err)
	}
	defer rows.Close()

	items := make([]Content, 0)
	for rows.Next() {
		var c Content
		if err := rows.Scan(&c.ID, &c.Title, &c.Body, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("content: scan list row: %w", err)
		}
		items = append(items, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("content: iterate list: %w", err)
	}

	return items, nil
}

// Get returns a Content record by id.
func (r *PostgresRepository) Get(ctx context.Context, id uuid.UUID) (*Content, error) {
	const query = `
		SELECT id, title, body, created_at, updated_at
		FROM content
		WHERE id = $1`

	var c Content
	err := r.pool.QueryRow(ctx, query, id).Scan(&c.ID, &c.Title, &c.Body, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("content: get: %w", err)
	}

	return &c, nil
}

// Update replaces the mutable fields of an existing Content record.
func (r *PostgresRepository) Update(ctx context.Context, c *Content) error {
	const query = `
		UPDATE content
		SET title = $2, body = $3, updated_at = $4
		WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, c.ID, c.Title, c.Body, c.UpdatedAt)
	if err != nil {
		return fmt.Errorf("content: update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete removes a Content record by id.
func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM content WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("content: delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
