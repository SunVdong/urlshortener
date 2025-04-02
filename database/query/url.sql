-- name: CreateURL :one
INSERT INTO urls (
    original_url,
    short_code,
    is_custom,
    expired_at
) VALUES (
    $1, $2, $3, $4
) RETURNING *;


-- name: IsShortCodeAvailable :one
select not exists (
    select 1 from urls where short_code=$1
) as is_available;


-- name: GetUrlByShortCode :one
select * from urls 
where short_code=$1 
and expired_at > CURRENT_TIMESTAMP;


-- name: DeleUrlExpired :exec
delete from urls where expired_at <= CURRENT_TIMESTAMP;