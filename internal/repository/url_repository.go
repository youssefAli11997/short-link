package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"url-shortener/internal/model"
)

type URLRepository interface {
	CreateOrGet(ctx context.Context, originalURL string) (*model.URL, error)
	FindByID(ctx context.Context, id int64) (*model.URL, error)
}

type postgresURLRepository struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

func NewPostgresURLRepository(db *pgxpool.Pool) URLRepository {
	return &postgresURLRepository{
		db:      db,
		timeout: 5 * time.Second,
	}
}

func (r *postgresURLRepository) CreateOrGet(
	ctx context.Context,
	originalURL string,
) (*model.URL, error) {
	ctx, cancel := context.WithTimeout(
		ctx,
		r.timeout,
	)
	defer cancel()

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

func (r *postgresURLRepository) FindByID(
	ctx context.Context,
	id int64,
) (*model.URL, error) {
	ctx, cancel := context.WithTimeout(
		ctx,
		r.timeout,
	)
	defer cancel()

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
