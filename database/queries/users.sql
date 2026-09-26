-- name: CreateUser :one
INSERT INTO users (
    id,
    first_name,
    last_name,
    username,
    email,
    password_hash
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING
    id,
    first_name,
    last_name,
    username,
    email,
    password_hash,
    created_at,
    updated_at,
    deleted_at;
