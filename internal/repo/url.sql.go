// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: url.sql

package repo

import (
	"context"
	"time"
)

const createURL = `-- name: CreateURL :one
INSERT INTO urls (
    original_url,
    short_code,
    is_custom,
    expired_at
) VALUES (
    $1, $2, $3, $4
) RETURNING id, original_url, short_code, is_custom, expired_at, created_at
`

type CreateURLParams struct {
	OriginalUrl string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	IsCustom    bool      `json:"is_custom"`
	ExpiredAt   time.Time `json:"expired_at"`
}

func (q *Queries) CreateURL(ctx context.Context, arg CreateURLParams) (Url, error) {
	row := q.db.QueryRowContext(ctx, createURL,
		arg.OriginalUrl,
		arg.ShortCode,
		arg.IsCustom,
		arg.ExpiredAt,
	)
	var i Url
	err := row.Scan(
		&i.ID,
		&i.OriginalUrl,
		&i.ShortCode,
		&i.IsCustom,
		&i.ExpiredAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleUrlExpired = `-- name: DeleUrlExpired :exec
delete from urls where expired_at <= CURRENT_TIMESTAMP
`

func (q *Queries) DeleUrlExpired(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleUrlExpired)
	return err
}

const getUrlByShortCode = `-- name: GetUrlByShortCode :one
select id, original_url, short_code, is_custom, expired_at, created_at from urls 
where short_code=$1 
and expired_at > CURRENT_TIMESTAMP
`

func (q *Queries) GetUrlByShortCode(ctx context.Context, shortCode string) (Url, error) {
	row := q.db.QueryRowContext(ctx, getUrlByShortCode, shortCode)
	var i Url
	err := row.Scan(
		&i.ID,
		&i.OriginalUrl,
		&i.ShortCode,
		&i.IsCustom,
		&i.ExpiredAt,
		&i.CreatedAt,
	)
	return i, err
}

const isShortCodeAvailable = `-- name: IsShortCodeAvailable :one
select not exists (
    select 1 from urls where short_code=$1
) as is_available
`

func (q *Queries) IsShortCodeAvailable(ctx context.Context, shortCode string) (bool, error) {
	row := q.db.QueryRowContext(ctx, isShortCodeAvailable, shortCode)
	var is_available bool
	err := row.Scan(&is_available)
	return is_available, err
}
