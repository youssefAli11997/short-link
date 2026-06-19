package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"url-shortener/internal/model"
)

type PostgresURLRepository struct {
	db *pgxpool.Pool
}

func NewPostgresURLRepository(db *pgxpool.Pool) *PostgresURLRepository {
	return &PostgresURLRepository{
		db: db,
	}
}

func (r *PostgresURLRepository) CreateOrGet(
	ctx context.Context,
	originalURL string,
) (*model.URL, error) {

	query := `
		INSERT INTO urls(original_url)
		VALUES ($1)
		ON CONFLICT (original_url)
		DO UPDATE
			SET original_url = EXCLUDED.original_url -- effectively noop, just to return the record, for idempotency
		RETURNING
			id,
			original_url,
			created_at
	`

	var u model.URL

	err := r.db.QueryRow(
		ctx,
		query,
		originalURL,
	).Scan(
		&u.ID,
		&u.OriginalURL,
		&u.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *PostgresURLRepository) FindByID(
	ctx context.Context,
	id int64,
) (*model.URL, error) {

	query := `
		SELECT
			id,
			original_url,
			created_at
		FROM urls
		WHERE id = $1
	`

	var u model.URL

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&u.ID,
		&u.OriginalURL,
		&u.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrURLNotFound
	}

	if err != nil {
		return nil, err
	}

	return &u, nil
}
