package main

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
)

// postgresStore persists Content in the default local postgres service. The
// Content Data Object is modelled in the System Design view, so a relational
// store backs the Content Management function.
type postgresStore struct {
	db *sql.DB
}

func newPostgresStore(dsn string) (*postgresStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	store := &postgresStore{db: db}
	if err := store.migrate(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *postgresStore) migrate(ctx context.Context) error {
	const schema = `
CREATE TABLE IF NOT EXISTS content (
    id         TEXT PRIMARY KEY,
    title      TEXT NOT NULL,
    body       TEXT NOT NULL DEFAULT '',
    author     TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);`
	_, err := s.db.ExecContext(ctx, schema)
	return err
}

func (s *postgresStore) List(ctx context.Context) ([]Content, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, title, body, author, created_at, updated_at
FROM content
ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Content, 0)
	for rows.Next() {
		item, err := scanContent(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *postgresStore) Get(ctx context.Context, id string) (Content, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, title, body, author, created_at, updated_at
FROM content
WHERE id = $1`, id)

	item, err := scanContent(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Content{}, ErrContentNotFound
	}
	return item, err
}

func (s *postgresStore) Create(ctx context.Context, item Content) (Content, error) {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO content (id, title, body, author, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)`,
		item.ID, item.Title, item.Body, item.Author, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return Content{}, err
	}
	return item, nil
}

func (s *postgresStore) Update(ctx context.Context, item Content) (Content, error) {
	result, err := s.db.ExecContext(ctx, `
UPDATE content
SET title = $2, body = $3, author = $4, updated_at = $5
WHERE id = $1`,
		item.ID, item.Title, item.Body, item.Author, item.UpdatedAt)
	if err != nil {
		return Content{}, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return Content{}, err
	}
	if affected == 0 {
		return Content{}, ErrContentNotFound
	}
	return item, nil
}

func (s *postgresStore) Delete(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM content WHERE id = $1`, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrContentNotFound
	}
	return nil
}

func (s *postgresStore) Close() {
	_ = s.db.Close()
}

// scanner is satisfied by both *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...any) error
}

func scanContent(row scanner) (Content, error) {
	var item Content
	if err := row.Scan(
		&item.ID,
		&item.Title,
		&item.Body,
		&item.Author,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return Content{}, err
	}
	return item, nil
}
