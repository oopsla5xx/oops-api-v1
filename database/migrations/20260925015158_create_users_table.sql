-- +goose Up

CREATE TABLE users (
    id UUID PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX users_username_unique
    ON users (username)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX users_email_unique
    ON users (email)
    WHERE deleted_at IS NULL;

-- +goose Down

DROP TABLE users;