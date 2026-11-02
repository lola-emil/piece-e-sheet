-- +goose Up
SELECT 'up SQL query';

CREATE TABLE users (
    id UUID PRIMARY KEY,

    email VARCHAR(255) NOT NULL UNIQUE,
    display_name VARCHAR(100) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
SELECT 'down SQL query';
DROP TABLE users;