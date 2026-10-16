package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// PostgresRepository stores the Product DataObject in a Postgres database.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository connects to the database identified by dsn, creates the
// product schema and seeds the catalog when it is empty.
func NewPostgresRepository(dsn string) (*PostgresRepository, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	repo := &PostgresRepository{db: db}
	if err := repo.migrate(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := repo.seed(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return repo, nil
}

const schema = `
CREATE TABLE IF NOT EXISTS product (
	id          TEXT PRIMARY KEY,
	name        TEXT    NOT NULL,
	category    TEXT    NOT NULL,
	description TEXT    NOT NULL DEFAULT '',
	price_cents BIGINT  NOT NULL DEFAULT 0,
	currency    TEXT    NOT NULL DEFAULT 'EUR',
	unit        TEXT    NOT NULL DEFAULT 'stuk',
	stock       INTEGER NOT NULL DEFAULT 0
);`

func (r *PostgresRepository) migrate(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("migrate product schema: %w", err)
	}
	return nil
}

func (r *PostgresRepository) seed(ctx context.Context) error {
	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM product`).Scan(&count); err != nil {
		return fmt.Errorf("count products: %w", err)
	}
	if count > 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin seed: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // rollback after commit is a no-op

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO product (id, name, category, description, price_cents, currency, unit, stock)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO NOTHING`)
	if err != nil {
		return fmt.Errorf("prepare seed: %w", err)
	}
	defer stmt.Close()

	for _, p := range SeedCatalog() {
		if _, err := stmt.ExecContext(ctx, p.ID, p.Name, string(p.Category), p.Description, p.PriceCents, p.Currency, p.Unit, p.Stock); err != nil {
			return fmt.Errorf("seed product %s: %w", p.ID, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit seed: %w", err)
	}
	return nil
}

// List returns all products ordered by category and name.
func (r *PostgresRepository) List(ctx context.Context) ([]Product, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, category, description, price_cents, currency, unit, stock
		FROM product
		ORDER BY category, name`)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		var category string
		if err := rows.Scan(&p.ID, &p.Name, &category, &p.Description, &p.PriceCents, &p.Currency, &p.Unit, &p.Stock); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		p.Category = Category(category)
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate products: %w", err)
	}
	return products, nil
}

// Get returns the product with the given id.
func (r *PostgresRepository) Get(ctx context.Context, id string) (Product, error) {
	var p Product
	var category string
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, category, description, price_cents, currency, unit, stock
		FROM product
		WHERE id = $1`, id).
		Scan(&p.ID, &p.Name, &category, &p.Description, &p.PriceCents, &p.Currency, &p.Unit, &p.Stock)
	if errors.Is(err, sql.ErrNoRows) {
		return Product{}, ErrNotFound
	}
	if err != nil {
		return Product{}, fmt.Errorf("get product %s: %w", id, err)
	}
	p.Category = Category(category)
	return p, nil
}

// Close releases the database connection.
func (r *PostgresRepository) Close() error { return r.db.Close() }
