package content

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository is the Repository implementation backed by the local
// Postgres service declared in .plot/package.json.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository returns a Repository backed by the given pool.
func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// Create inserts a new Content record.
func (r *PostgresRepository) Create(ctx context.Context, content *Content) error {
	const query = `
		INSERT INTO content (id, title, body, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	if _, err := r.pool.Exec(ctx, query,
		content.ID, content.Title, content.Body, content.CreatedAt, content.UpdatedAt,
	); err != nil {
		return fmt.Errorf("insert content: %w", err)
	}
	return nil
}

// Get returns a single Content record by id.
func (r *PostgresRepository) Get(ctx context.Context, id string) (*Content, error) {
	const query = `
		SELECT id, title, body, created_at, updated_at
		FROM content
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)
	return scanContent(row)
}

// List returns all Content records ordered by creation time.
func (r *PostgresRepository) List(ctx context.Context) ([]Content, error) {
	const query = `
		SELECT id, title, body, created_at, updated_at
		FROM content
		ORDER BY created_at ASC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list content: %w", err)
	}
	defer rows.Close()

	items := make([]Content, 0)
	for rows.Next() {
		item, err := scanContent(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list content: %w", err)
	}
	return items, nil
}

// Update persists the mutable fields of an existing Content record.
func (r *PostgresRepository) Update(ctx context.Context, content *Content) error {
	const query = `
		UPDATE content
		SET title = $2, body = $3, updated_at = $4
		WHERE id = $1
	`
	tag, err := r.pool.Exec(ctx, query,
		content.ID, content.Title, content.Body, content.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update content: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete removes a Content record by id.
func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM content WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete content: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// scanner is satisfied by both pgx.Row and pgx.Rows.
type scanner interface {
	Scan(dest ...any) error
}

func scanContent(row scanner) (*Content, error) {
	var content Content
	if err := row.Scan(
		&content.ID, &content.Title, &content.Body,
		&content.CreatedAt, &content.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("scan content: %w", err)
	}
	return &content, nil
}
